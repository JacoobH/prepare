package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)

	//generated cmd
	cmd = exec.Command("/bin/bash", "-c", "sleep5; ls -l")

	//The command was executed to capture the output of the child process(pipe)
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}

	//Print child process output
	fmt.Println(string(output))
}
