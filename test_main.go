package main

import (
	"fmt"
	"goadb"
)

const (
	Test_Adb string = "/home/holmes/test/AppCool_203.apk"
	Test_Pkg string = "com.mgyapp.android"
)

func main() {
	goadb.Debug = false
	adbPath := "/home/holmes/prosoft/android-sdk-linux/platform-tools/adb"
	adb := goadb.NewGoAdb(adbPath)
	fmt.Println("adb path is:")
	fmt.Println(adb.GetAdbPath())
	fmt.Println("adb deivces:")
	devices := adb.Devices()
	fmt.Println("connected", len(devices), "device")
	for _, dev := range devices {
		if dev != nil {
			fmt.Println("  -- deviceId", dev.GetDeviceId(), "status", dev.GetStatus())
		}
	}
	fmt.Println("adb versoin", adb.Version())

	fmt.Println("will install", Test_Adb)
	inststallResult, isError := adb.Install(Test_Adb, true, false, false, false)
	fmt.Println("error", isError, "install", inststallResult)

	//fmt.Println("Will uninstall", Test_Pkg)
	//uninstall, isError := adb.Uninstall(Test_Pkg, false)
	//fmt.Println("error", isError, "Uninstall", uninstall)
}
