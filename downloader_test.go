package downloader

import (
	"testing"
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

	file, err = New(wrongUrl).Retries(5).SaveToTempFile()
	if err == nil {
		t.Fail()
	}

	if file != "" {
		t.Fail()
	}
}
