package updater

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/kardianos/osext"
	log "github.com/sirupsen/logrus"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// SelfUpdater allows a binary to update itself
type SelfUpdater struct {
	baseUrl             string // url base for remote artifacts (ex : "https://update.plugis.com/")
	remoteUrl           string // url of the remote binary
	localBinaryFilePath string
}

// NewSelfUpdater creates a SelfUpdatr
// baseUrl sample : https://update.plugis.com/ (must end with /)
func NewSelfUpdater(baseUrl, remoteUrl string) (*SelfUpdater, error) {
	localBinaryFilePath, err := osext.Executable()
	if err != nil {
		return nil, err
	}

	if remoteUrl == "" {
		exeName := filepath.Base(localBinaryFilePath) // get binary name
		uri := fmt.Sprintf("%s/%s/%s", runtime.GOOS, runtime.GOARCH, exeName)
		remoteUrl = baseUrl + uri
	}
	return &SelfUpdater{baseUrl: baseUrl, remoteUrl: remoteUrl, localBinaryFilePath: localBinaryFilePath}, nil
}

// regular expression to extract a md5 hash
var MD5RE = regexp.MustCompile("([0-9a-fA-F]{32})")

// GetRemoteMD5 gets the md5 contained in the remote file with .md5 extension
func (su *SelfUpdater) GetRemoteMD5() (string, error) {
	line, err := get(su.remoteUrl + ".md5")
	if err != nil {
		return "", fmt.Errorf("%s : %w", su.remoteUrl+".md5", err)
	}
	md5 := MD5RE.FindString(string(line))
	if md5 == "" {
		return "", fmt.Errorf("no md5 hash in file %s", su.remoteUrl+".md5")
	}
	return md5, nil
}

// GetLocalMD5 computes the local md5 of the current binary
func (su *SelfUpdater) GetLocalMD5() (string, error) {
	return hash_file_md5(su.localBinaryFilePath)
}

// NeedsUpdate returns true if the md5 of the remote binary is different from the md5 of the local binary
func (su *SelfUpdater) NeedsUpdate() (bool, error) {
	remoteMD5, err := su.GetRemoteMD5()
	if err != nil {
		return false, fmt.Errorf("get remote md5 : %w", err)
	}
	localMD5, err := su.GetLocalMD5()
	if err != nil {
		return false, fmt.Errorf("get local md5 : %w", err)
	}
	slog.Info("compare md5 for update", "remote", remoteMD5, "local", localMD5)
	return !strings.EqualFold(localMD5, remoteMD5), nil
	// return strings.ToUpper(localMD5) != strings.ToUpper(remoteMD5), nil
}

// SelfUpdate updates the current binary
// if force==true, don't check md5 and update anyway
func (su *SelfUpdater) SelfUpdate(force bool) (bool, error) {
	var needsUpdate bool
	var err error
	if !force {
		needsUpdate, err = su.NeedsUpdate()
		if err != nil {
			return false, err
		}
	}
	if needsUpdate || force {
		err = downloadAndReplace(su.remoteUrl)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// SelfUpdate downloads the executable and replaces the current running binary
// baseUrl sample : https://update.domain.com/
func SelfUpdate(baseUrl string) (bool, error) {
	localBinaryFilePath, err := osext.Executable()

	if err != nil {
		return false, err
	}

	exeName := filepath.Base(localBinaryFilePath)
	url := fmt.Sprintf("%s/%s/%s", runtime.GOOS, runtime.GOARCH, exeName)
	// compare remote and local md5
	remoteMd5, err := get(baseUrl + url + ".md5")
	if err == nil { // if md5 file present on remote server

		localMd5, err := hash_file_md5(localBinaryFilePath)
		if err != nil {
			return false, err
		}
		fmt.Printf("local=%s, remote=%s\n", localMd5, remoteMd5)
		if strings.Contains(string(remoteMd5), localMd5) {
			fmt.Println("no update needed")
			return false, nil
		}

	} else {
		log.WithError(err).WithField("url", baseUrl+url+".md5").Warn("download md5")
	}
	err = downloadAndReplace(baseUrl + url)
	if err != nil {
		return false, err
	}
	return true, nil
}

func get(url string) ([]byte, error) {
	var body []byte
	resp, err := http.Get(url)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return body, errors.New(resp.Status)
	}
	return body, err
}

func hash_file_md5(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}

func downloadAndReplace(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = apply(resp.Body)
	if err != nil {
		// TODO : error handling
	}
	return err
}

func apply(update io.Reader) error {

	// get target path
	var err error
	var newBytes []byte
	if newBytes, err = io.ReadAll(update); err != nil {
		return err
	}

	//Folder, _ := os.Getwd()
	targetPath, _ := osext.Executable()

	// get the directory the executable exists in
	updateDir := filepath.Dir(targetPath)
	filename := filepath.Base(targetPath)

	// Copy the contents of newbinary to a new executable file
	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0775)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = io.Copy(fp, bytes.NewReader(newBytes))
	if err != nil {
		return err
	}

	err = fp.Sync()
	if err != nil {
		return err
	}

	// if we don't call fp.Close(), windows won't let us move the new executable
	// because the file will still be "in use"
	fp.Close()

	// Check if file can run by getting version
	_, err = exec.Command(newPath, "version").Output()
	//log.Printf("UPDATE check new executable %v %v", err, out)
	if err != nil {
		_ = os.Remove(newPath)
		return err
	}

	// this is where we'll move the executable to so that we can swap in the updated replacement
	oldPath := filepath.Join(updateDir, fmt.Sprintf(".%s.old", filename))
	// delete any existing old exec file - this is necessary on Windows for two reasons:
	// 1. after a successful update, Windows can't remove the .old file because the process is still running
	// 2. windows rename operations fail if the destination file already exists
	_ = os.Remove(oldPath)

	// move the existing executable to a new file in the same directory
	err = os.Rename(targetPath, oldPath)
	if err != nil {
		return err
	}

	// move the new exectuable in to become the new program
	err = os.Rename(newPath, targetPath)

	if err != nil {
		// move unsuccessful
		//
		// The filesystem is now in a bad state. We have successfully
		// moved the existing binary to a new location, but we couldn't move the new
		// binary to take its place. That means there is no file where the current executable binary
		// used to be!
		// Try to rollback by restoring the old binary to its original path.
		rerr := os.Rename(oldPath, targetPath)
		if rerr != nil {
			return rerr
		}

		return err
	}

	return nil
}
