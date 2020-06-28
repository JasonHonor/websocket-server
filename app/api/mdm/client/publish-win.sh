CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X 'main.BuildID=$1'" -o win-client.exe
scp ./win-client.exe root@192.168.61.11:/opt/syscenter/update/
