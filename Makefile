build:
	@go build -o xontainer main.go rootfs.go

run: build
	@./xontainer


.PHONY: build run
