package telphin

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func getTestData(t *testing.T, file string) []byte {
	path := filepath.Join("./../testdata", file) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
