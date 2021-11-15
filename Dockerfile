FROM --platform=${BUILDPLATFORM} golang:alpine AS builder
# Builder Image
WORKDIR /certsync
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN apk --no-cache add ca-certificates
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o certsync_server server/*.go
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o certsync_client client/*.go

# Server Image
FROM --platform=${TARGETPLATFORM} alpine:latest AS server
# RUN apt-get update && apt-get install --yes ca-certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /certsync/certsync_server /
ENTRYPOINT ["certsync_server"]
CMD ["--help"]

# Client Image
FROM --platform=${TARGETPLATFORM} alpine:latest AS client
# RUN apt-get update && apt-get install --yes ca-certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /certsync/certsync_client /
ENTRYPOINT ["certsync_client"]
CMD ["--help"]
