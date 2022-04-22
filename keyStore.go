package keyStore

type KeyStore struct {
	backend BackEnd
}

func NewKeyStore(backend BackEnd) *KeyStore {
	return &KeyStore{
		backend: backend,
	}
}

func (k *KeyStore) List() []*KeyStoreItem {
	return k.backend.List()
}

func (k *KeyStore) Add(item *KeyStoreItem) {
	k.backend.Add(item)
}

func (k *KeyStore) Remove(item *KeyStoreItem) {
	k.backend.Remove(item)
}
