package files

import "os"

// FileExists returns true if the file exists, false if not
func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
