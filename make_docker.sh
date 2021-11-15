#!/bin/bash

docker buildx build \
  --target server \
  -t icemarkom/certsync-server \
  --platform=linux/amd64,linux/arm64,linux/arm \
  --push \
  .

docker buildx build \
  --target client \
  -t icemarkom/certsync-server \
  --platform=linux/amd64,linux/arm64,linux/arm \
  --push \
  .
