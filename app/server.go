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

	// the request is split on '\r\n' into diff parts: info, header, and body
	parts := bytes.Split(requestBody, crlfBytes)
	url := extractURLPath(parts)
	fmt.Printf("URL: %q", url)
	urlParts := strings.Split(url, "/")
	switch len(urlParts) {
	case 2:
		switch urlParts[1] {
		// handle "/"
		case "":
			conn.Write(CreateResponseWithHeader(200, "", nil))
		default:
			conn.Write(CreateResponseWithHeader(404, "", nil))
		}
	case 3:
		switch urlParts[1] {
		// handle "/echo/{str}"
		case "echo":
			conn.Write(CreateResponseWithHeader(200, "text/plain", []byte(urlParts[2])))
		default:
			conn.Write(CreateResponseWithHeader(404, "", nil))
			fmt.Printf("URL is %q, can not handle it", url)
		}
	default:
		conn.Write(CreateResponseWithHeader(404, "", nil))
		fmt.Printf("URL is %q, can not handle it", url)
	}
}

// createBasicResponse the base for creating a simple response with only status,
// no header and no body
//
// E.g: HTTP/1.1 200 Ok, HTTP/1.1 404 Not Found
func createBasicResponse(status int) []byte {
	var b strings.Builder
	b.Write(protocolBytes)
	b.Write(spaceBytes)
	b.WriteString(fmt.Sprintf("%d", status))
	b.Write(spaceBytes)
	b.WriteString(http.StatusText(status))
	return []byte(b.String())
}

// CreateResponseWithHeader will create response with status, header according
// to content type and body
//
// NOTE:
//   - For no header pass contentType as empty string
//   - For no body pass body as nil
func CreateResponseWithHeader(status int, contentType string, body []byte) []byte {
	statusPart := createBasicResponse(status)
	header := respHeader{
		ContentType:   contentType,
		ContentLength: len(body),
	}
	headerBytes := header.toBytes()
	if headerBytes == nil {
		return slices.Concat(statusPart, crlfBytes, crlfBytes)
	}
	return slices.Concat(statusPart, crlfBytes, headerBytes, crlfBytes, body)
}

// extractURLPath is stored in the first part of the request
func extractURLPath(req [][]byte) string {
	reqLine := req[0]
	reqLineSplits := bytes.Split(reqLine, spaceBytes)
	// first is the protocol info and 2nd one is the URL path
	return string(reqLineSplits[1])
}
