package embed

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"unsafe"
)

func byteString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func stringBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func md5sum(data []byte) string {
	var hash = md5.Sum(data)

	return hex.EncodeToString(hash[:])
}

func toHexEscaped(input string) string {
	escaped := ""
	for _, b := range []byte(input) {
		escaped += fmt.Sprintf("\\x%02x", b)
	}
	return escaped
}

func toUnicodeEscaped(s string) string {
	result := ""
	for _, r := range s {
		if r <= 0xFFFF {
			result += fmt.Sprintf("\\u%04X", r)
		} else {
			result += fmt.Sprintf("\\U%08X", r)
		}
	}
	return result
}
