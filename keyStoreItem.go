package keyStore

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type KeyStoreItem struct {
	Uuid string `json:"uuid"`
	// UTC--<created_at UTC ISO8601>--uuid.json
	FileName string `json:"file_name"`
}

func NewKeyStoreItem() *KeyStoreItem {
	uuid := uuid.NewString()
	return &KeyStoreItem{
		Uuid: uuid,
		FileName: fmt.Sprintf(
			"UTC--%s--%s.json",
			toISO8601(time.Now().UTC()),
			uuid,
		),
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
