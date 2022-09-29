main: linux

linux: init
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux  go build -o bin/fax_calc_linux_amd64 main.go

windows: init
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o bin/fax_calc_windows_amd64.exe main.go

darwin: init
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o bin/fax_calc_darwin_amd64 main.go

init:
	echo "开始打包..."