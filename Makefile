version=`date "+%Y%m%d%H%M%S"`
githash=`git rev-parse --short HEAD`

# Does not make Docker -- use docker_all
all: client server

#
# Standalone binaries
#
client: deps *.go client/*
	go build \
		--ldflags "-s -w -X main.version=${version} -X main.gitHash=${githash} -X main.binaryName=certsync_client" \
		-o bin/certsync_client \
		github.com/icemarkom/certsync/client

server: deps *.go server/*
	go build \
		--ldflags "-s -w -X main.version=${version} -X main.gitHash=${githash} -X main.binaryName=certsync_server" \
		-o bin/certsync_server \
		server/*.go

#
# Docker
#
docker_all: docker_server docker_client

docker_server:
	docker buildx build \
		--target server \
		-t icemarkom/certsync:server \
		--platform=linux/amd64,linux/arm64,linux/arm \
		--push \
		.

docker_client:
	docker buildx build \
		--target client \
		-t icemarkom/certsync:client \
		--platform=linux/amd64,linux/arm64,linux/arm \
		--push \
		.

deps: Makefile go.mod