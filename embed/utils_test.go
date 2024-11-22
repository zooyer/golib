package embed

import (
	"bytes"
	"testing"
)

func TestByteString(t *testing.T) {
	var (
		str = "Hello World"
		src = []byte(str)
	)

	if byteString(src) != str {
		t.Fatal("byteString error")
	}
}

func TestStringBytes(t *testing.T) {
	var (
		src = "Hello World"
		tar = []byte(src)
	)

	if !bytes.Equal(stringBytes(src), tar) {
		t.Fatal("stringBytes error")
	}
}

func TestMd5sum(t *testing.T) {
	var (
		data []byte
		hash = "d41d8cd98f00b204e9800998ecf8427e"
	)

	if md5sum(data) != hash {
		t.Fatal("md5sum error")
	}
}

func TestToHexEscaped(t *testing.T) {
	var (
		src = "Hello World"
		dst = `\x48\x65\x6c\x6c\x6f\x20\x57\x6f\x72\x6c\x64`
	)

	if toHexEscaped(src) != dst {
		t.Fatal("toHexEscaped error")
	}
}

func TestToUnicodeEscaped(t *testing.T) {
	var (
		src = "Hello World"
		dst = `\u0048\u0065\u006C\u006C\u006F\u0020\u0057\u006F\u0072\u006C\u0064`
	)

	if toUnicodeEscaped(src) != dst {
		t.Fatal("toUnicodeEscaped error")
	}
}
