.PHONY: build-linux

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/Hivemind-Studio/isi-core/cmd/api/