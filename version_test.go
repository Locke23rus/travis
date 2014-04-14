package main

import (
	"fmt"
	"os/exec"
	"testing"
)

var versionArgs = [][]string{{"version"}, {"--version"}, {"-v"}}

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
