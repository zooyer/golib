package daemon

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"syscall"
	"testing"
	"time"
)

const (
	envProcessName    = "GLIB_PROCESS_NAME"    // 原始进程名
	envProcessDaemon  = "GLIB_PROCESS_DAEMON"  // 守护进程标识
	envProcessRunning = "GLIB_PROCESS_RUNNING" // 运行标识符
)

const defaultDaemonSuffix = "(glib/daemon)"

func exit() {
	if testing.Testing() {
		syscall.Exit(0)
	}

	os.Exit(0)
}

func fork(envs []string, args ...string) (proc *os.Process, err error) {
	// 当前程序的路径
	name, err := os.Executable()
	if err != nil {
		return
	}

	var attr = os.ProcAttr{
		Dir: "",   // 继承工作目录
		Env: envs, // 继承环境变量
		Sys: nil,  // 不设置进程属性
		Files: []*os.File{
			os.Stdin,  // 标准输入
			os.Stdout, // 标准输出
			os.Stderr, // 标准错误
		}, // 继承文件描述符
	}

	return os.StartProcess(name, args, &attr)
}

func daemon(dup bool) (err error) {
	// 守护进程
	if isDaemon() {
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

	var (
		proc *os.Process
		envs = os.Environ()
		args = slices.Clone(os.Args)
	)

	// 设置进程名
	args[0] = fmt.Sprintf("%s%s", os.Args[0], defaultDaemonSuffix)
	envs = append(envs, fmt.Sprintf("%s=%s", envProcessName, os.Args[0]))

	// 设置进程标识
	envs = append(envs, fmt.Sprintf("%s=1", envProcessDaemon))

	// 启动守护进程
	if proc, err = fork(envs, args...); err != nil {
		return
	}

	// 释放子进程资源管理权（交给操作系统清理）
	if err = proc.Release(); err != nil {
		return
	}

	// 这里直接退出，子进程会启动失败，需要延迟退出
	for {
		time.Sleep(time.Second)

		if _, err = os.FindProcess(proc.Pid); err == nil {
			break
		}
	}

	exit()

	return
}

func start() (err error) {
	var (
		envs = os.Environ()
		args = slices.Clone(os.Args)
	)

	// 运行中
	envs = append(envs, fmt.Sprintf("%s=1", envProcessRunning))

	// 还原进程名
	if name := os.Getenv(envProcessName); name != "" {
		args[0] = name
	}

	// 启动子进程
	proc, err := fork(envs, args...)
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

// 孤儿进程
func isDaemon() bool {
	// 这里防止子进程启动时，父进程还没退出，ppid还是父进程的，判断守护进程失败，采用传递环境变量，可以同步取到值
	//time.Sleep(time.Second)
	//return os.Getppid() == 1

	return os.Getenv(envProcessDaemon) == "1"
}

// 运行中
func isRunning() bool {
	return os.Getenv(envProcessRunning) == "1"
}

// Daemon initializes a process to run as a daemon.
// Params:
//   - dup: Duplicate standard file descriptors.
//   - always: Specifies whether the daemon should always run.
func Daemon(dup, always bool) (err error) {
	// 1. 已运行
	if isRunning() {
		return
	}

	// 2. 守护进程

	if err = daemon(dup); err != nil {
		log.Println(err)
	}

	// 3. 守护进程运行中
	for {
		if err = start(); err != nil {
			log.Println(err)
		}

		// 正常退出
		if err == nil && !always {
			exit()
		}

		time.Sleep(time.Second)
	}
}
