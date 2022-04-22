// SPDX-License-Identifier: MIT

package keyStore

import (
	"crypto/aes"
	"crypto/ecdsa"
	"fmt"
	uuidPackage "github.com/google/uuid"
	"golang.org/x/crypto/scrypt"
	"log"
	"time"
)

const (
	headerKDF        = "scrypt"
	cipherAes128Crt  = "aes-128-ctr"
	ScryptN          = 1 << 18
	ScryptP          = 1
	scryptR          = 8
	scryptDKLen      = 32
	privateKeyMinLen = 32
	gethVersion      = 3
	socVersion       = 1
)

type KeyStoreItem struct {
	Address  string            `json:"address"`
	Crypto   Web3SecretStorage `json:"crypto"`
	Uuid     string            `json:"id"`
	Version  int               `json:"version"`
	Mnemonic *Mnemonic         `json:"x-soc-mnemonic-entropy-crypto"`
	FileName string            `json:"x-sos-fileName"`
}

type MnemonicInfo struct {
	entropy    []byte
	langLocale string
	path       string
}

func NewKeyStoreItem(
	uuid *string,
	privateKey ecdsa.PrivateKey,
	mnemonicData *MnemonicInfo,
	password string,
	scryptN int,
	scryptP int,
) *KeyStoreItem {
	if uuid == nil {
		uuidStr := uuidPackage.NewString()
		uuid = &uuidStr
	}
	fn := fmt.Sprintf("UTC--%s--%s.json", toISO8601(time.Now().UTC()), *uuid)
	bytes := paddedBytes(privateKey.D, privateKeyMinLen)
	return &KeyStoreItem{
		Address:  "",
		Crypto:   *NewWeb3SecretStorage(bytes, []byte(password), scryptN, scryptP),
		Uuid:     *uuid,
		Version:  gethVersion,
		Mnemonic: NewMnemonic(mnemonicData, []byte(password), scryptN, scryptP),
		FileName: fn,
	}
}

type Web3SecretStorage struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams CipherParams           `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type CipherParams struct {
	IV string `json:"iv"`
}

func NewWeb3SecretStorage(data, pswd []byte, n, p int) *Web3SecretStorage {
	salt := cryptRandBytes(scryptDKLen)
	dk, err := scrypt.Key(pswd, salt, n, scryptR, p, scryptDKLen)
	if err != nil {
		log.Panicln("Failed to derive key from password", err)
	}
	encryptKey := dk[:16]
	iv := cryptRandBytes(aes.BlockSize)
	cipherText, err := aesCTRXOR(encryptKey, data, iv)
	if err != nil {
		log.Panicln("Failed to encrypt data", err)
	}
	return &Web3SecretStorage{
		Cipher:       cipherAes128Crt,
		CipherText:   bytesToHexStr(cipherText),
		CipherParams: CipherParams{IV: bytesToHexStr(iv)},
		KDF:          headerKDF,
		KDFParams: map[string]interface{}{
			"n":     n,
			"r":     scryptR,
			"p":     p,
			"dklen": scryptDKLen,
			"salt":  bytesToHexStr(salt),
		},
		MAC: bytesToHexStr(keccak256(dk[16:32], cipherText)),
	}
}

type Mnemonic struct {
	Crypto     Web3SecretStorage `json:"crypto"`
	LangLocale string            `json:"langLocale"`
	Path       string            `json:"path"`
	Version    int               `json:"version"`
}

func NewMnemonic(data *MnemonicInfo, pswd []byte, n, p int) *Mnemonic {
	if data == nil {
		return nil
	}
	if len(data.entropy) < privateKeyMinLen {
		log.Panicln("Invalid entropy len", len(data.entropy))
	}
	return &Mnemonic{
		Crypto:     *NewWeb3SecretStorage(data.entropy, []byte(pswd), n, p),
		LangLocale: data.langLocale,
		Path:       data.path,
		Version:    socVersion,
	}
}
