package main

import (
	"fmt"
)

type respHeader struct {
	ContentType   string
	ContentLength int
}

// toBytes implements json.Marshaler.
func (r respHeader) toBytes() []byte {
	if r.ContentType == "" {
		return nil
	}
	str := fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d", r.ContentType, r.ContentLength)
	return []byte(str)
}
