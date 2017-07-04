package odin

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(value []byte) []byte {
	var m = md5.New()
	m.Write(value)
	return m.Sum(nil)
}

func MD5String(value string) string {
	return hex.EncodeToString(MD5([]byte(value)))
}
