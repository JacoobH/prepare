package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	output []byte
	err    error
}

func main() {
	//Execute a CMD, let it run in a coroutine, let it run for 2 seconds(sleep 2sec)
	//At 1 second, kill cmd
	var (
		ctx        context.Context
		cancelFunc context.CancelFunc
		cmd        *exec.Cmd
		resultChan chan *result
		res        *result
	)
	resultChan = make(chan *result, 1000)
	ctx, cancelFunc = context.WithCancel(context.TODO())
	go func() {
		var (
			output []byte
			err    error
		)
		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 2; echo hello;")
		output, err = cmd.CombinedOutput()
		resultChan <- &result{
			output: output,
			err:    err,
		}
	}()

	//Keep going
	time.Sleep(1 * time.Second)

	//Cancel context
	cancelFunc()

	res = <-resultChan
	fmt.Println(res.err, string(res.output))
}
