$(shell mkdir -p bin)

.PHONY: linux
linux:
	@GOOS=linux GOARCH=amd64 go build -o bin/ipstore main.go