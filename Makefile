.PHONY: build-linux build-arm

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/Hivemind-Studio/isi-core/cmd/api/

build-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build  github.com/Hivemind-Studio/isi-core/cmd/api/