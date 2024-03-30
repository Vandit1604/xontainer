package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
)

// init is called before main function. This automatically registers the commands inside the nsInitialisation
func init() {
	// via `Register` we can register the functions that we will use inside the namespace that we are creating for a container
	reexec.Register("nsInitialisation", nsInitialisation)
	// via Init we check if the registered function was actually exec'd or not
	if reexec.Init() {
		os.Exit(0)
	}
}

func nsInitialisation() {
	fmt.Printf("\n>> ANYTHING THAT WE WANT TO DO INSIDE THE NAMESPACE <<\n")
	// this runs the command in the new namespace
	nsRun()
}

func nsRun() {
	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// setting the env for better distinction between the namespaces
	cmd.Env = []string{"PS1=-[xontainer]- # "}

	// running the command
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error while running the command %v:", err)
		os.Exit(cmd.ProcessState.ExitCode())
	}
}

func main() {
	cmd := reexec.Command("nsInitialisation")

	// pipe the stdin/out/err of os to cmd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	/*
					FOR GENERAL INFORMATION: https://man7.org/linux/man-pages/man7/namespaces.7.html
					syscall.CLONE_NEWUTS - new uts namespace which gives us a namespaced hostname and domain name to the process
					syscall.CLONE_NEWPID - new process id for this process in the namespace it will be 1
					syscall.CLONE_NEWIPC - https://www.man7.org/linux/man-pages/man7/ipc_namespaces.7.html
					syscall.CLONE_NEWNET - https://man7.org/linux/man-pages/man7/network_namespaces.7.html
					syscall.CLONE_NEWUSER - https://man7.org/linux/man-pages/man7/user_namespaces.7.html

		CLONE_NEWNS: https://man7.org/linux/man-pages/man7/mount_namespaces.7.html
		This flag has the same effect as the clone(2) CLONE_NEWNS
		flag.  Unshare the mount namespace, so that the calling
		process has a private copy of its namespace which is not
		shared with any other process.  Specifying this flag
		automatically implies CLONE_FS as well.  Use of
		CLONE_NEWNS requires the CAP_SYS_ADMIN capability.  For
		further information, see mount_namespaces(7).

		We're creating a new usernamespace which enables us to run the program not as a root user too
		ALTHOUGH, we're not mapping the user id's to the new namespace, so the user will not be root in the new namespace
	*/
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER, // these mappings will give the new user in user namespace a root identity
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	// starting the command
	if err := cmd.Start(); err != nil {
		log.Fatalf("Error while starting the command %v:", err)
		os.Exit(cmd.ProcessState.ExitCode())
	}

	// waiting for the command
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Error while waiting for the command %v:", err)
		os.Exit(cmd.ProcessState.ExitCode())
	}
}
