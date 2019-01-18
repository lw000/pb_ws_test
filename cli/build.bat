cd ../../../../
set GOPATH=%cd%
cd src/demo/pb_ws_test/cli/
set GOARCH=amd64
set GOOS=windows
go build -v -ldflags="-s -w"