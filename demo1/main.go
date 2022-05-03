package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd *exec.Cmd
		err error
	)
	cmd = exec.Command("/bin/bash", "-c", "echo 1")
	cmd.Run()
	fmt.Println(err)
}
