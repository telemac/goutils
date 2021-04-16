package ansibleutils

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// Downloads and stores http files
type HttpFiles struct {
	baseFolder      string
	downloadedFiles []string
}

func NewHttpFiles(baseFolder string) *HttpFiles {
	return &HttpFiles{baseFolder: baseFolder}
}

// GetFile downloads the file at fileUrl and stores it in baseFolder folder
// returns the full path of the file stored
func (hf *HttpFiles) GetFile(fileUrl string) (string, error) {
	// fileUrl is a normal file path, just return it
	if !(strings.HasPrefix(fileUrl, "http://") || strings.HasPrefix(fileUrl, "https://")) {
		// check if file exists
		_, err := os.Stat(fileUrl)
		return fileUrl, err
	}

	u, err := url.Parse(fileUrl)
	if err != nil {
		return fileUrl, err
	}
	fileName := path.Base(u.Path)
	destPath := path.Join(hf.baseFolder, fileName)

	err = httpGetFile(fileUrl, destPath)
	if err != nil {
		return fileUrl, err
	}

	hf.downloadedFiles = append(hf.downloadedFiles, destPath)

	return destPath, nil
}

func httpGetFile(url, destFile string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(destFile, data, 0777)
	if err != nil {
		return err
	}
	err = os.Chmod(destFile, 0777)
	if err != nil {
		return err
	}
	return nil
}
