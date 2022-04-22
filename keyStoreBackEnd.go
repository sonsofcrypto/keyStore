// SPDX-License-Identifier: MIT

package keyStore

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// BackEnd `Item`s can be stored at disk, system keychain etc.
type BackEnd interface {
	// List all items in `BackEnd`
	List() []*KeyStoreItem
	// Add `KeyStoreItem` to `BackEnd`
	Add(item *KeyStoreItem)
	// Remove `KeyStoreItem` to `BackEnd`
	Remove(item *KeyStoreItem)
}

// DiskBackEnd stores `Item`s at disk
type DiskBackEnd struct {
	storePath string
}

// NewDiskBackEnd checks if there is existing folder and path, it not attempts
// to create one. Panics if unable to create folder when needed
func NewDiskBackEnd(storePath string) *DiskBackEnd {
	cleanPath := path.Clean(storePath)
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		if err = os.MkdirAll(cleanPath, os.ModePerm); err != nil {
			log.Panicln("Failed to create `DiskBackEnd`", cleanPath, err)
		}
	}
	return &DiskBackEnd{
		storePath: path.Clean(cleanPath),
	}
}

// List all the files at `storePath` and unmarshals them to `KeyStoreItem`s
func (d *DiskBackEnd) List() ([]*KeyStoreItem, error) {
	items := make([]*KeyStoreItem, 0)
	files, err := ioutil.ReadDir(d.storePath)
	if err != nil {
		log.Panicln(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		f, err := os.Open(d.storePath + file.Name())
		if err != nil {
			return nil, err
		}
		defer f.Close()
		item := new(KeyStoreItem)
		if err := json.NewDecoder(f).Decode(item); err != nil {
			log.Println("Failed to unmarshal file", f, err)
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

// Add json serializes `KeyStoreItem` to file
func (d *DiskBackEnd) Add(item *KeyStoreItem) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return os.WriteFile(d.storePath+item.uuid+".json", data, 0644)
}

// Remove file from store
func (d *DiskBackEnd) Remove(item *KeyStoreItem) {
	os.Remove(d.filePath(item))
}

func (d *DiskBackEnd) filePath(item *KeyStoreItem) string {
	return d.storePath + item.uuid + ".json"
}
