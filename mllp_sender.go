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
)

var (
	host         = flag.String("host", "", "Host to send data to.")
	port         = flag.Int("port", -1, "Port # to send data to.")
	mllpStartStr = flag.String("mllp-start", "11", "MLLP start character(s).")
	mllpEndStr   = flag.String("mllp-end", "28,13", "MLLP ending character(s).")
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
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	var dataToSend []byte
	dataToSend = append(dataToSend, mllpStart[:]...)
	dataToSend = append(dataToSend, []byte(msg)...)
	dataToSend = append(dataToSend, mllpEnd[:]...)

	_, err = conn.Write(dataToSend)

	if err != nil {
		log.Fatal(err)
	}

	connBuf := bufio.NewReader(conn)

	for {

		slc, err := connBuf.ReadSlice(mllpEnd[len(mllpEnd)-1])

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		if len(slc) > 0 {
			fmt.Printf("%s\n", string(slc))
		}

	}
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
