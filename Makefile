build:
	@go build -o xontainer main.go

run: build
	@./xontainer


.PHONY: build run
