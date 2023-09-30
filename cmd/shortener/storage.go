package main

import (
	"crypto/md5"
	"fmt"
	"io"
)

const MAX_LINK_ID_LENGTH = 8

var links map[string]string

func makeAndStoreShortURL(url string) string {
	hash := md5.New()
	io.WriteString(hash, url)
	encodedString := fmt.Sprintf("%x", hash.Sum(nil))
	if len([]rune(encodedString)) < MAX_LINK_ID_LENGTH {
		links[encodedString] = url
		return encodedString
	} else {
		links[encodedString[:MAX_LINK_ID_LENGTH]] = url
		return encodedString[:MAX_LINK_ID_LENGTH]
	}
}
