
FROM golang:1.24-alpine AS builder
WORKDIR /app
# Устанавливаем git (необходим для загрузки зависимостей)
RUN apk update && apk add --no-cache git
# Копируем модули и скачиваем зависимости
COPY go.mod ./
RUN go mod download
 # Копируем исходный код
COPY . .
# Собираем сервис статически (CGO отключён, чтобы можно было использовать минимальный образ)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go
FROM alpine
COPY --from=builder /app/scripts/migrations /scripts/migrations
COPY --from=builder /app/main /main
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
CMD ["./main"]
