#!/bin/bash

cd "$(dirname "$0")/project"
echo "Установка и проверка зависимостей..."
go mod download
go mod tidy

echo "Запуск приложения..."
go run main.go