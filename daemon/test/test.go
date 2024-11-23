package main

import (
	"fmt"
	"os"
	"time"

	"github.com/zooyer/golib/daemon"
)

func init() {
	if err := daemon.Daemon(true, true); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Started.")
	fmt.Println("Please kill me to verify whether it can be restarted.")

	// 1分钟后全部退出，守护进程也退出
	for i := 0; i < 30; i++ {
		fmt.Println("This pid:", os.Getpid(), "num:", i)
		time.Sleep(time.Second)
	}
}
