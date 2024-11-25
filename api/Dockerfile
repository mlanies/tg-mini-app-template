# Этап сборки
FROM golang:1.19-alpine AS builder

# Установка рабочей директории
WORKDIR /build

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . ./

# Сборка
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o main ./main.go

# Финальный образ
FROM alpine:3.18
RUN apk add --no-cache bash curl
COPY --from=builder /build/main /main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 3000
ENTRYPOINT ["/main"]
