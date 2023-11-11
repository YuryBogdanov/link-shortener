package storage

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/YuryBogdanov/link-shortener/internal/config"
)

const MaxLinkIDLength = 8

var (
	Links          map[string]string
	lock           = sync.RWMutex{}
	errNoSuchValue = errors.New("no such value!")
)

func MakeAndStoreShortURL(url string) (string, error) {
	if Links == nil {
		Links = make(map[string]string)
	}

	hash := md5.New()
	io.WriteString(hash, url)
	encodedString := fmt.Sprintf("%x", hash.Sum(nil))
	if len([]rune(encodedString)) < MaxLinkIDLength {
		setLinkForKey(encodedString, url)
		resultLink := getShortenedLink(encodedString)
		return resultLink, nil
	} else {
		maxID := encodedString[:MaxLinkIDLength]
		setLinkForKey(maxID, url)
		resultLink := getShortenedLink(maxID)
		return resultLink, nil
	}
}

func GetLinkForKey(key string) (string, error) {
	lock.RLock()
	linkToReturn, ok := Links[key]
	lock.RUnlock()
	if ok {
		return linkToReturn, nil
	} else {
		return "", errNoSuchValue
	}

}

func getShortenedLink(linkID string) string {
	return config.BaseConfig.ShoretnedBaseURL.Value + "/" + linkID
}

func setLinkForKey(key string, link string) {
	lock.Lock()
	Links[key] = link
	lock.Unlock()
}
