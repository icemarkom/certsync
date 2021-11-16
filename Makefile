all: docker_all

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