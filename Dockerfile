# Specify the version of Go to use
FROM golang:1.13 AS builder

# Install upx (upx.github.io) to compress the compiled action
RUN apt-get update && apt-get -y install upx

# Turn on Go modules support and disable CGO
ENV GO111MODULE=on CGO_ENABLED=0

# Copy all the files from the host into the container
COPY . /tmp/s3-proxy

WORKDIR /tmp/s3-proxy
# Compile the action - the added flags instruct Go to produce a
# standalone binary
RUN go build \
  -a \
  -trimpath \
  -ldflags "-s -w -extldflags '-static'" \
  -installsuffix cgo \
  -tags netgo \
  -o /bin/s3-proxy\
  .

# Strip any symbols - this is not a library
RUN strip /bin/s3-proxy

# Compress the compiled action
RUN upx -q -9 /bin/s3-proxy

# Step 2

FROM alpine:3.11

# Copy over SSL certificates from the first step - this is required
# if our code makes any outbound SSL connections because it contains
# the root CA bundle.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /bin/s3-proxy/ bin/s3-proxy
ENTRYPOINT ["/bin/s3-proxy"]
