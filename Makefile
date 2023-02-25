version=`date "+%Y%m%d%H%M%S"`
gitcommit=`git rev-parse --short HEAD`

all: certsync_client certsync_server

#
# Standalone binaries
#
certsync_client: *.go client/*
	go build \
		--ldflags "-s -w -X main.version=${version} -X main.gitCommit=${gitcommit} -X main.binaryName=certsync_client" \
		-o certsync_client \
		client/*.go

certsync_server: *.go server/*
	go build \
		--ldflags "-s -w -X main.version=${version} -X main.gitCommit=${gitcommit} -X main.binaryName=certsync_server" \
		-o certsync_server \
		server/*.go
