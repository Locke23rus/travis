package main

import (
	"os/exec"
	"strings"
	"testing"
)

const helpHint = "Usage: travis COMMAND"

var helpArgs = [][]string{{}, {"help"}, {"--help"}, {"-h"}, {"-?"}}

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
