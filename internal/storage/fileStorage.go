package storage

import (
	"encoding/json"
	"os"
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
	file, err := os.OpenFile(fs.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

func (fs FileStorage) Get(id string, objectType reflect.Type) interface{} {
	return nil
}
