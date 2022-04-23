// SPDX-License-Identifier: MIT

package keyStore

import (
	"crypto/ecdsa"
	cryptorand "crypto/rand"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"log"
	"os"
	"testing"
)

func TestDiskKeyStore(t *testing.T) {
	keyStore := NewDiskBackEnd("~/.soc/keystore" + uuid.New().String())
	items, err := keyStore.List()
	if err != nil {
		t.Error("Failed to list keyStore ", err)
	}
	if len(items) > 0 {
		t.Error("KeyStore meant to be empty, but has items", items)
	}
	keyStore.Add(newTestKeyStoreItem())
	keyStore.Add(newTestKeyStoreItem())
	items, err = keyStore.List()
	if err != nil {
		t.Error("Failed to list keyStore ", err)
	}
	if len(items) != 2 {
		t.Error("Expected 2 items, instead", items)
	}
	keyStore.Remove(items[0])
	items, err = keyStore.List()
	if err != nil {
		t.Error("Failed to list keyStore ", err)
	}
	if len(items) != 1 {
		t.Error("Expected 1 item, instead", items)
	}
	os.RemoveAll(keyStore.path())
}

func newTestKeyStoreItem() *KeyStoreItem {
	id := uuid.NewString()
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), cryptorand.Reader)
	if err != nil {
		log.Panicln("Error generating key", err)
	}
	return NewKeyStoreItem(
		&id,
		*privateKey,
		bytesToHexStr(cryptRandBytes(20)),
		&MnemonicInfo{
			entropy:    cryptRandBytes(32),
			langLocale: "en",
			path:       "m/44'/60'/0'/0/0",
		},
		"SomeGoodPass",
		ScryptN,
		ScryptP,
	)
}
