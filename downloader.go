package downloader

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"time"
)

var ErrNotFound = errors.New("file not found")

type Downloader struct {
	url     string
	timeout time.Duration
	retries int
}

func New(url string) *Downloader {
	return &Downloader{url, 5 * time.Second, 3}
}

func (self *Downloader) Timeout(timeout time.Duration) *Downloader {
	if timeout > 0 {
		self.timeout = timeout
	}

	return self
}

func (self *Downloader) Retries(i int) *Downloader {
	if i > 1 {
		self.retries = i
	}

	return self
}

func (self *Downloader) SaveToTempFile() (string, error) {
	file, err := self.download()
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func (self *Downloader) SaveToFile(filename string) error {
	file, err := self.download()
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(filename), 0755)
	if err != nil {
		return err
	}

	return os.Rename(file.Name(), filename)
}

func (self *Downloader) download() (*os.File, error) {
	var lerr error
	var file *os.File

	for attempt := 0; attempt < self.retries; attempt++ {
		var resp *http.Response

		if attempt > 0 {
			time.Sleep(time.Duration(250*int64(math.Exp2(float64(attempt-1)))) * time.Millisecond)
		}

		resp, lerr = http.Get(self.url)
		if lerr != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusNotFound {
				return nil, ErrNotFound
			}

			lerr = fmt.Errorf("wrong status code: %d", resp.StatusCode)
			continue
		}

		file, lerr = ioutil.TempFile(os.TempDir(), "downloader")
		if lerr != nil {
			continue
		}

		_, lerr = io.Copy(file, resp.Body)
		resp.Body.Close()

		if lerr != nil {
			continue
		}
	}

	if lerr != nil {
		return nil, lerr
	}

	return file, nil
}
