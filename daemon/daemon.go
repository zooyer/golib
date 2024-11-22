package daemon

import (
	"errors"
	"log"
	"os"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/zooyer/golib/xos"
)

const envOrphan = "GOLIB_ORPHAN_PROCESS"

var args = slices.Clone(os.Args)

func init() {
	for i, arg := range args {
		args[i] = strings.Clone(arg)
	}
}

func exit() {
	if testing.Testing() {
		syscall.Exit(0)
	}

	os.Exit(0)
}

// 孤儿进程
func isOrphan() bool {
	return os.Getppid() == 1 || os.Getenv(envOrphan) == "true"
}

// 守护进程(非linux守护进程，这里指本代码的守护进程)
func isDaemon(env string) bool {
	return os.Getenv(env) == ""
}

func fork() (proc *os.Process, err error) {
	// 当前程序的路径
	name, err := os.Executable()
	if err != nil {
		return
	}

	var (
		envs = os.Environ()
		attr = os.ProcAttr{
			Dir: "",   // 继承工作目录
			Env: envs, // 继承环境变量
			Sys: nil,  // 不设置进程属性
			Files: []*os.File{
				os.Stdin,  // 标准输入
				os.Stdout, // 标准输出
				os.Stderr, // 标准错误
			}, // 继承文件描述符
		}
	)

	return os.StartProcess(name, args, &attr)
}

func orphan(dup bool) (err error) {
	// TODO 这里防止子进程启动时，父进程还没退出，判断孤儿进程失败
	//time.Sleep(time.Second)

	// 孤儿进程当做linux守护进程
	if isOrphan() {
		if err = os.Unsetenv(envOrphan); err != nil {
			return
		}

		// 创建新会话
		if _, err = syscall.Setsid(); err != nil {
			return
		}

		// 改变工作目录
		if err = os.Chdir("/"); err != nil {
			return
		}

		// 改变文件访问权限掩码
		syscall.Umask(0022)

		if dup {
			// 关闭stdin、stdout、stderr
			if err = os.Stdin.Close(); err != nil {
				return
			}

			if err = os.Stdout.Close(); err != nil {
				return
			}

			if err = os.Stderr.Close(); err != nil {
				return
			}

			// 打开/dev/null
			var null *os.File
			if null, err = os.OpenFile(os.DevNull, os.O_RDWR, 0); err != nil {
				return
			}

			// 重定向stdin、stdout、stderr到/dev/null
			if err = syscall.Dup2(int(null.Fd()), int(os.Stdin.Fd())); err != nil {
				return
			}

			if err = syscall.Dup2(int(null.Fd()), int(os.Stdout.Fd())); err != nil {
				return
			}

			if err = syscall.Dup2(int(null.Fd()), int(os.Stderr.Fd())); err != nil {
				return
			}

			os.Stdin = null
			os.Stdout = null
			os.Stderr = null
		}

		return
	}

	// 设置孤儿进程环境变量
	if err = os.Setenv(envOrphan, "true"); err != nil {
		return
	}

	var proc *os.Process
	if proc, err = fork(); err != nil {
		return
	}

	// TODO 这里启动子进程并释放后退出，不是原子操作
	// TODO 可能子进程启动快，导致子进程判断还不是孤儿进程，导致递归fork
	// TODO 可以在判断前isOrphan增加延时解决
	// TODO 原生fork不会出现这个问题
	if err = proc.Release(); err != nil {
		return
	}

	// TODO 这里直接退出，子进程会启动失败，需要延迟退出
	time.Sleep(time.Second * 10)

	exit()

	return
}

func rename(name string) {
	xos.SetProcessName(name)
}

func start(env string) (err error) {
	// 设置非守护进程环境变量
	if err = os.Setenv(env, "true"); err != nil {
		return
	}

	// 启动子进程
	proc, err := fork()
	if err != nil {
		return
	}

	// 等待子进程退出，释放子进程资源
	state, err := proc.Wait()
	if err != nil {
		return
	}

	state.ExitCode()
	if !state.Success() {
		return errors.New(state.String())
	}

	return
}

func daemon(name, env string, dup, always bool) (err error) {
	// 1. 非守护进程
	if !isDaemon(env) {
		return
	}

	// 2. 启动孤儿进程
	if err = orphan(dup); err != nil {
		log.Println(err)
	}

	// 3. 修改进程名
	rename(name)

	// 4. 启动守护进程
	for {
		if err = start(env); err != nil {
			log.Println(err)
		}

		// 正常退出
		if err == nil && !always {
			exit()
		}

		time.Sleep(time.Second)
	}
}

// Daemon initializes a process to run as a daemon.
// Parameters:
//   - name: The name of the daemon process.
//   - env: The environment configuration for the daemon process.
//   - dup: Indicates whether file descriptors should be duplicated.
//   - always: Specifies whether the daemon should always run.
//     -- true: The child process will be restarted when it exits, regardless of the exit status code.
//     -- false: If the child process exits with a status code other than 0, restart.
func Daemon(name, env string, dup, always bool) (err error) {
	return daemon(name, env, dup, always)
}
