version=`date "+%Y%m%d%H%M%S"`
gitcommit=`git rev-parse --short HEAD`

all: client server

#
# Standalone binaries
#
client: deps *.go client/*
	go build \
		--ldflags "-s -w -X main.version=${version} -X main.gitCommit=${gitcommit} -X main.binaryName=certsync_client" \
		-o bin/certsync_client \
		github.com/icemarkom/certsync/client

server: deps *.go server/*
	go build \
		--ldflags "-s -w -X main.version=${version} -X main.gitCommit=${gitcommit} -X main.binaryName=certsync_server" \
		-o bin/certsync_server \
		server/*.go
