package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		_ = fmt.Sprintf("Usage: %s <file> <srouce_string> <target_string>\n", os.Args[0])
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var index = bytes.Index(data, []byte(os.Args[2]))
	if index == -1 {
		fmt.Println(os.Args[2], "Not found")
		return
	}

	copy(data[index:], os.Args[3])

	if err = os.WriteFile(os.Args[1], data, 0600); err != nil {
		panic(err)
	}

	fmt.Println("Replace", index, os.Args[2], "to", os.Args[3])
}
