package downloader

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"time"
)

type Downloader struct {
	url     string
	retries int
}

func New(url string) *Downloader {
	return &Downloader{url, 3}
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
