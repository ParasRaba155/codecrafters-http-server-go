package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
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
	colonBytes    = []byte{':'}
)

var dirToLook = flag.String("directory", "", "to look into the directory")

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	fmt.Println("Ready to serve on PORT=4221")
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		}
		go handleConnection(conn)
	}
}

// handleConnection will handle each connection corresponding to the
// rules of the problem
//
// NOTE: Should be used as a separate go routine, to handle each connection concurrently
func handleConnection(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("Could not close the connection: ", err.Error())
			os.Exit(1)
		}
	}()

	// read the request in sufficiently large buffer
	requestBody := make([]byte, 1024)
	_, err := conn.Read(requestBody)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Println("Connection Closed")
			os.Exit(0)
		}
		fmt.Println("Could not read the connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Printf("request: %s\n", requestBody)

	// the request is split on '\r\n' into diff parts: info, header, and body
	parts := bytes.Split(requestBody, crlfBytes)
	url := extractURLPath(parts)
	fmt.Printf("URL: %q\n", url)
	urlParts := strings.Split(url, "/")
	switch len(urlParts) {
	case 2:
		switch urlParts[1] {
		// handle "/"
		case "":
			conn.Write(CreateResponseWithHeader(200, "", nil))
		case "user-agent":
			conn.Write(CreateResponseWithHeader(200, "text/plain", []byte(getUserAgent(parts))))
		default:
			conn.Write(CreateResponseWithHeader(404, "", nil))
		}
	case 3:
		switch urlParts[1] {
		// handle "/echo/{str}"
		case "echo":
			conn.Write(CreateResponseWithHeader(200, "text/plain", []byte(urlParts[2])))
		// handle "/files/{filename}"
		case "files":
			filePath := urlParts[2]
			conn.Write(GetFile(filePath))
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

// extractHeaders will read request bytes and extract headers from it
func extractHeaders(req [][]byte) http.Header {
	headers := http.Header{}

	// first part is GET <url> <protocol>, ignore that
	for i := 1; i < len(req); i++ {
		headerBytes := bytes.SplitN(req[i], colonBytes, 2)
		if len(headerBytes) != 2 {
			continue
		}
		key, value := bytes.TrimLeft(headerBytes[0], " "), bytes.TrimLeft(headerBytes[1], " ")
		headers.Add(string(key), string(value))
	}
	return headers
}

func getUserAgent(req [][]byte) string {
	headers := extractHeaders(req)
	return headers.Get("User-Agent")
}

func GetFile(filename string) []byte {
	fmt.Printf("DEBUG dir: %q, filename: %q\n", *dirToLook, filename)
	return nil
}
