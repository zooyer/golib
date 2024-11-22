package embed

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
)

const (
	magic           = "\uEEEE\u0047\u004F\uEEEE"
	emptyDataLen    = "\x00\x00\x00\x00"
	emptyHeaderRest = "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"
	headerSize      = 52
)

type header struct {
	Magic      uint64 // 魔术标志
	CRC32      uint32 // CRC32校验
	DataLen    uint32 // 数据大小
	DataCap    uint32 // 数据容量
	DataCRC32  uint32 // 数据CRC32
	NextOffset uint32 // 下个数据偏移
	CreateTime int64  // 首次写入时间
	UpdateTime int64  // 最后写入时间
	Reserve1   uint32 // 保留字段1
	Reserve2   uint32 // 保留字段2
}

var emptyHeader = header{
	Magic: binary.BigEndian.Uint64([]byte(magic)),
	CRC32: 0xf151cd41,
}

func (h *header) String() string {
	var buf = make([]byte, headerSize)

	n, err := binary.Encode(buf, binary.BigEndian, h)
	if err != nil {
		return fmt.Sprintf("<err>: %s", err)
	}

	return string(buf[:n])
}

func (h *header) IsInit() bool {
	return h.CreateTime != 0
}

func (h *header) Encode() (data []byte, err error) {
	var buf = make([]byte, headerSize)

	n, err := binary.Encode(buf, binary.BigEndian, h)
	if err != nil {
		return
	}

	return buf[:n], nil
}

type Header struct {
	header
	Offset int64 // 数据偏移
}

func (h *Header) Verify(data []byte) (err error) {
	// 校验magic
	if h.Magic != emptyHeader.Magic {
		return errors.New("invalid header magic")
	}

	// 计算头crc32
	var (
		crc32Data   []byte
		crc32Header = *h
	)
	crc32Header.CRC32 = 0
	if crc32Data, err = crc32Header.Encode(); err != nil {
		return
	}

	// 校验头crc32
	if h.CRC32 != crc32.ChecksumIEEE(crc32Data) {
		return errors.New("invalid header crc32")
	}

	// 校验数据大小
	if h.DataLen > h.DataCap || len(data) != int(h.DataLen) {
		return errors.New("invalid data length")
	}

	// 校验数据crc32
	if h.DataCRC32 != crc32.ChecksumIEEE(data) {
		return errors.New("invalid data checksum")
	}

	return
}
