package goadb

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type AdbError struct {
	msg string
}

func (a *AdbError) Error() string {
	return a.msg
}

func newAdbError(msg string) *AdbError {
	return &AdbError{msg: msg}
}

const (
	Unknown         string = "Unknown"
	Empty           string = ""
	MODE_RECOVERY          = "recovery"
	MODE_BOOTLOADER        = "bootloader"
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

// Create a Adb, auto detact adb path which in Env
func NewGoAdbWithEnv() (*GoAdb, error) {
	androiHome := os.Getenv("ANDROID_HOME")
	if androiHome != "" {
		adbpath := androiHome + "/platform-tools/adb"
		return NewGoAdb(adbpath), nil
	} else {
		return nil, errors.New("ANDROID_HOME not exists in Env")
	}
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
func (g *GoAdb) Install(apk string, reinstall bool,
	forward bool, downgrade bool, toSd bool) (string, error) {
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

	args += " " + safeArg(strings.TrimSpace(apk))

	result, isError := g.runAdbCmd("install" + args)

	return result, isError
}

// Uninstall package
//   adb uninstall [-k] <package> - remove this app package from the device
func (g *GoAdb) Uninstall(pkg string, keepData bool) (string, error) {
	var args string
	if keepData {
		args += " -k"
	}

	args += " " + safeArg(strings.TrimSpace(pkg))
	result, isError := g.runAdbCmd("uninstall" + args)
	return result, isError
}

// Run command by Adb shell
// adb shell $cmd
func (g *GoAdb) ShellCmd(cmd string) (string, error) {
	cmd = safeArg(strings.TrimSpace(cmd))
	result, err := g.runAdbCmd("shell " + cmd)
	return result, err
}

// Push file to phone
// src is the file path on compute, and dst is the path in phone
func (g *GoAdb) Push(src string, dst string) error {
	args := []string{
		"push",
		src,
		dst,
	}
	_, err := g.runAdb(args...)
	return err
}

// Pull file from phone
// src is the path in phone, and dst is the path on compute
func (g *GoAdb) Pull(src string, dst string) error {
	_, err := g.runAdb("pull", src, dst)
	return err
}

// Kill-server
func (g *GoAdb) KillServer() (string, error) {
	return g.runAdb("kill-server")
}

// start-server
func (g *GoAdb) StartServer() (string, error) {
	return g.runAdb("start-server")
}

// Reboot phone
func (g *GoAdb) Reboot() (string, error) {
	return g.runAdb("reboot")
}

// RebootTo like command
// "reboot [bootloader|recovery]"
// "to" can be MODE_BOOTLOADER or MODE_RECOVERY
func (g *GoAdb) RebootTo(to string) (string, error) {
	return g.runAdb("reboot", to)
}

// run adb cmd string
// Use "\ " instead of " " like shell
func (g *GoAdb) runAdbCmd(cmd string) (string, error) {
	// cmdArgs := strings.Split(cmd, " ")
	cmdArgs := splitCmdAgrs(cmd)
	return g.runAdb(cmdArgs...)
}

// run a adb cmd which is a slice or array
// get output
// When exec error, error != nil
func (g *GoAdb) runAdb(cmd ...string) (string, error) {
	adbExec := exec.Command(g.adbPath, cmd...)

	in, _ := adbExec.StdinPipe()
	errorOut, _ := adbExec.StderrPipe()
	out, _ := adbExec.StdoutPipe()
	defer closeIO(in)
	defer closeIO(errorOut)
	defer closeIO(out)

	if err := adbExec.Start(); err != nil {
		return Empty, errors.New("start adb process error")
	}

	outData, _ := ioutil.ReadAll(out)
	errorData, _ := ioutil.ReadAll(errorOut)

	var adbError error = nil

	if err := adbExec.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			adbError = newAdbError("adb return error")
			outData = errorData
		} else {
			return Empty, errors.New("start adb process error")
		}
	}

	return string(outData), adbError
}

// close a stream
func closeIO(c io.Closer) {
	if c != nil {
		c.Close()
	}
}
