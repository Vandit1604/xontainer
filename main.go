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

	/*
		FOR GENERAL INFORMATION: https://man7.org/linux/man-pages/man7/namespaces.7.html

				syscall.CLONE_NEWUTS - new uts namespace which gives us a namespaced hostname and domain name to the process
				syscall.CLONE_NEWPID - new process id for this process in the namespace it will be 1
				syscall.CLONE_NEWIPC - https://www.man7.org/linux/man-pages/man7/ipc_namespaces.7.html
				syscall.CLONE_NEWNET - https://man7.org/linux/man-pages/man7/network_namespaces.7.html
				syscall.CLONE_NEWUSER - https://man7.org/linux/man-pages/man7/user_namespaces.7.html
	*/
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER,
	}

	// running the command
	if err := cmd.Run(); err != nil {
		log.Fatalf("Erro while runningt the command %v:", err)
		os.Exit(cmd.ProcessState.ExitCode())
	}
}
