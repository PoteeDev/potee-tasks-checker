package main

import (
	"bufio"
	"context"
	"os/exec"
	"strings"
	"time"
)

func Execute(command string) (string, error) {
	args := strings.Split(command, " ")

	ctx := context.Background()
	if *timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	r, _ := cmd.StdoutPipe()
	//cmd.Stderr = cmd.Stdout

	done := make(chan struct{})

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(r)

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	output := ""

	go func() {
		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			output += line + "\n"
			//log.Println(line)
		}
		// We're all done, unblock the channel
		done <- struct{}{}
	}()

	// Start the command and check for errors
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	// Wait for all output to be processed
	<-done

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	return output, nil
}
