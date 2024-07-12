//go:build linux || darwin

package task

import "syscall"

// Abort aborts the process by sending a SIGABRT signal
func Abort() {
	syscall.Kill(syscall.Getpid(), syscall.SIGABRT)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}

// Interrupt aborts the process by sending a SIGINT signal
func Interrupt() {
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
