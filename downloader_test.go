package downloader

import (
	"os"
	"path"
	"testing"
	"time"
)

const (
	rightUrl = "http://download.geonames.org/export/dump/countryInfo.txt"
	wrongUrl = "http://download.geonames.org/export/dump/countryInfo.txt_wrong"
)

func TestSaveToTempFile(t *testing.T) {
	file, err := New(rightUrl).SaveToTempFile()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("downloaded file: %s", file)

	file, err = New(wrongUrl).Timeout(time.Second).Retries(5).SaveToTempFile()
	if err != ErrNotFound {
		t.Fail()
	}

	if file != "" {
		t.Fail()
	}
}

func TestSaveToFile(t *testing.T) {
	filename := path.Join(os.TempDir(), "downloader_test")

	err := New(rightUrl).SaveToFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("downloaded file: %s", filename)

	err = New(wrongUrl).Timeout(time.Second).Retries(5).SaveToFile(filename)
	if err != ErrNotFound {
		t.Fail()
	}
}
