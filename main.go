package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("/bin/sh")

	// pipe the stdin/out/err of os to cmd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// setting the env for better distinction between the namespaces
	cmd.Env = []string{"PS1=-[xontainer]- # "}

	// new uts namespace which gives us a namespaced hostname and domain name to the process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	// running the command
	if err := cmd.Run(); err != nil {
		log.Fatalf("Erro while runningt the command %v:", err)
		os.Exit(cmd.ProcessState.ExitCode())
	}
}
