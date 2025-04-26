package main

import (
	"fmt"
)

// respHeader the simple Header type
type respHeader struct {
	ContentType     string
	ContentEncoding string
	ContentLength   int
	CloseConnection bool
}

// toBytes to get headers in byte format.
func (r respHeader) toBytes() []byte {
	if r.ContentType == "" {
		return nil
	}
	closeConnectionStr := ""
	if r.CloseConnection {
		closeConnectionStr = "Connection: close\r\n"
	}
	if r.ContentEncoding == "" {
		str := fmt.Sprintf(
			"Content-Type: %s\r\nContent-Length: %d\r\n%s",
			r.ContentType,
			r.ContentLength,
			closeConnectionStr,
		)
		return []byte(str)
	}
	str := fmt.Sprintf(
		"Content-Encoding: %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n%s",
		r.ContentEncoding,
		r.ContentType,
		r.ContentLength,
		closeConnectionStr,
	)
	return []byte(str)
}
