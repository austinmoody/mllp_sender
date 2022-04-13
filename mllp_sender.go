package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	host         = flag.String("host", "", "Host to send data to.")
	port         = flag.Int("port", -1, "Port # to send data to.")
	mllpStartStr = flag.String("mllp-start",
		"11",
		"MLLP start character(s), specified as decimal values.\nFor example Vertical Tab = 11.")
	mllpEndStr = flag.String("mllp-end",
		"28,13",
		"MLLP ending character(s), specified as decimal values.\nFor example File Separator = 28.")
	timeout = flag.String("timeout",
		"10s",
		"Timeout to stop listening for response.\nTo be specified in format understood by ParseDuration.\nSee https://pkg.go.dev/time#ParseDuration")
)

func main() {

	flag.Parse()

	if *host == "" {
		flag.PrintDefaults()
		fmt.Printf("\n\n*** Host was not specified. ***\n\n")
		os.Exit(0)
	}

	if *port == -1 {
		flag.PrintDefaults()
		fmt.Printf("\n\n*** Port was not specified. ***\n\n")
		os.Exit(0)
	}

	mllpStart := StringSplitToByteArray(*mllpStartStr, ",")
	mllpEnd := StringSplitToByteArray(*mllpEndStr, ",")

	msg := GetStdinData()

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer CloseConnection(conn)

	var dataToSend []byte
	dataToSend = append(dataToSend, mllpStart[:]...)
	dataToSend = append(dataToSend, []byte(msg)...)
	dataToSend = append(dataToSend, mllpEnd[:]...)

	_, err = conn.Write(dataToSend)

	if err != nil {
		log.Fatal(err)
	}

	// Set timeout to close connection and exit
	// if we aren't getting a response after our submission.
	// Example: An HL7 listener that isn't setup to respond w/ ACK messages
	timeoutDuration, err := time.ParseDuration(*timeout)
	if err != nil {
		log.Fatal(err)
	}
	time.AfterFunc(timeoutDuration, func() {
		CloseConnection(conn)
		fmt.Printf("Timeout (%s) while listening for response.", timeoutDuration.String())
		os.Exit(0)
	})

	fmt.Println("Data Sent")
	fmt.Printf("\tMLLP start character(s): %d (decimal)\n", mllpStart)
	fmt.Printf("\tMLLP end   character(s): %d (decimal)\n", mllpEnd)

	connBuf := bufio.NewReader(conn)
	var fullResponse []byte
	for {
		bt, err := connBuf.ReadByte()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		fullResponse = append(fullResponse, bt)

		// Check to see if we have gotten MLLP End character(s)
		if len(fullResponse) >= len(mllpEnd) && bytes.Equal(mllpEnd, fullResponse[len(fullResponse)-len(mllpEnd):]) {
			break
		}
	}

	responseString := string(fullResponse)
	// MLLP start and end characters can interfere w/ output to command line
	ripNonPrintable := func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return '\n'
	}
	responseString = strings.Map(ripNonPrintable, responseString)

	fmt.Printf("Server Response:\n\n%s\n\n", responseString)
}

func StringSplitToByteArray(stringToSplit string, delim string) []byte {
	i := 0
	var returnBytes []byte
	for _, mllpChar := range strings.Split(stringToSplit, delim) {
		mllpInt, err := strconv.Atoi(mllpChar)
		if err != nil {
			panic(err)
		}
		returnBytes = append(returnBytes, byte(mllpInt))
		i++
	}

	return returnBytes
}

func GetStdinData() string {
	buf := bytes.Buffer{}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		buf.Write(s.Bytes())
	}

	return buf.String()
}

func CloseConnection(conn *net.TCPConn) {
	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
