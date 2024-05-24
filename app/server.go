package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"slices"
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
	url := extractURLPath(parts)
	urlParts := strings.Split(url, "/")
	switch len(urlParts) {
	case 2:
		if urlParts[1] != "" {
			resp := createResponseWithHeader(404, "", nil)
			conn.Write(resp)
			return
		}
		conn.Write(createResponseWithHeader(200, "", nil))
	case 3:
		if urlParts[1] != "echo" {
			panic(fmt.Sprintf("URL is %q, can not handle it", url))
		}
		conn.Write(createResponseWithHeader(200, "text/plain", []byte(urlParts[2])))
	default:
		panic(fmt.Sprintf("URL is %q, can not handle it", url))
	}
}

func createResponse(status int) []byte {
	var b strings.Builder
	b.Write(protocolBytes)
	b.Write(spaceBytes)
	b.WriteString(fmt.Sprintf("%d", status))
	b.Write(spaceBytes)
	b.WriteString(http.StatusText(status))
	return []byte(b.String())
}

func createResponseWithHeader(status int, contentType string, body []byte) []byte {
	resp := createResponse(status)
	header := respHeader{
		ContentType:   contentType,
		ContentLength: len(body),
	}
	headerBytes := header.toBytes()
	if headerBytes == nil {
		return slices.Concat(resp, crlfBytes)
	}
	return slices.Concat(resp, crlfBytes, headerBytes, crlfBytes, crlfBytes, body)
}

func extractURLPath(req [][]byte) string {
	reqLine := req[0]
	reqLineSplits := bytes.Split(reqLine, spaceBytes)
	return string(reqLineSplits[1])
}
