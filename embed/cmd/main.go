package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/zooyer/golib/embed"
)

type Command string

const (
	Show   Command = "show"
	Print  Command = "print"
	Import Command = "import"
	Export Command = "export"
	Help   Command = "help"
)

var commands = []Command{Show, Print, Import, Export, Help}

var (
	block1 = embed.MustMalloc(embed.Size1KB + "1")
	block2 = embed.MustMalloc(embed.Size1KB + "2")
)

var _, this = filepath.Split(os.Args[0])

func help(format string, v ...any) {
	fmt.Println(fmt.Sprintf(format, v...))
	fmt.Println("See 'embed help'")
	os.Exit(2)
}

func usage() {
	fmt.Printf("Usage: %s source_file <COMMAND> <BLOCK> <import_file | export_file>\n", this)
	fmt.Println()
	fmt.Println("desc...")
	fmt.Println()
	fmt.Println("Blocks(all|number):")
	fmt.Println("  all\t\tAll blocks")
	fmt.Println("  0\t\tBlock number 0")
	fmt.Println("  1\t\tBlock number 1")
	fmt.Println("  2\t\tBlock number 2")
	fmt.Println("  ...\t\tBlock number ...")
	fmt.Println()

	fmt.Println("Commands:")
	fmt.Println("  show\t\tPrints blocks info")
	fmt.Println("  import\t\tImport files into blocks")
	fmt.Println("  export\t\tExport blocks to files")
	fmt.Println("  help\tPrints this help message")
	fmt.Println()
	// embed file COMMAND BLOCK file...
}

func isAll(str string) bool {
	return str == "all"
}

func isNumber(str string) bool {
	for _, r := range str {
		if !unicode.IsNumber(r) {
			return false
		}
	}

	return true
}

func openBlocks(file string) (*embed.Embed, []embed.Block) {
	emd, err := embed.Open(file)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", file, err)
		os.Exit(1)
	}

	blocks, err := emd.Blocks()
	if err != nil {
		fmt.Printf("Error getting blocks: %s\n", err)
		os.Exit(1)
	}

	return emd, blocks
}

func showID(file string, id int) {
	emd, blocks := openBlocks(file)

	if id >= len(blocks) {
		fmt.Printf("Block %d not found\n", id)
		os.Exit(1)
	}

	if err := emd.Close(); err != nil {
		fmt.Printf("Error closing file %s: %s\n", file, err)
		os.Exit(1)
	}

	var block = blocks[id]
	fmt.Printf("Block %d:\n", id)
	fmt.Println(block.String())
	fmt.Println()
}

func showAll(file string) {
	emd, blocks := openBlocks(file)

	if err := emd.Close(); err != nil {
		fmt.Printf("Error closing file %s: %s\n", file, err)
		os.Exit(1)
	}

	for id, block := range blocks {
		fmt.Printf("Block %d:\n", id)
		fmt.Println(block.String())
		fmt.Println()
	}
}

func importID(file string, id int, filename string) {
	emd, blocks := openBlocks(file)

	if id >= len(blocks) {
		fmt.Printf("Block %d not found\n", id)
		os.Exit(1)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	if _, err = blocks[id].Write(data); err != nil {
		fmt.Printf("Error writing block %d: %s\n", id, err)
		os.Exit(1)
	}

	if err = emd.Close(); err != nil {
		fmt.Printf("Error closing block %d: %s\n", id, err)
		os.Exit(1)
	}

	fmt.Printf("Import block %d successful.\n", id)
}

func importAll(file string, filenames ...string) {
	emd, blocks := openBlocks(file)

	if len(blocks) != len(filenames) {
		fmt.Printf("Block count mismatch: %d != %d\n", len(blocks), len(filenames))
		os.Exit(1)
	}

	for id := range blocks {
		data, err := os.ReadFile(filenames[id])
		if err != nil {
			fmt.Printf("Error reading file %s: %s\n", filenames[id], err)
			os.Exit(1)
		}

		if _, err = blocks[id].Write(data); err != nil {
			fmt.Printf("Error writing block %d: %s\n", id, err)
			os.Exit(1)
		}
	}

	if err := emd.Close(); err != nil {
		fmt.Printf("Error closing blocks: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Import blocks successful.\n")
}

func exportID(file string, id int, filename string) {
	emd, blocks := openBlocks(file)

	if id >= len(blocks) {
		fmt.Printf("Block %d not found\n", id)
		os.Exit(1)
	}

	var (
		block = blocks[id]
		buf   = make([]byte, block.Len())
	)

	if _, err := block.Read(buf); err != nil {
		fmt.Printf("Error reading block %d: %s\n", id, err)
		os.Exit(1)
	}

	if err := emd.Close(); err != nil {
		fmt.Printf("Error closing block %d: %s\n", id, err)
		os.Exit(1)
	}

	if err := os.WriteFile(filename, buf, 0644); err != nil {
		fmt.Printf("Error writing block %d: %s\n", id, err)
		os.Exit(1)
	}

	fmt.Printf("Export block %d successful.\n", id)
}

func exportAll(file, filename string) {
	emd, blocks := openBlocks(file)

	var format = fmt.Sprintf("%%0%dd", len(strconv.Itoa(len(blocks))))
	for i, block := range blocks {
		var buf = make([]byte, block.Len())
		if _, err := block.Read(buf); err != nil {
			fmt.Printf("Error reading block %d: %s\n", i, err)
			os.Exit(1)
		}

		if err := os.WriteFile(fmt.Sprintf(filename+"."+format, i), buf, 0644); err != nil {
			fmt.Printf("Error writing block %d: %s\n", i, err)
			os.Exit(1)
		}
	}

	if err := emd.Close(); err != nil {
		fmt.Printf("Error closing blocks: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Export blocks successful.\n")
}

func printID(file string, id int) {
	emd, blocks := openBlocks(file)

	if id >= len(blocks) {
		fmt.Printf("Block %d not found\n", id)
		os.Exit(1)
	}

	var (
		block = blocks[id]
		buf   = make([]byte, block.Len())
	)

	if _, err := block.Read(buf); err != nil {
		fmt.Printf("Error reading block %d: %s\n", id, err)
		os.Exit(1)
	}

	if err := emd.Close(); err != nil {
		fmt.Printf("Error closing block %d: %s\n", id, err)
		os.Exit(1)
	}

	fmt.Printf("Block %d:\n", id)
	fmt.Printf("%s\n", buf)
}

func printAll(file string) {
	emd, blocks := openBlocks(file)

	for id, block := range blocks {
		var buf = make([]byte, block.Len())
		if _, err := block.Read(buf); err != nil {
			fmt.Printf("Error reading block %d: %s\n", block, err)
			os.Exit(1)
		}

		fmt.Printf("Block %d:\n", id)
		fmt.Printf("%s\n", buf)
	}

	if err := emd.Close(); err != nil {
		fmt.Printf("Error closing blocks: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var (
		args = os.Args
		file = args[1]
	)

	// 帮助
	if Command(file) == Help {
		usage()
		os.Exit(0)
	}

	// 解析命令
	var command = Show
	if args = args[2:]; len(args) > 0 {
		command = Command(args[0])
		args = args[1:]
	}

	// 校验命令
	if !slices.Contains(commands, command) {
		help("%s: '%s' is not a embed command.", this, command)
	}

	// block
	var block string
	if len(args) > 0 {
		block = args[0]
		args = args[1:]
	}

	if command == Show || command == Print && block == "" {
		block = "all"
	}

	// 校验block
	if !isAll(block) && !isNumber(block) {
		help("%s: '%s' is not a block id.", this, block)
	}

	// 获取block id
	var id int
	if !isAll(block) {
		var err error
		if id, err = strconv.Atoi(block); err != nil {
			help("%s: '%s' is not a block id.", this, block)
		}
	}

	var files = args

	// 校验目标文件
	if command != Show && command != Print && command != Help && len(files) == 0 {
		var params = strings.Join(os.Args[1:], " ")
		help("%s: %s %s <%s_file>, the %s file is missing.", this, this, params, command, command)
	}

	// 执行命令
	switch command {
	case Show:
		if isAll(block) {
			showAll(file)
		} else {
			showID(file, id)
		}
	case Print:
		if isAll(block) {
			printAll(file)
		} else {
			printID(file, id)
		}
	case Import:
		if isAll(block) {
			importAll(file, files...)
		} else {
			importID(file, id, files[0])
		}
	case Export:
		if isAll(block) {
			exportAll(file, files[0])
		} else {
			exportID(file, id, files[0])
		}
	case Help:
		usage()
		os.Exit(0)
	}
}
