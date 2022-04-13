# mllp_sender

This is a simple TCP/IP client which reads data via STDIN and posts it to an MLLP server.

```
cat ~/Downloads/single.hl7| mllp_sender -host 10.0.0.50 -port 6661
Data Sent
	MLLP start character(s): [11] (decimal)
	MLLP end   character(s): [28 13] (decimal)
Server Response:


MSH|^~\&|HIE|HIE|Mirth|Hospital|20220412205011.204||ACK^A07^ACK|20220412205011.204|P|2.5.1
MSA|AA|e64b7cda-1e27-40f2-86ef-97016909ae16
```

## Getting It

Grab a file for your platform from the latest release: https://github.com/austinmoody/mllp_sender/releases

Install via go install:

```
go install github.com/austinmoody/mllp_sender@latest
```
## Usage

```
Usage of mllp_sender:
  -host string
    	Host to send data to.
  -mllp-end string
    	MLLP ending character(s), specified as decimal values.
    	For example File Separator = 28. (default "28,13")
  -mllp-start string
    	MLLP start character(s), specified as decimal values.
    	For example Vertical Tab = 11. (default "11")
  -port int
    	Port # to send data to. (default -1)
  -timeout string
    	Timeout to stop listening for response.
    	To be specified in format understood by ParseDuration.
    	See https://pkg.go.dev/time#ParseDuration (default "10s")
```

### MLLP Start & End

If these are not specified, a _default_ is used which are the values typically seen in Healthcare Integration:

* MLLP Start: Set to ASCII Vertical Tab (decimal 11).
* MLLP End: Set to ASCII File Separator (decimal 28) + ASCII Carriage Return (decimal 13).

If the MLLP server requires different MLLP wrapping, you may specify that with the _mllp-start_ and _mllp_end_ options.

The values should be specified using the decimal representation of the necessary character.  For example, if MLLP start of Group Separator and MLLP end of Vertical Tab + Carriage Return are required, you would specify with:

```
mllp_sender -mllp-start 29 -mllp-end 11,13
```

### Timeout

This utility will close its connection to the MLLP server after it has submitted data and received a response.  

When a response is received, it determines the completion using the same MLLP wrapping as used to send the message.  After this the utility will close the connection and exit.

If the MLLP server is not configured to respond with a message **or** the MLLP server is configured to keep connections open then the utility will close the connection after the specified timeout.

If not specified, the timeout is 10 seconds.  You may specify this using any time span which is able to be parsed by ParseDuration.  See https://pkg.go.dev/time#ParseDuration.