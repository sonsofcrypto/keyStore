// SPDX-License-Identifier: MIT

// Package keyStore stores private keys according to Web3 Secret Storage
// Definition. Extends format to optionally store mnemonic similarly to
// Ethers.js. Supports storing of `KeyStoreItem`s to variety `Backend`s
// (eg disk, system keychain).
package keyStore

type KeyStore struct {
	backends       []Backend
	itemsByBackend map[Backend][]*KeyStoreItem
}

func NewKeyStore(backends []Backend) *KeyStore {
	return &KeyStore{
		backends:       backends,
		itemsByBackend: make(map[Backend][]*KeyStoreItem),
	}
}

func (k *KeyStore) List() ([]*KeyStoreItem, error) {
	var allItems = make([]*KeyStoreItem, 0)
	var anyErr error = nil
	for _, backend := range k.backends {
		items, err := backend.List()
		if err != nil {
			anyErr = err
		}
		allItems = append(allItems, items...)
	}
	return allItems, anyErr
}

func (k *KeyStore) Add(item *KeyStoreItem, backend Backend) error {
	return backend.Add(item)
}

func (k *KeyStore) Remove(item *KeyStoreItem, backend Backend) error {
	return backend.Remove(item)
}
