package xos

import (
	"os"
	"unsafe"
)

func SetProcessName(name string) {
	var (
		addr0 = unsafe.StringData(os.Args[0])
		argv0 = (*[1 << 30]byte)(unsafe.Pointer(addr0))[:len(os.Args[0])]
	)

	if n := copy(argv0, name); n < len(argv0) {
		argv0[n] = 0
	}

	return
}

//func SetProcessName2(name string) error {
//	bytes := append([]byte(name), 0)
//	ptr := unsafe.Pointer(&bytes[0])
//	if _, _, errno := syscall.RawSyscall6(syscall.SYS_PRCTL, syscall.PR_SET_NAME, uintptr(ptr), 0, 0, 0, 0); errno != 0 {
//		return syscall.Errno(errno)
//	}
//	return nil
//}
