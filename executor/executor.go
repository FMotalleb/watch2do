package executor

import (
	"os/exec"
	"strings"

	"github.com/fmotalleb/watch2do/cmd"
)

func RunCommands() {
	setupLog()

	proc := strings.Split(cmd.Params.Shell, " ")

	for _, command := range cmd.Params.Commands {

		log.Infof("executing `%s` with args: %v %v", proc[0], proc[1:], command)
		process := exec.Command(proc[0], append(proc[1:], command)...)
		initErr := process.Start()
		if initErr != nil {
			log.Warnf("Launching process %v failed: %v", command, initErr)
		}
		waitErr := process.Wait()
		if waitErr != nil {
			log.Warnf("Launching process %v failed: %v", command, waitErr)
		}
	}
}
