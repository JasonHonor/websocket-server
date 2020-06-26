go build -ldflags="-s -w -X 'main.BuildID=$1'" -o client2
#gzip client2
go build -ldflags="-s -w -X 'main.BuildID=$1'" -o client
./client run