package main

import (
	"log"
	"os"
)

func limitCPUandMemory(newroot string, memorylimit string, cpulimit string) error {
	// making directory for memory limiting
	if err := os.MkdirAll((newroot + `/sys/fs/cgroup/memory/container`), 0700); err != nil {
		log.Fatalf(`Error creating directory for memory limit: %v`, err)
		return err
	}

	// making directory for cpu limiting
	if err := os.MkdirAll((newroot + `/sys/fs/cgroup/cpu/container`), 0700); err != nil {
		log.Fatalf(`Error creating directory for cpu limit: %v`, err)
		return err
	}

	// setting the memory limit
	if err := os.WriteFile((newroot + `/sys/fs/cgroup/memory/container/memory.limit_in_bytes`), []byte(memorylimit), 0666); err != nil {
		log.Fatalf(`Error setting the memory limit : %v`, err)
		return err
	}

	// setting the cpu limit
	if err := os.WriteFile((newroot + `/sys/fs/cgroup/cpu/container/cpu.shares`), []byte(cpulimit), 0666); err != nil {
		log.Fatalf(`Error setting the cpu limit : %v`, err)
		return err
	}

	if err := os.WriteFile((newroot + `/sys/fs/cgroup/memory/container/cgroup.procs`), []byte(`1`), 0666); err != nil {
		log.Fatalf(`Error adding container process to the cgroup processes /proc : %v`, err)
		return err
	}

	if err := os.WriteFile((newroot + `/sys/fs/cgroup/cpu/container/cgroup.procs`), []byte(`1`), 0666); err != nil {
		log.Fatalf(`Error adding container process to the cgroup processes /proc : %v`, err)
		return err
	}

	return nil
}
