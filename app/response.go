package main

import (
	"fmt"
)

// respHeader the simple Header type
type respHeader struct {
	ContentType     string
	ContentEncoding string
	ContentLength   int
}

// toBytes to get headers in byte format.
func (r respHeader) toBytes() []byte {
	if r.ContentType == "" {
		return nil
	}
	if r.ContentEncoding == "" {
		str := fmt.Sprintf(
			"Content-Type: %s\r\nContent-Length: %d\r\n",
			r.ContentType,
			r.ContentLength,
		)
		return []byte(str)
	}
	str := fmt.Sprintf(
		"Content-Encoding: %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n",
		r.ContentEncoding,
		r.ContentType,
		r.ContentLength,
	)
	return []byte(str)
}
