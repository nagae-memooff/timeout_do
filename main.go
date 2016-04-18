package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	} else {
		fmt.Print(stdout)
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
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		// 超时处理
		if err := process.Process.Kill(); err != nil {
			process.Process.Signal(os.Kill)
		}
		return "killed", errors.New("timeout")

	case err = <-done:
		// 完成处理
		ori_output := cmdOutput.Bytes()
		return string(ori_output), err
	}

	return
}