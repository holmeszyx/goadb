package aapt

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestFindLastAaptPath(t *testing.T) {
	appPath, err := FindLastAaptPath()
	if err != nil {
		t.Error(err)
		return
	}
	fileNames := strings.Split(appPath, string(filepath.Separator))
	if !strings.HasPrefix(fileNames[len(fileNames)-1], AAPT_EXEC_NAME) {
		t.Error("It not a aapt cmd", appPath)
	}
	if fileNames[len(fileNames)-2] != "30.0.2" {
		t.Error("Got error version path,", appPath)
	}
}
