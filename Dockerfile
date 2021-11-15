FROM golang:alpine AS builder
# Builder Image
WORKDIR /certsync
COPY . .
RUN GOOS="linux" GOARCH=$(uname -m | sed -e "s/aarch64/arm64/" -e "s/x86_64/amd64/" -e "s/armv7l/arm/") go build -o certsync_server server/*.go
RUN GOOS="linux" GOARCH=$(uname -m | sed -e "s/aarch64/arm64/" -e "s/x86_64/amd64/" -e "s/armv7l/arm/") go build -o certsync_client client/*.go

# Server Image
FROM alpine:latest AS server
RUN apk --no-cache add ca-certificates
# RUN apt-get update && apt-get install --yes ca-certificates
COPY --from=builder /certsync/certsync_server /usr/local/bin/certsync_server
ENTRYPOINT ["/usr/local/bin/certsync_server"]
CMD ["--help"]

# Client Image
FROM alpine:latest AS client
RUN apk --no-cache add ca-certificates
# RUN apt-get update && apt-get install --yes ca-certificates
COPY --from=builder /certsync/certsync_client /usr/local/bin/certsync_client
ENTRYPOINT ["/usr/local/bin/certsync_client"]
CMD ["--help"]
