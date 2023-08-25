FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/sample ./tests/sampleplugin/main.go

FROM nginx:latest

RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /app/sample /usr/share/nginx/html/sample/1.0.0/sample

RUN chmod 755 /usr/share/nginx/html/sample/1.0.0/sample

EXPOSE 80
