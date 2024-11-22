package xos

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestSetProcessName(t *testing.T) {
	var tests = []string{"test", "Hello", "abc", "666"}

	for i, test := range tests {
		SetProcessName(test)

		var (
			pid  = os.Getpid()
			cmd  = "ps"
			args = []string{"-o", "command", "-p", strconv.Itoa(pid)}
		)
		output, err := exec.Command(cmd, args...).Output()
		if err != nil {
			t.Fatal(i, err)
		}

		var lines = bytes.Split(output, []byte("\n"))
		if len(lines) < 2 {
			t.Fatal(i, "not found pid:", pid)
		}

		var cmdline = strings.Split(string(lines[1]), " ")
		if cmdline[0] != test {
			t.Fatal(i, "not match:", test, cmdline[0])
		}
	}
}
