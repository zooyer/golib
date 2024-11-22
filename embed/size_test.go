package embed

import (
	"fmt"
	"testing"
)

func TestGenSize(t *testing.T) {
	const (
		kb = 1024
		mb = 1024 * 1024
	)
	var sizes = []uint32{
		1 * kb, 2 * kb, 4 * kb, 8 * kb, 16 * kb, 32 * kb, 64 * kb, 128 * kb, 512 * kb,
		1 * mb, 2 * mb, 4 * mb, 8 * mb, 16 * mb, 32 * mb,
	}
	for _, size := range sizes {
		fmt.Println(genSizeDefine("Size", size))
	}
}
