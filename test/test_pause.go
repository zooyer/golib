package main

import "github.com/zooyer/golib/xos"

func main() {
	defer xos.PauseExit()
	xos.Pause()
	xos.Pause("请按任意键继续。。。")
}
