package odin

import (
	"crypto/md5"
	"encoding/hex"
)

func md5String(value string) string {
	var m = md5.New()
	m.Write([]byte(value))
	var bs = m.Sum(nil)
	return hex.EncodeToString(bs)
}
