go build -ldflags="-s -w -X 'main.BuildID=$1'" -o gfx
scp ./gfx root@192.168.61.11:/opt/syscenter/syscenter
scp ./template/deploy.html root@192.168.61.11:/opt/syscenter/template/deploy.html