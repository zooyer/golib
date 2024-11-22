package embed

import (
	"testing"
	"time"
)

func TestHeader_Encode(t *testing.T) {
	var h = header{
		Magic:      emptyHeader.Magic,
		CRC32:      emptyHeader.CRC32,
		DataLen:    0,
		DataCap:    0,
		DataCRC32:  0,
		NextOffset: 0,
		CreateTime: time.Now().Unix(),
		UpdateTime: 0,
		Reserve1:   0,
		Reserve2:   0,
	}

	if !h.IsInit() {
		t.Fatal("IsInit fail")
	}

	data, err := h.Encode()
	if err != nil {
		t.Fatal(err)
	}

	if len(data) != headerSize {
		t.Fatal("header encode fail")
	}
}

func TestHeader_Decode(t *testing.T) {
	var h = Header{
		header: emptyHeader,
		Offset: 0,
	}

	if err := h.Verify(nil); err != nil {
		t.Fatal(err)
	}
}
