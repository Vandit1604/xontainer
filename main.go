package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

type SysProcIDMap struct {
}

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

		We're creating a new usernamespace which enables us to run the program not as a root user too
		ALTHOUGH, we're not mapping the user id's to the new namespace, so the user will not be root in the new namespace
	*/
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER,
		// these mappings will give the new user in user namespace a root identity
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        0,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        0,
			},
		},
	}

	// running the command
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error while running the command %v:", err)
		os.Exit(cmd.ProcessState.ExitCode())
	}
}
