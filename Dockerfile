# syntax = docker/dockerfile:1.0-experimental
# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.23 as builder

# Set the Current Working Directory inside the container
WORKDIR /foodfinder-api

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 go build -a -tags timetzdata -o ./deploy/foodfinder-api .

FROM scratch
# Create an empty world-writable /tmp folder so temp files can be used.
RUN --mount=from=busybox:1.34.0,src=/bin/,dst=/bin/ mkdir -m 1755 /tmp
# Copy the binary from the builder stage
COPY --from=builder /foodfinder-api/deploy/foodfinder-api /
# Copy the TLS root certs from the upstream image which already has ca-certificates:
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose port 8080 to the outside world
EXPOSE 8080
ENTRYPOINT [ "/foodfinder-api" ]
