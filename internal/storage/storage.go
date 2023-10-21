package storage

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
)

const MaxLinkIDLength = 8

var (
	Links                        map[string]string
	ErrEmptyLinkError            = errors.New("the link is empty")
	ErrLinkContainsJustURLScheme = errors.New("the link contains only URL scheme")
)

func MakeAndStoreShortURL(url string) (string, error) {
	if err := validateURL(url); err != nil {
		return "", err
	}
	if Links == nil {
		Links = make(map[string]string)
	}

	hash := md5.New()
	io.WriteString(hash, url)
	encodedString := fmt.Sprintf("%x", hash.Sum(nil))
	if len([]rune(encodedString)) < MaxLinkIDLength {
		Links[encodedString] = url
		return encodedString, nil
	} else {
		Links[encodedString[:MaxLinkIDLength]] = url
		return encodedString[:MaxLinkIDLength], nil
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
