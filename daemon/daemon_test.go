package daemon

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestDaemon(t *testing.T) {
	const (
		env  = "DAEMON_TEST"
		name = "daemon.hello"
	)

	if err := Daemon(name, env, false, false); err != nil {
		panic(err)
	}

	fmt.Println(env, ":", os.Getenv(env))
	time.Sleep(time.Second * 15)
}
