package aapt

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	// XML_Manifest is the useful name of "AndroidManifest.xml"
	XML_Manifest = "AndroidManifest.xml"
	// AAPT executeable bin name
	AAPT_EXEC_NAME = "aapt2"
)

/*
The Aapt control
*/
type Aapt struct {
	appt string
}

// NewAapt is to create a AAPT by aapt executer path
func NewAapt(aaptPath string) *Aapt {
	return &Aapt{
		appt: aaptPath,
	}
}

// to sort files
type dirSort []os.FileInfo

func (d *dirSort) Len() int {
	return len(*d)
}

func (d *dirSort) Less(i, j int) bool {
	return (*d)[i].Name() > (*d)[j].Name()
}

func (d *dirSort) Swap(i, j int) {
	(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
}

// FindLastAaptPath is to find the aapt with last version.
// Must declare ANDROID_HOME in ENV.
// Error will be not nil when aapt found, and not empty path return too.
func FindLastAaptPath() (string, error) {
	androiHome := os.Getenv("ANDROID_HOME")
	if androiHome != "" {
		toolsDir := filepath.Join(androiHome, "build-tools")
		files, err := ioutil.ReadDir(toolsDir)
		if err != nil {
			return "", err
		}

		if len(files) > 0 {
			// sort
			// sorted := dirSort(files)
			// sort.Sort(&sorted)
			// find last version
			findLast := func(pool []os.FileInfo) os.FileInfo {
				last := pool[0]
				for _, item := range pool {
					if last.Name() < item.Name() {
						last = item
					}
				}
				return last
			}
			lastTools := findLast(files)

			aaptExec := AAPT_EXEC_NAME
			if runtime.GOOS == "windows" {
				aaptExec = aaptExec + ".exe"
			}
			return filepath.Join(toolsDir, lastTools.Name(), aaptExec), nil
		}

		return "", errors.New("No build-tools found in " + androiHome)
	}
	return "", errors.New("ANDROID_HOME not set")
}

// LineFilterFunc is to check and filter a text line
// if passed return "true", or return "false" to break
type LineFilterFunc func(line []byte) bool

// DumpXmlTrees to dump the content of xml file which name is "xmlname" in apk.
// If filter not nil, can use it to process ever line the aapt return, if filter return false
// the line process will be broken, and DumpXmlTrees will be return
func (a *Aapt) DumpXmlTrees(apk string, xmlname string, filter LineFilterFunc) (result string, err error) {
	cmd := exec.Command(a.appt, "dump", "xmltree", "--file", xmlname, apk)
	return runCmd(cmd, filter)
}

// DumpBadging will return the badgings in apk. like cmd "appt d badging apk".
// Can use it to get the packageName of apk.
// If filter not nil, can use it to process ever line the aapt return, if filter return false
// the line process will be broken, and DumpBadging will be return
func (a *Aapt) DumpBadging(apk string, filter LineFilterFunc) (result string, err error) {
	cmd := exec.Command(a.appt, "dump", "badging", "--include-meta-data", apk)
	return runCmd(cmd, filter)
}

// runCmd to run a exec.Cmd, and process cmd results by filter when it not nil
func runCmd(cmd *exec.Cmd, filter LineFilterFunc) (result string, err error) {

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	bufout := bufio.NewReader(stdout)
	readLines := bytes.NewBuffer(make([]byte, 0, 8192))

	func() {
		defer f(bufout)
		for {
			lineBuf, _, rErr := bufout.ReadLine()
			if rErr != nil {
				break
			}
			readLines.Write(lineBuf)
			readLines.WriteString("\n")
			if filter != nil && !filter(lineBuf) {
				break
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return
	}

	return readLines.String(), err
}

func f(rd *bufio.Reader) {
	for {
		_, _, err := rd.ReadLine()
		if err != nil {
			break
		}
	}
}
