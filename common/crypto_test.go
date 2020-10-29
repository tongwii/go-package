package common

import "testing"

func TestHashFile(t *testing.T) {
	file := "C:\\Windows"
	mode := "crc32"
	if hash, err := HashFile(file, mode); err != nil {
		t.Error("dir error", err)
	} else {
		t.Log(hash)
	}
}
