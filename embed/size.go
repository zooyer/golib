package embed

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

type Size string

const (
	size128Byte Size = "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	size512Byte      = size128Byte + size128Byte + size128Byte + size128Byte
)

const (
	size1KB   = size512Byte + size512Byte
	size2KB   = size1KB + size1KB
	size4KB   = size2KB + size2KB
	size8KB   = size4KB + size4KB
	size16KB  = size8KB + size8KB
	size32KB  = size16KB + size16KB
	size64KB  = size32KB + size32KB
	size128KB = size64KB + size64KB
	size256KB = size128KB + size128KB
	size512KB = size256KB + size256KB

	size1MB  = size512KB + size512KB
	size2MB  = size1MB + size1MB
	size4MB  = size2MB + size2MB
	size8MB  = size4MB + size4MB
	size16MB = size8MB + size8MB
	size32MB = size16MB + size16MB
)

const (
	Size1KB   = magic + "\x63\x58\x89\x99" + emptyDataLen + "\x00\x00\x04\x00" + emptyHeaderRest + size1KB
	Size2KB   = magic + "\x0e\x32\x42\xb0" + emptyDataLen + "\x00\x00\x08\x00" + emptyHeaderRest + size2KB
	Size4KB   = magic + "\xd4\xe7\xd4\xe2" + emptyDataLen + "\x00\x00\x10\x00" + emptyHeaderRest + size4KB
	Size8KB   = magic + "\xba\x3d\xfe\x07" + emptyDataLen + "\x00\x00\x20\x00" + emptyHeaderRest + size8KB
	Size16KB  = magic + "\x67\x89\xab\xcd" + emptyDataLen + "\x00\x00\x40\x00" + emptyHeaderRest + size16KB
	Size32KB  = magic + "\x07\x90\x06\x18" + emptyDataLen + "\x00\x00\x80\x00" + emptyHeaderRest + size32KB
	Size64KB  = magic + "\x3e\xcf\xda\x89" + emptyDataLen + "\x00\x01\x00\x00" + emptyHeaderRest + size64KB
	Size128KB = magic + "\xb5\x1c\xe4\x90" + emptyDataLen + "\x00\x02\x00\x00" + emptyHeaderRest + size128KB
	Size512KB = magic + "\x3b\x14\x6c\x44" + emptyDataLen + "\x00\x08\x00\x00" + emptyHeaderRest + size512KB

	Size1MB  = magic + "\xbe\xab\x89\x0a" + emptyDataLen + "\x00\x10\x00\x00" + emptyHeaderRest + size1MB
	Size2MB  = magic + "\x6e\xa5\x45\xd7" + emptyDataLen + "\x00\x20\x00\x00" + emptyHeaderRest + size2MB
	Size4MB  = magic + "\x15\xc9\xda\x2c" + emptyDataLen + "\x00\x40\x00\x00" + emptyHeaderRest + size4MB
	Size8MB  = magic + "\xe3\x10\xe5\xda" + emptyDataLen + "\x00\x80\x00\x00" + emptyHeaderRest + size8MB
	Size16MB = magic + "\x64\x21\x19\xd4" + emptyDataLen + "\x01\x00\x00\x00" + emptyHeaderRest + size16MB
	Size32MB = magic + "\x00\xc1\x62\x2a" + emptyDataLen + "\x02\x00\x00\x00" + emptyHeaderRest + size32MB
)

func genSize(size uint32) string {
	var buf = make([]byte, 4)

	binary.BigEndian.PutUint32(buf, size)

	return string(buf)
}

func genSizeDefine(prefix string, size uint32) string {
	var (
		units   = []string{"B", "KB", "MB", "GB", "TB", "PB"}
		hexSize = toHexEscaped(genSize(size))
		header  = header{
			Magic:      emptyHeader.Magic,
			CRC32:      0,
			DataLen:    0,
			DataCap:    size,
			DataCRC32:  0,
			NextOffset: 0,
			CreateTime: 0,
			UpdateTime: 0,
			Reserve1:   0,
			Reserve2:   0,
		}
	)

	var index int
	for size >= 1024 && index < len(units)-1 {
		size /= 1024 // 不考虑小数部分
		index++
	}

	var unit = units[index]

	data, err := header.Encode()
	if err != nil {
		panic(err)
	}

	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, crc32.ChecksumIEEE(data))
	var headerCRC32 = toHexEscaped(string(buf))

	return fmt.Sprintf("%s%d%s = magic + \"%s\" + emptyDataLen + \"%s\" + emptyHeaderRest + size%d%s", prefix, size, unit, headerCRC32, hexSize, size, unit)
}
