package main

import (
	"bytes"
	"fmt"
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
	requestBody := make([]byte, 1024)
	_, err = conn.Read(requestBody)
	if err != nil {
		fmt.Println("Could not read the connection: ", err.Error())
		os.Exit(1)
	}
	fmt.Printf("request: %s\n", requestBody)
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
	b.Write(spaceBytes)
	b.WriteString(http.StatusText(status))
	b.WriteString(string(crlfBytes))
	b.WriteString(string(crlfBytes))
	return []byte(b.String())
}

func isSlashRequest(req [][]byte) bool {
	reqLine := req[0]
	reqLineSplits := bytes.Split(reqLine, spaceBytes)
	return bytes.Equal(reqLineSplits[1], []byte("/"))
}
