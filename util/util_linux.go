// +build !arm

package util

import (
	"os"
	"syscall"
)

func Fadvise(file *os.File, off, len int64, advice uint32) error {
	_, _, errno := syscall.Syscall6(syscall.SYS_FADVISE64, file.Fd(), uintptr(off), uintptr(len), uintptr(advice), 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
