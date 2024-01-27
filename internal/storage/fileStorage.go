package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Storager interface {
	Store(record StorableLink) error
	Get(id string) StorableLink
	GetAllItems() []StorableLink
}

type FileStorage struct {
	FilePath string
}

var Storage Storager

func (fs FileStorage) Store(record StorableLink) error {
	dir := filepath.Dir(fs.FilePath)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}
	file, err := os.OpenFile(fs.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	json, err := json.Marshal(record)
	if err != nil {
		return err
	}
	file.WriteString(string(json))
	file.WriteString("\n")
	file.Close()
	return nil
}

func (fs FileStorage) Get(id string) StorableLink {
	return StorableLink{}
}

func (fs FileStorage) GetAllItems() []StorableLink {
	file, err := os.OpenFile(fs.FilePath, os.O_RDONLY, 0777)
	if err != nil {
		lg.Info("couldn't open storage file, skipping data restoration")
		return make([]StorableLink, 0)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	links := make([]StorableLink, 0)

	for scanner.Scan() {
		text := scanner.Text()
		data := []byte(text)
		var obj *StorableLink
		json.Unmarshal(data, &obj)
		fmt.Println(obj)
		links = append(links, *obj)
	}

	return links
}
