version=`date "+%Y%m%d%H%M%S"`
githash=`git rev-parse HEAD`

all: client server

docker_all: docker_server docker_client

client: *.go client/*
	go build \
  		--ldflags "-X main.version=${version} -X main.gitHash=${githash}" \
    	-o certsync_client \
		client/*.go

server:
	go build \
  	--ldflags "-X main.version ${version} -X main.gitHash ${githash}" \
    -o certsync_server \
		client/

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
