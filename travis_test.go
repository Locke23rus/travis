package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const (
	helpHint       = "Usage: travis COMMAND"
	exitStatusZero = "exit status 0"
)

var (
	compiled    = false
	helpArgs    = [][]string{{}, {"help"}, {"--help"}, {"-h"}, {"-?"}}
	versionArgs = [][]string{{"version"}, {"--version"}, {"-v"}}
	travis      = ""
)

func compile() {
	if compiled {
		return
	}
	cwd, err := os.Getwd()
	exitIfErr(err)
	fmt.Println(cwd)
	err = exec.Command("go", "build").Run()
	exitIfErr(err)
	compiled = true
	travis = cwd + "/travis"
}

func TestHelpCommand(t *testing.T) {
	compile()
	for _, args := range helpArgs {
		cmd := exec.Command(travis, args...)
		out, err := cmd.Output()
		if err != nil {
			t.Errorf("%v %v", args, err)
		}
		if cmd.ProcessState.String() != exitStatusZero {
			t.Fail()
		}
		if !strings.Contains(string(out), helpHint) {
			t.Fail()
		}
	}

}

func TestVersionCommand(t *testing.T) {
	compile()
	for _, args := range versionArgs {
		cmd := exec.Command(travis, args...)
		out, err := cmd.Output()
		if err != nil {
			t.Errorf("%v %v", args, err)
		}
		if cmd.ProcessState.String() != exitStatusZero {
			t.Fail()
		}
		if string(out) != fmt.Sprintln(Version) {
			t.Fail()
		}
	}
}
