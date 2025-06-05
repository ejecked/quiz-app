FROM golang:1.24-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы зависимостей и устанавливаем модули
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код проекта
COPY main.go ./
COPY handlers ./handlers
COPY models ./models
COPY db ./db

# Создаём папку для загрузок (например, изображений)
RUN mkdir -p uploads
COPY uploads ./uploads

# Если у тебя есть HTML/CSS/JS фронтенд
COPY public ./public

# Собираем бинарный файл
RUN go build -o main .

# Открываем порт
EXPOSE 3000

# Команда по умолчанию
CMD ["./main"]
