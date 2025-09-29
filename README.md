# Запуск без docker

Для запуска api нужен golang 1.25+. [Если нету go](https://go.dev/doc/install)
```bash
bash run.sh
```
> в run.sh за место go mod tidy, поставьте go mod download

**ИЛИ**

```bash
cd project
go mod download
go run main.go
```