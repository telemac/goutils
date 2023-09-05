package task

import "syscall"

// Abort aborts the process by sending a SIGABRT signal
func Abort() {
	syscall.Kill(syscall.Getpid(), syscall.SIGABRT)
}
