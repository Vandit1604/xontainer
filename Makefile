build:
	@go build -o xontainer main.go rootfs.go cgroup.go

run: build
	@./xontainer


.PHONY: build run
