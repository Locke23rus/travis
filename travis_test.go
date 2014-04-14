package main

import (
	"os"
	"os/exec"
)

const (
	exitStatusZero = "exit status 0"
)

var (
	compiled = false
	travis   = ""
)

func compile() {
	if compiled {
		return
	}
	cwd, err := os.Getwd()
	exitIfErr(err)
	err = exec.Command("go", "build").Run()
	exitIfErr(err)
	compiled = true
	travis = cwd + "/travis"
}
