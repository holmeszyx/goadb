package goadb

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

const (
	Unknown string = "Unknown"
)

var (
	// It's debug
	Debug bool = false
)

// A connected Android Device
type Device struct {
	deviceId string
	status   string
}

// new a Deivce
func newDevice(line []string) *Device {
	device := new(Device)
	device.deviceId = line[0]
	device.status = line[1]
	return device
}

// Get device Id
func (d *Device) GetDeviceId() string {
	return d.deviceId
}

// get device status
// device
// offline
func (d *Device) GetStatus() string {
	return d.status
}

// Create a Adb
func NewGoAdb(adbPath string) *GoAdb {
	adb := new(GoAdb)
	adb.adbPath = adbPath
	adb.init()
	return adb
}

type GoAdb struct {
	adbPath         string
	connectedDevice []*Device
}

// Init
func (g *GoAdb) init() {

}

// Get conntect Devices
// maybe []
func (g *GoAdb) GetConnectedDevice() []*Device {
	return g.connectedDevice
}

// Set Adb path
func (g *GoAdb) SetAdbPath(path string) {
	g.adbPath = path
}

// Get Adb path
func (g *GoAdb) GetAdbPath() string {
	return g.adbPath
}

// Get Adb verison
func (g *GoAdb) Version() string {
	result, _ := g.runAdb("version")
	index := strings.LastIndex(result, "version")
	if index > -1 {
		result = result[index+8:]
	}
	if len(result) == 0 {
		return Unknown
	}
	return result
}

// Get devices
func (g *GoAdb) Devices() []*Device {
	result, _ := g.runAdb("devices")
	lines := strings.Split(result, "\n")
	var devices = []*Device{}
	for i, line := range lines {
		if Debug {
			fmt.Println(i, line, len(line))
		}
		if len(line) > 0 {
			deviceLine := strings.Split(line, "	")
			//fmt.Println("device line", deviceLine, len(deviceLine))
			if len(deviceLine) < 2 {
				continue
			}
			dev := newDevice(deviceLine)
			//fmt.Println(dev)
			if dev != nil {
				devices = append(devices, dev)
			}
		}
	}
	g.connectedDevice = devices
	return devices
}

// Install Apk
//  adb install [-l] [-r] [-d] [-s]
func (g *GoAdb) Install(apk string, reinstall bool, forward bool, downgrade bool, toSd bool) (string, bool) {
	var args string
	if reinstall {
		args += " -r"
	}
	if forward {
		args += " -l"
	}
	if downgrade {
		args += " -d"
	}
	if toSd {
		args += " -s"
	}

	args += " " + apk

	result, isError := g.runAdb("install" + args)

	return result, isError
}

// Uninstall package
//   adb uninstall [-k] <package> - remove this app package from the device
func (g *GoAdb) Uninstall(pkg string, keepData bool) (string, bool) {
	var args string
	if keepData {
		args += " -k"
	}

	args += " " + pkg
	result, isError := g.runAdb("uninstall" + args)
	return result, isError
}

// run a adb cmd
// get output
// When exec error, isError = true
func (g *GoAdb) runAdb(cmd string) (string, bool) {
	cmdArgs := strings.Split(cmd, " ")
	adbExec := exec.Command(g.adbPath, cmdArgs...)

	isError := false

	in, _ := adbExec.StdinPipe()
	error, _ := adbExec.StderrPipe()
	out, _ := adbExec.StdoutPipe()
	defer closeIO(in)
	defer closeIO(error)
	defer closeIO(out)

	if err := adbExec.Start(); err != nil {
		panic("start adb process error")
	}

	outData, _ := ioutil.ReadAll(out)
	errorData, _ := ioutil.ReadAll(error)

	if err := adbExec.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			fmt.Println("adb error return")
			outData = errorData
			isError = true
		} else {
			panic("wait adb process error")
		}
	}

	return string(outData), isError
}

// close a stream
func closeIO(c io.Closer) {
	if c != nil {
		c.Close()
	}
}
