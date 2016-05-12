GoAdb
=====

adb with go interface.

Let's us use adb by golang simplely.

See docs on [Godoc](https://godoc.org/github.com/holmeszyx/goadb)


Usage
======

```go
adb, err := goadb.NewAdbWithEnv()
if err != nil {
    return
}

// connect the phone

var output string

// install apk by adb
output, err = adb.Install("/home/xxx/xx.apk", true, false, false, false)
if err != nil {
    fmt.Println(err)
    return
}

fmt.Println(output)

// execute shell command in phone
// only the pure shell command, not need "adb shell"
output, err = adb.ShellCmd("getprop")
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(output)
// output, _ := adb.ShellCmd("pm list package -f -3")

// reboot phone
_, err := adb.Reboot()
if err != nil {
    fmt.Println(err)
    return
}

// reboot to recovery
_, err := adb.RebootTo(goadb.MODE_RECOVERY)
if err != nil {
    fmt.Println(err)
    return
}


```
