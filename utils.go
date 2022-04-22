// SPDX-License-Identifier: MIT

package keyStore

import (
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"log"
	"math/big"
	"time"
)

const (
	// number of bits in big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in big.Word
	wordBytes = wordBits / 8
)

func keccak256(data ...[]byte) []byte {
	return crypto.Keccak256(data...)
}

func cryptRandBytes(n int) []byte {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(cryptorand.Reader, bytes); err != nil {
		log.Panicln("Failed to generate random bytes", err)
	}
	return bytes
}

func bytesToHexStr(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func aesCTRXOR(key, data, iv []byte) ([]byte, error) {
	output := make([]byte, len(data))
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cipher.NewCTR(aesBlock, iv).XORKeyStream(output, data)
	return output, nil
}

func paddedBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	bytes := make([]byte, n)
	readBits(bigint, bytes)
	return bytes
}

func readBits(bigint *big.Int, buf []byte) {
	i := len(buf)
	for _, d := range bigint.Bits() {
		for j := 0; j < wordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}

func toISO8601(t time.Time) string {
	name, offset := t.Zone()
	tz := "Z"
	if name != "UTC" {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}
	return fmt.Sprintf(
		"%04d-%02d-%02dT%02d-%02d-%02d.%09d%s",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
		tz,
	)
}
