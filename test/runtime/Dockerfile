FROM golang:1.21 as builder

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
RUN CGO_ENABLED=0 GOOS=linux go build -o /busy ./test/plugins/busy.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /hello ./test/plugins/hello.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /azurecli ./test/plugins/azurecli.go

# Use a minimal base image to package the runtime binary
# Python so we can install azure cli
FROM mcr.microsoft.com/azure-cli:latest

WORKDIR /app

COPY --from=builder /runtime ./runtime
COPY --from=builder /busy ./plugins/busy/1.0.0/busy
COPY --from=builder /hello ./plugins/hello/1.0.0/hello
COPY --from=builder /azurecli ./plugins/azurecli/1.0.0/azurecli

COPY test/config/docker/config.yml .
COPY test/config/assessments/ ./assessments/

# Run
CMD ["./runtime"]
