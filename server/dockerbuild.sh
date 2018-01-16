docker run --rm -v "$PWD":/mope/server -w /mope/server/src/roomserver -e GOPATH=/mope/server golang go build -o ../../bin/roomserver 
