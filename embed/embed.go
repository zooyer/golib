package embed

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"slices"
	"sort"
	"time"
)

func getOffset(file *os.File, magic []byte) (offsets []int64, err error) {
	// 1. 获取文件大小
	var offset int64
	if offset, err = file.Seek(0, io.SeekEnd); err != nil {
		return
	}

	// 2. 重置文件游标
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return
	}

	// 3. 小于4MB直接读取查找
	if offset <= 1024*1024*4 {
		var data []byte
		if data, err = io.ReadAll(file); err != nil {
			return
		}

		offset = 0
		for {
			var index = bytes.Index(data[offset:], magic)
			if index == -1 {
				break
			}

			offsets = append(offsets, offset+int64(index))

			offset += int64(index) + int64(len(magic))
		}

		return
	}

	offset = 0
	var scanner = bufio.NewScanner(file)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if index := bytes.Index(data, magic); index >= 0 {
			offsets = append(offsets, offset+int64(index))
			advance = index + len(magic)

			return advance, data[:advance], nil
		}

		if atEOF {
			return 0, data, bufio.ErrFinalToken
		}

		return len(data), nil, nil
	})

	for scanner.Scan() {
		offset += int64(len(scanner.Bytes()))
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}

func getHeader(file *os.File, offset int64) (_ *Header, err error) {
	var h Header

	h.Offset = offset + headerSize

	if _, err = file.Seek(offset, io.SeekStart); err != nil {
		return
	}

	if err = binary.Read(file, binary.BigEndian, &h.header); err != nil {
		return
	}

	// 已初始化
	if h.IsInit() {
		return &h, nil
	}

	return &h, nil
}

func getHeaders(file *os.File, offsets []int64) (headers []Header, err error) {
	var h *Header
	for _, offset := range offsets {
		if h, err = getHeader(file, offset); err != nil {
			return
		}
		headers = append(headers, *h)
	}

	return
}

func linkHeaders(headers []Header) {
	// 按照偏移量排序
	sort.Slice(headers, func(i, j int) bool {
		return headers[i].Offset < headers[j].Offset
	})

	for i, h := range headers {
		if i >= len(headers)-1 {
			h.NextOffset = 0
			break
		}

		headers[i].NextOffset = uint32(h.Offset - headerSize)
	}
}

func checkHeaders(file *os.File, headers []Header) []Header {
	var (
		err        error
		newHeaders = make([]Header, 0, len(headers))
	)
	for _, h := range headers {
		if h.Reserve1 != 0 || h.Reserve2 != 0 {
			continue
		}

		// 校验数据容量
		if _, err = file.Seek(h.Offset+int64(h.DataCap), io.SeekStart); err != nil {
			continue
		}

		// data为空，则读取文件校验
		if _, err = file.Seek(h.Offset, io.SeekStart); err != nil {
			continue
		}

		var buf = make([]byte, h.DataLen)
		if h.DataLen > 0 {
			if _, err = io.ReadFull(file, buf); err != nil {
				continue
			}
		}

		if h.Verify(buf) == nil {
			newHeaders = append(newHeaders, h)
		}
	}

	return newHeaders
}

func initHeaders(headers []Header) []Header {
	for i, h := range headers {
		// 首次初始化时间
		if h.CreateTime == 0 {
			headers[i].CreateTime = time.Now().Unix()
		}
	}

	return headers
}

func writeHeaders(file *os.File, headers []Header) (err error) {
	var data []byte

	for _, h := range headers {
		// 计算头crc32
		h.CRC32 = 0
		if data, err = h.Encode(); err != nil {
			return
		}
		h.CRC32 = crc32.ChecksumIEEE(data)

		// 序列化头
		if data, err = h.Encode(); err != nil {
			return
		}

		// 移动游标
		if _, err = file.Seek(h.Offset-headerSize, io.SeekStart); err != nil {
			return
		}

		// 刷写落盘
		if _, err = file.Write(data); err != nil {
			return
		}
	}

	return
}

func readHeaders(file *os.File) (headers []Header, err error) {
	// 重置文件游标
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return
	}

	// 获取偏移量
	offsets, err := getOffset(file, []byte(magic))
	if err != nil {
		return
	}

	// 解析头
	if headers, err = getHeaders(file, offsets); err != nil {
		return
	}

	// 校验头，过滤掉非法头
	headers = checkHeaders(file, headers)

	// TODO 这里定义时也会调用，可能导致下面写入失败（text file busy），导致运行中断
	return

	var cloneHeaders = slices.Clone(headers)

	// 建立连接
	linkHeaders(headers)

	// 初始化
	headers = initHeaders(headers)

	// 一致则不用写入文件
	if slices.Equal(headers, cloneHeaders) {
		return
	}

	// 保存到文件
	if err = writeHeaders(file, headers); err != nil {
		return
	}

	// 同步落盘
	if err = file.Sync(); err != nil {
		return
	}

	return
}

var (
	embed  *Embed
	malloc = make(map[string]bool)
)

func init() {
	var err error

	this, err := os.Executable()
	if err != nil {
		panic(err)
	}

	if embed, err = Open(this); err != nil {
		panic(err)
	}
}

type Block struct {
	file   *os.File
	header Header
}

func (b Block) String() string {
	data, _ := json.Marshal(b.header.header)
	return string(data)
}

func (b Block) Len() uint32 {
	return b.header.DataLen
}

func (b Block) Cap() uint32 {
	return b.header.DataCap
}

func (b Block) Read(buf []byte) (n int, err error) {
	if len(buf) == 0 {
		return
	}

	if uint32(len(buf)) > b.header.DataLen {
		buf = buf[:b.header.DataLen]
	}

	if _, err = b.file.Seek(b.header.Offset, io.SeekStart); err != nil {
		return
	}

	return b.file.Read(buf)
}

func (b Block) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return
	}

	// 校验数据大小
	if uint32(len(data)) > b.header.DataCap {
		return 0, fmt.Errorf("data too large")
	}

	// 移动游标
	if _, err = b.file.Seek(b.header.Offset, io.SeekStart); err != nil {
		fmt.Println("seek error:???", err)
		return
	}

	// 写入数据
	if _, err = b.file.Write(data); err != nil {
		return
	}

	// 更新头
	b.header.DataLen = uint32(len(data))
	b.header.DataCRC32 = crc32.ChecksumIEEE(data)
	b.header.UpdateTime = time.Now().Unix()
	if b.header.CreateTime == 0 {
		b.header.CreateTime = time.Now().Unix()
	}

	// 计算头crc32
	var crcData []byte
	b.header.CRC32 = 0
	if crcData, err = b.header.Encode(); err != nil {
		return
	}
	b.header.CRC32 = crc32.ChecksumIEEE(crcData)

	// 校验头
	if err = b.header.Verify(data); err != nil {
		return
	}

	// 序列化头
	if data, err = b.header.Encode(); err != nil {
		return
	}

	// 移动游标
	if _, err = b.file.Seek(b.header.Offset-headerSize, io.SeekStart); err != nil {
		return
	}

	// 写入头
	if _, err = b.file.Write(data); err != nil {
		return
	}

	// 刷写落盘
	if err = b.file.Sync(); err != nil {
		return
	}

	return
}

type Embed struct {
	file *os.File
}

func (e *Embed) Blocks() (blocks []Block, err error) {
	headers, err := readHeaders(e.file)
	if err != nil {
		return
	}

	for _, h := range headers {
		blocks = append(blocks, Block{
			file:   e.file,
			header: h,
		})
	}

	return
}

func (e *Embed) Close() (err error) {
	return e.file.Close()
}

func Open(filename string) (embed *Embed, err error) {
	var flag = os.O_RDWR

	this, err := os.Executable()
	if err != nil {
		return
	}

	if filename == this {
		flag = os.O_RDONLY
	}

	file, err := os.OpenFile(filename, flag, 0644)
	if err != nil {
		return
	}

	return &Embed{file: file}, err
}

func Blocks() (blocks []Block, err error) {
	return embed.Blocks()
}

func Export(filename string, block Block) (err error) {
	var buf = make([]byte, block.header.DataLen)

	if _, err = block.Read(buf); err != nil {
		return
	}

	if err = os.WriteFile(filename, buf, 0644); err != nil {
		return
	}

	return
}

func Malloc(size Size) (_ *Block, err error) {
	if len(size) < headerSize {
		return nil, errors.New("invalid size")
	}

	if malloc[string(size)] {
		return nil, errors.New("block already malloced")
	}

	blocks, err := Blocks()
	if err != nil {
		return
	}

	for _, b := range blocks {
		//if data, err = b.header.Encode(); err != nil {
		//	return
		//}

		// 数据长度
		var baseSize = headerSize + b.header.DataCap
		if uint32(len(size)) < baseSize {
			continue
		}

		// 匹配头
		//if !bytes.Equal(data, []byte(size[:len(data)])) {
		//	continue
		//}

		// 匹配尾
		if uint32(len(size)) > baseSize {
			var (
				suffix = size[baseSize:]
				buffer = make([]byte, len(suffix))
			)

			if _, err = b.file.Seek(b.header.Offset+int64(b.header.DataCap), io.SeekStart); err != nil {
				return
			}

			if _, err = io.ReadFull(b.file, buffer); err != nil {
				return
			}

			if !bytes.Equal([]byte(suffix), buffer) {
				continue
			}
		}

		var (
			data = stringBytes(string(size))
			hash = md5sum(data)
		)
		malloc[hash] = true

		return &b, nil
	}

	return nil, errors.New("invalid size")
}

func MustMalloc(size Size) *Block {
	block, err := Malloc(size)
	if err != nil {
		panic(err)
	}

	return block
}

func MallocBytes(size Size) (buf []byte, err error) {
	return
}
