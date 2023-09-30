package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
)

const MAX_LINK_ID_LENGTH = 8

var (
	links                        map[string]string
	ErrEmptyLinkError            = errors.New("The link is empty")
	ErrLinkContainsJustURLScheme = errors.New("The link contains only URL scheme")
)

func makeAndStoreShortURL(url string) (string, error) {
	if err := validateURL(url); err != nil {
		return "", err
	}
	if links == nil {
		links = make(map[string]string)
	}

	hash := md5.New()
	io.WriteString(hash, url)
	encodedString := fmt.Sprintf("%x", hash.Sum(nil))
	if len([]rune(encodedString)) < MAX_LINK_ID_LENGTH {
		links[encodedString] = url
		return encodedString, nil
	} else {
		links[encodedString[:MAX_LINK_ID_LENGTH]] = url
		return encodedString[:MAX_LINK_ID_LENGTH], nil
	}
}

func validateURL(url string) error {
	if len(url) == 0 {
		return ErrEmptyLinkError
	}
	if url == "https://" || url == "http://" {
		return ErrLinkContainsJustURLScheme
	}
	return nil
}
