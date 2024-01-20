package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

type Storager interface {
	Store(record interface{}) error
	Get(id string, objectType reflect.Type) interface{}
}

type FileStorage struct {
	FilePath string
}

var Storage Storager

func (fs FileStorage) Store(record interface{}) error {
	dir := filepath.Dir(fs.FilePath)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}
	file, err := os.OpenFile(fs.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	fmt.Println(err)
	if err != nil {
		return err
	}
	fmt.Println(fs.FilePath)
	json, err := json.Marshal(record)
	if err != nil {
		return err
	}
	file.WriteString(string(json))
	file.WriteString("\n")
	file.Close()
	return nil
}

func (fs FileStorage) Get(id string, objectType reflect.Type) interface{} {
	return nil
}
