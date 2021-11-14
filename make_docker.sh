#!/bin/bash

docker buildx build \
  -f server.Dockerfile \
  -t icemarkom/certsync-server \
  --platform=linux/amd64,linux/arm64,linux/arm \
  --push \
  .

docker buildx build \
  -f client.Dockerfile \
  -t icemarkom/certsync-server \
  --platform=linux/amd64,linux/arm64,linux/arm \
  --push \
  .
