package ansibleutils

import (
	"os/exec"
	"syscall"
)

// command runs a command and sends back the output, exitcode and error
func command(params []string) ([]byte, int, error) {
	var cmd *exec.Cmd
	if len(params) == 1 {
		cmd = exec.Command(params[0])
	} else if len(params) > 1 {
		cmd = exec.Command(params[0], params[1:]...)
	}

	//params = append([]string{"sh", "-c"}, params...)

	out, err := cmd.CombinedOutput()

	exitCode := 0
	if msg, ok := err.(*exec.ExitError); ok { // there is error code
		exitCode = msg.Sys().(syscall.WaitStatus).ExitStatus()
	}

	return out, exitCode, err
}
