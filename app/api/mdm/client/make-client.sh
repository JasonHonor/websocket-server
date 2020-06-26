go build -ldflags="-X 'main.BuildID=$1'" -o client2
go build -ldflags="-X 'main.BuildID=$1'" -o client
./client run