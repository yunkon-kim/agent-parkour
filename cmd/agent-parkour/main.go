package main

import (
	"os"
	"os/exec"
	"syscall"
)

// Main entrypoint for agent-parkour wrapper (invokes parkour or runs embedded command)
func main() {
	// If parkour binary is available in PATH, exec it directly
	if binPath, err := exec.LookPath("parkour"); err == nil {
		_ = syscall.Exec(binPath, os.Args, os.Environ())
	}
	// Fallback to calling parkour in the same directory or printing notice
	cmd := exec.Command("parkour", os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}
