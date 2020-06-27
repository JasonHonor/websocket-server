go build -ldflags="-s -w -X 'main.BuildID=$1'" -o client2
scp ./client2 root@192.168.61.11:/opt/syscenter/update/
