package keyStore

import (
	"github.com/google/uuid"
)

type KeyStoreItem struct {
	uuid string
}

func NewKeyStoreItem() *KeyStoreItem {
	return &KeyStoreItem{
		uuid: uuid.New().String(),
	}
}
