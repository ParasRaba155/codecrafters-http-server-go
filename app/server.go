package main

import (
	"fmt"
	"net"
	"os"
	"slices"
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
	_, err = conn.Write(createStatusOkResponse())
	if err != nil {
		fmt.Println("Error Writing Status OK: ", err.Error())
		os.Exit(1)
	}
}

func createStatusOkResponse() []byte {
	statusBytes := []byte(`200 OK`)
	return slices.Concat(protocolBytes, spaceBytes, statusBytes, crlfBytes, crlfBytes)
}
