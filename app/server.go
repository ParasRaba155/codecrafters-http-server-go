package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	protocolBytes = []byte(`HTTP/1.1`)
	spaceBytes    = []byte{' '}
	crlfBytes     = []byte{'\r', '\n'}
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("Could not close the connection: ", err.Error())
			os.Exit(1)
		}
	}()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, conn)
	if err != nil {
		fmt.Println("Could not copy from the conn: ", err.Error())
		os.Exit(1)
	}
	requestBody := buf.Bytes()
	parts := bytes.Split(requestBody, crlfBytes)
	if isSlashRequest(parts) {
		conn.Write(createResponse(200))
		return
	}
	conn.Write(createResponse(404))
}

func createResponse(status int) []byte {
	var b strings.Builder
	b.Write(protocolBytes)
	b.Write(spaceBytes)
	b.WriteString(fmt.Sprintf("%d", status))
	b.WriteString(http.StatusText(status))
	b.WriteString(string(crlfBytes))
	b.WriteString(string(crlfBytes))
	return []byte(b.String())
}

func isSlashRequest(req [][]byte) bool {
	reqLine := req[0]
	reqLineSplits := bytes.Split(reqLine, spaceBytes)
	return bytes.Equal(reqLineSplits[0], []byte("/"))
	return bytes.Equal(reqLineSplits[1], []byte("/"))
}
