FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/busy ./test/plugins/busy.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/hello ./test/plugins/hello.go

FROM ghcr.io/compliance-framework/nano-web:latest

COPY --from=builder /app/busy /public/busy/1.0.0/busy
COPY --from=builder /app/hello /public/hello/1.0.0/hello

RUN chmod 755 /public/busy/1.0.0/busy
RUN chmod 755 /public/hello/1.0.0/hello

