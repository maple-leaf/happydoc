package helpers

import (
	"os/exec"
	"strings"
)

// run a command using the shell; no need to split args
// from https://stackoverflow.com/questions/6182369/exec-a-shell-command-in-go
func RunShellCmd(cmd string, shell bool) (out []byte, err error) {
	if shell {
		out, err = exec.Command("bash", "-c", cmd).Output()
		return
	}
	// https://sj14.gitlab.io/post/2018-07-01-go-unix-shell/
	cmds := strings.Split(cmd, " ")
	out, err = exec.Command(cmds[0], cmds[1:]...).Output()

	return
}
