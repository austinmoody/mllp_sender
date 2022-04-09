package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	mllpStartStr = flag.String("mllp-start", "11", "MLLP start character(s)")
	mllpEndStr   = flag.String("mllp-end", "28,13", "MLLP ending character(s)")
)

func main() {

	flag.Parse()

	mllpStart := StringSplitToByteArray(*mllpStartStr, ",")
	mllpEnd := StringSplitToByteArray(*mllpEndStr, ",")

	hl7 := `MSH|^~\&|Mirth|Hospital|HIE|HIE|20220405211727||ADT^A07|62f241e7-37de-404a-b1d0-a428160ef940|P|2.5.1\n
EVN|A07|20220405211727\n
PID|1||PG23FK030030^2^BCV^^MR~TO64MN065772^2^BCV^^AN~LR78DM885590^8^BCV^^AND~ID75SM136214^4^NPI^^ANT||Patrickson^Lillis^I^^|Tallboy|19980118|F||2076-8^Native Hawaiian or Other Pacific Islander^HL70005|76255 Harbort Park^Number 18^Memphis^TN^38136^US^M||5791303873^PRN^CP^^+234^579^1303873^^^^^2345791303873~4108438947^ORN^PH^^+86^410^8438947^^^^^864108438947~^NET^X.400^rscudders3@hexun.com^^^^^^^^~^NET^X.400^ogreed4@globo.com^^^^^^^^|5689724029^WPN^CP^^+62^568^9724029^^^^^625689724029|pl|O^Other^HL70002|HOT|TO64MN065772^2^BCV^^AN|661-41-7685|1913688611^TN^20230328||N||Y|5||||20220327|Y\n
PV1|1|C`

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "10.0.0.50", 6661))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	_, err = conn.Write(mllpStart)
	_, err = conn.Write([]byte(hl7))
	_, err = conn.Write(mllpEnd)

	connBuf := bufio.NewReader(conn)

	for {

		slc, err := connBuf.ReadSlice(mllpEnd[len(mllpEnd)-1])

		if err != nil {
			break
		}

		if len(slc) > 0 {
			fmt.Printf("|%s|\r\n", string(slc))
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
