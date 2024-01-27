package storage

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/YuryBogdanov/link-shortener/internal/config"
	"github.com/YuryBogdanov/link-shortener/internal/logger"
	uuid "github.com/satori/go.uuid"
)

const MaxLinkIDLength = 8

var (
	Links          = make(map[string]string)
	lock           = sync.RWMutex{}
	errNoSuchValue = errors.New("no such value")
	lg             = logger.DefaultLogger{}
)

type StorableLink struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func SetupPersistentStorage(fileName string) {
	lg.Setup()
	defer lg.Finish()
	Storage = FileStorage{FilePath: fileName}
	storedLinks := Storage.GetAllItems()
	for _, link := range storedLinks {
		Links[link.ShortURL] = link.OriginalURL
	}
}

func MakeAndStoreShortURL(url string) (string, error) {
	hash := md5.New()
	io.WriteString(hash, url)
	encodedString := fmt.Sprintf("%x", hash.Sum(nil))

	var storageKey string
	if len([]rune(encodedString)) < MaxLinkIDLength {
		storageKey = encodedString
	} else {
		storageKey = encodedString[:MaxLinkIDLength]
	}

	resultLink := getShortenedLink(storageKey)

	_, linkAlredyStored := Links[storageKey]
	if linkAlredyStored {
		return resultLink, nil
	}

	setLinkForKey(storageKey, url)
	return resultLink, nil
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

	l := StorableLink{
		UUID:        string(uuid.NewV4().String()),
		ShortURL:    key,
		OriginalURL: link,
	}
	err := Storage.Store(l)
	if err != nil {
		lg.Error("couldn't store a new record", l)
	}
}
