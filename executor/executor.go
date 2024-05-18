package executor

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fmotalleb/watch2do/cmd"
	"github.com/fmotalleb/watch2do/fallback"
	"github.com/fmotalleb/watch2do/logger"
	"github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func setupLog() {
	log = logger.SetupLogger("Executor")
}

func panicOn(log *logrus.Entry, note string, err error) {
	if err != nil {
		log.WithField("error", err).Panicln(note)
	}
}
func warnOn(log *logrus.Entry, note string, err error) {
	if err != nil {
		log.WithField("error", err).Warnln(note)
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

	logger := log.WithField("cycle_number", counter)
	killOldInstances(logger)

	proc := strings.Split(cmd.Params.Shell, " ")
	go func() {
		for _, command := range cmd.Params.Commands {
			logger := logger.WithField("process", command)
			args := append(proc[1:], command)
			logger.Infof("executing `%s` with args: %v", proc[0], args)
			process := exec.Command(proc[0], args...)
			if cmd.Params.LogLevel == logrus.DebugLevel {
				stdout, err := process.StdoutPipe()
				log.WithFields(logrus.Fields{
					"error": err,
					"pipe":  "stdout",
				}).Warningln("cannot get processes stdout pipe")
				warnOn(logger, "cannot get stdout of child", err)
				stderr, err := process.StderrPipe()
				warnOn(logger, "cannot get stderr of child", err)
				go io.Copy(os.Stderr, stdout)
				go io.Copy(os.Stderr, stderr)
			}
			logger.Debugln("stdout and stderr attached")

			initErr := process.Start()
			panicOn(logger, "failed to start child process", initErr)
			pid := process.Process.Pid
			push(pid)
			logger = logger.WithField("pid", pid)
			logger.Debugln("process started")

			waitErr := process.Wait()
			warnOn(logger, "cannot wait for process", waitErr)
			logger.Debugln("process done")
			pop()
		}
		// close(pids)
	}()

}

func killOldInstances(logger *logrus.Entry) {
	for _, pid := range pids {
		killProcess(logger, pid)
	}
}

func killProcess(logger *logrus.Entry, pid int) {
	processes, err := process.Processes()
	warnOn(logger, "Failed to retrieve processes", err)
	if err != nil {
		return
	}

	for _, p := range processes {
		ppid, err := p.Ppid()
		if err != nil {
			continue
		}
		if int(ppid) == pid {
			killProcess(logger, int(p.Pid))
			p.Kill()
		}
		if int(p.Pid) == pid {
			p.Kill()
		}
	}
}

// func killProcess(logger *logrus.Entry, pid int) {
// 	pidStr := strconv.Itoa(pid)

// 	var cmd *exec.Cmd

// 	if runtime.GOOS == "windows" {
// 		cmd = exec.Command("taskkill", "/PID", pidStr, "/F", "/T")
// 	} else {
// 		cmd = exec.Command("pkill", "-P", pidStr)
// 	}

// 	err := cmd.Run()
// 	if err != nil {
// 		logger.Errorf("Failed to kill process with PID %d: %v", pid, err)
// 		return
// 	}

// 	logger.Debugf("Successfully killed process with PID %d", pid)
// }
