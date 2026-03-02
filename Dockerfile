# Используем Go 1.26 на Alpine для лёгкости
FROM golang:1.26-alpine AS build

# Рабочая директория внутри контейнера
WORKDIR /app

# Копируем только go.mod и go.sum, чтобы кэшировать зависимости
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN go build -o timo ./cmd/timo

# --- runtime stage ---
FROM alpine:latest

# Создаём рабочую папку
WORKDIR /app

# Копируем бинарник из билд стадии
COPY --from=build /app/timo .

# Экспортируем порт
EXPOSE 8080

# Запуск контейнера: только HTTP сервер
# Здесь мы передаем флаг --http, который мы добавим в main.go
CMD ["./timo", "--http"]