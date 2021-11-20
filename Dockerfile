FROM golang:alpine AS builder
# Builder Image
WORKDIR /certsync
COPY . .
RUN apk --no-cache add ca-certificates
RUN make server
RUN make client

# Server Image
FROM alpine:latest AS server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /certsync/certsync_server /certsync:server
WORKDIR /
ENV PATH "/"
ENTRYPOINT ["certsync:server"]
CMD ["--help"]

# Client Image
FROM alpine:latest AS client
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /certsync/certsync_client /certsync:client
WORKDIR /
ENV PATH "/"
ENTRYPOINT ["certsync:client"]
CMD ["--help"]