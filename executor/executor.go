package executor

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fmotalleb/watch2do/cmd"
	"github.com/fmotalleb/watch2do/fallback"
	"github.com/fmotalleb/watch2do/logger"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func setupLog() {
	log = logger.SetupLogger("Executor")
}

func panicOn(note string, err error) {
	if err != nil {
		log.Panicf("%s, Error: %s", note, err)
	}
}
func warnOn(note string, err error) {
	if err != nil {
		log.Warnf("%s, Error: %s", note, err)
	}
}

var pids []int = []int{}
var counter int = 0

func pop() (a int) {
	if len(pids) > 0 {
		a = pids[0]
		pids = pids[1:]
	}
	return
}
func push(pid int) {
	pids = append(pids, pid)
}

func RunCommands() {
	setupLog()
	counter++
	defer func() {
		fallback.CaptureError(log, recover())
	}()

	killOldInstances()

	proc := strings.Split(cmd.Params.Shell, " ")
	go func() {
		for _, command := range cmd.Params.Commands {
			logger := log.WithField("signal id", counter).WithField("process", command)
			args := append(proc[1:], command)
			logger.Infof("executing `%s` with args: %v", proc[0], args)
			process := exec.Command(proc[0], args...)
			if cmd.Params.LogLevel == logrus.DebugLevel {
				stdout, err := process.StdoutPipe()
				warnOn("cannot get stdout of child", err)
				stderr, err := process.StderrPipe()
				warnOn("cannot get stderr of child", err)
				go io.Copy(os.Stderr, stdout)
				go io.Copy(os.Stderr, stderr)
			}
			logger.Debugln("stdout and strerr attached")

			initErr := process.Start()
			panicOn("failed to start child process", initErr)
			push(process.Process.Pid)
			logger.Debugln("process started")

			waitErr := process.Wait()
			warnOn("cannot wait for process", waitErr)
			logger.Debugln("process done")
			pop()
		}
		// close(pids)
	}()

}

func killOldInstances() {
	for _, pid := range pids {
		log.Debugf("trying to kill process with pid %d\n", pid)
		process, err := os.FindProcess(pid)
		process.Release()
		if err != nil {
			warnOn("cannot find a process from pervious execution", err)
			continue
		}
		warnOn("cannot kill a process from pervious execution", process.Kill())
		// log.Debugf("killed process with pid %d",pid)
	}
	pids = []int{}
}
