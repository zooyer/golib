package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zooyer/golib/daemon"
)

var (
	daemonEnvName  = "DAEMON_TEST"
	daemonProcName = filepath.Base(os.Args[0]) + ".daemon_test"
)

func init() {
	if err := daemon.Daemon(daemonProcName, daemonEnvName, false, false); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Started.")
	fmt.Println("Please kill me to verify whether it can be restarted.")

	// 1分钟后全部退出，守护进程也退出
	for i := 0; i < 60; i++ {
		fmt.Println("This pid:", os.Getpid())
		time.Sleep(time.Second)
	}
}
