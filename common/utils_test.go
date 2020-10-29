package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetCurrentDirectory(t *testing.T) {
	dir := GetCurrentDirectory()
	t.Log(filepath.Join(dir, "demo"))
}

func TestIsPathExist(t *testing.T) {
	dir := "C:\\Windows"
	if _, err := IsPathExist(dir); err != nil {
		t.Error("dir error", err)
	}
}

func TestStringIPToByte(t *testing.T) {
	ip := StringIPToByte("192.168.0.1")
	t.Logf("%+v", ip)
}

func TestGetExeName(t *testing.T) {
	t.Logf("%s", GetExeName())
}

func TestCopyFile(t *testing.T) {
	if err := CopyFile("/ccms/ccms", "/mnt/ccms2", 0); err != nil {
		t.Error("dir error", err)
	}
}

func TestRestart(t *testing.T) {
	t.Logf("%s", os.Args)
}
