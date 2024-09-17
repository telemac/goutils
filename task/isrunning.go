package task

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

// IsRunning checks if another instance of the program is running by using file locking.
func IsRunning(lockFile string) (bool, error) {
	// Open or create the lock file
	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return false, fmt.Errorf("error opening lock file: %w", err)
	}

	// Try to lock the file using platform-specific syscall
	err = unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB)
	if err != nil {
		// If we can't acquire the lock, it means another instance is running
		if err == unix.EWOULDBLOCK {
			return true, nil
		}
		return false, fmt.Errorf("error acquiring file lock: %w", err)
	}

	// The file lock is held, no other instance is running
	return false, nil
}
