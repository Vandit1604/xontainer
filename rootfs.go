package main

import (
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func pivotRoot(newRoot string) error {
	// .pivot_root is a hidden directory in the new root
	putOld := filepath.Join(newRoot, "/.pivot_root")

	// mount command attaches the newRoot inside the newRoot which means we attach the newRoot to newRoot; e.g if we mount /somedir into /dir. The content of /somedir gets inside the /dir
	if err := syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		log.Fatalf("Couldn't mount the newRoot: %v", err)
	}

	if err := os.MkdirAll(putOld, 0700); err != nil {
		log.Fatalf("Error creating a hidden pivot_root directory in root filesystem of the host: %v", err)
	}

	if err := syscall.PivotRoot(newRoot, putOld); err != nil {
		log.Fatalf("Could not pivot root: %v", err)
	}

	// ensure current working directory is set to new root
	if err := os.Chdir("/"); err != nil {
		log.Fatalf("Couldn't change directory to /: %v", err)
	}

	// umount putold, which now lives at /.pivot_root
	putOld = "/.pivot_root"
	if err := syscall.Unmount(
		putOld,
		syscall.MNT_DETACH,
	); err != nil {
		return err
	}

	// remove putold
	if err := os.RemoveAll(putOld); err != nil {
		return err
	}

	return nil
}

func mountProc(newroot string) error {
	source := "proc"
	target := filepath.Join(newroot, "/proc")
	fstype := "proc"
	flags := 0
	data := ""

	os.MkdirAll(target, 0755)
	if err := syscall.Mount(
		source,
		target,
		fstype,
		uintptr(flags),
		data,
	); err != nil {
		return err
	}

	return nil
}
