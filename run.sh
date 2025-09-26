#!/bin/bash

cd "$(dirname "$0")/project"
echo "Проверка и установка зависимостей..."
go mod tidy

echo "Запуск приложения..."
go run main.go