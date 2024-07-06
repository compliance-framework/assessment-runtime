FROM golang:1.22 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /runtime

# Use a minimal base image to package the runtime binary
# Python so we can install azure cli
FROM mcr.microsoft.com/azure-cli:latest

WORKDIR /app

COPY --from=builder /runtime ./runtime

COPY test/config/docker/config.yml .
COPY test/config/assessments/ ./assessments/

# Run
CMD ["./runtime"]
