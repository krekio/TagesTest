.PHONY: build run

# Имя бинарного файла
BINARY_NAME := tages-server

# Цель по умолчанию
default: build

# Сборка проекта
build:
	go build -o bin/$(BINARY_NAME) ./cmd/startserver

# Запуск сервера
run: build
	./bin/$(BINARY_NAME)

# Очистка
clean:
	rm -rf bin/