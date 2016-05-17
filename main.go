package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	//   for _, arg := range os.Args {
	//     fmt.Println(arg)
	//   }
	timeout, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("invalid timeout argument: " + err.Error())
		os.Exit(1)
	}
	cmd := os.Args[2]
	args := os.Args[3:len(os.Args)]
	stdout, err := sysexec(timeout, cmd, args...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Print(stdout)
		os.Exit(0)
	}
}

func sysexec(timeout int, cmd string, args ...string) (stdout string, err error) {
	arg := append([]string{cmd}, args...)
	arg_str := fmt.Sprintf("%s", strings.Join(arg, " "))
	process := exec.Command("/bin/bash", "-l", "-c", arg_str)

	cmdOutput := &bytes.Buffer{}
	process.Stdout = cmdOutput
	process.Stderr = cmdOutput

	done := make(chan error)
	process.Start()
	go func() {
		done <- process.Wait()
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		// 超时处理
		if err := process.Process.Signal(syscall.SIGKILL); err != nil {
			return fmt.Sprintf("kill -9 %d", process.Process.Pid), errors.New(fmt.Sprintf("(PID:%d)ERROR:kill falure", process.Process.Pid))
		}
		return "killed", errors.New("timeout")

	case err = <-done:
		// 完成处理
		ori_output := cmdOutput.Bytes()
		return string(ori_output), err
	}

	return
}
