package common

import (
	"testing"
)

func TestYaml(t *testing.T) {
	file := "server.yml"
	dir := "C:\\TongSync\\development\\code\\go\\src\\ccms\\bin\\server\\config\\"
	var data interface{}
	err := ReadConfig(file, dir, &data)
	if err != nil {
		t.Error(err.Error())
		return
	} else {
		t.Logf("%+v", data)
	}

	err = WriteConfig(file, dir, data)
	if err != nil {
		t.Error(err.Error())
		return
	} else {
		t.Logf("%+v", data)
	}
}
