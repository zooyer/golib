package xos

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func pause() {
	fd := int(os.Stdin.Fd())

	if !term.IsTerminal(fd) {
		// stdin 不是终端，也直接读取，但是必须回车提交
		_, _ = os.Stdin.Read(make([]byte, 1))
		return
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return
	}
	defer term.Restore(fd, oldState)

	_, _ = os.Stdin.Read(make([]byte, 1))

	fmt.Println()
}

func Pause(v ...any) {
	if len(v) == 0 {
		v = append(v, "请按任意键继续...")
	}

	fmt.Print(v...)
	defer fmt.Println()

	pause()
}

func PauseExit() {
	Pause("请按任意键退出程序...")
}
