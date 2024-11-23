package daemon

import (
	"testing"
	"time"
)

func TestDaemon(t *testing.T) {
	if err := Daemon(false, false); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 15)
}
