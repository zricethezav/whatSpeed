.PHONY: test build-all deploy

test:
	go get github.com/golang/lint/golint
	golint
build-all:
	rm -rf build
	mkdir build
	env GOOS="windows" GOARCH="amd64" go build -o "build/whatspeed-windows-amd64.exe"
	env GOOS="windows" GOARCH="386" go build -o "build/whatspeed-windows-386.exe"
	env GOOS="linux" GOARCH="amd64" go build -o "build/whatspeed-linux-amd64"
	env GOOS="linux" GOARCH="arm" go build -o "build/whatspeed-linux-arm"
	env GOOS="linux" GOARCH="mips" go build -o "build/whatspeed-linux-mips"
	env GOOS="linux" GOARCH="mips" go build -o "build/whatspeed-linux-mips"
	env GOOS="darwin" GOARCH="amd64" go build -o "build/whatspeed-darwin-amd64"
