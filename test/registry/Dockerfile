FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/busy ./test/plugins/busy.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/hello ./test/plugins/hello.go

FROM nginx:latest

RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /app/busy /usr/share/nginx/html/busy/1.0.0/busy
COPY --from=builder /app/hello /usr/share/nginx/html/hello/1.0.0/hello

RUN chmod 755 /usr/share/nginx/html/busy/1.0.0/busy
RUN chmod 755 /usr/share/nginx/html/hello/1.0.0/hello

EXPOSE 80
