FROM golang:alpine AS builder

# Builder image
WORKDIR /certsync
COPY . .
RUN GOOS="linux" GOARCH=$(uname -m | sed -e "s/aarch64/arm64/" -e "s/x86_64/amd64/" -e "s/armv7l/arm/") go build -o certsync_server server/*.go

# Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
# RUN apt-get update && apt-get install --yes ca-certificates
COPY --from=builder /certsync/certsync_server /usr/local/bin/certsync_server
ENTRYPOINT ["/usr/local/bin/certsync_server"]
CMD ["--help"]
