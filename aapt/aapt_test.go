package aapt

import "testing"
import "strings"

func TestFindLastAaptPath(t *testing.T) {
	appPath, err := FindLastAaptPath()
	if err != nil {
		t.Error(err)
		return
	}
	fileNames := strings.Split(appPath, "/")
	if fileNames[len(fileNames)-1] != "aapt" {
		t.Error("It not a aapt cmd", appPath)
	}
	if fileNames[len(fileNames)-2] != "25.0.2" {
		t.Error("Got error version path,", appPath)
	}
}
