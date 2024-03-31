# xontainer

NOTE: Download the mini-root filesystem from https://www.alpinelinux.org/downloads/. We will test the alpine linux image inside the container

Do read about mount namespaces as they were really interesting: https://man7.org/linux/man-pages/man7/mount_namespaces.7.html
TL;DR:

---
“When a process creates a new mount namespace using clone(2) or unshare(2) with the CLONE_NEWNS flag, the mount point list for the new namespace is a copy of the caller’s mount point list.”
---

## Steps before using

1. Run this command as we will be using this as our alpine filesystem that we will be mounted.

From now on, xontainer will expect a root filesystem to exist at `/tmp/xontainer/rootfs` and will raise an error if one can’t be found. Note that although we’re using BusyBox for this particular example, you could just as easily use any other distro.

```bash
$ mkdir -p /tmp/xontainer/rootfs
$ tar -C /tmp/xontainer/rootfs -xf assets/alpine-minirootfs-3.19.1-x86_64.tar.gz
```

