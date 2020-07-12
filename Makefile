build: build_windows build_linux build_mac

build_windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./net-wait-go-windows.exe . && upx ./net-wait-go-windows.exe

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./net-wait-go-linux . && upx ./net-wait-go-linux

build_mac:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./net-wait-go-mac . && upx ./net-wait-go-mac