package util

import (
	"os"
	"syscall"
)

func Fadvise(file *os.File, off, len int64, advice uint32) error {
	_, _, errno := syscall.Syscall6(syscall.SYS_ARM_FADVISE64_64, file.Fd(), uintptr(advice), uintptr(off), uintptr(len), 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
