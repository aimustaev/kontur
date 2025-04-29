# Service Gateway

Сервис для обработки новых сообщений через gRPC.

## Требования

- Go 1.21 или выше
- Protocol Buffers (protoc)
- Go plugins for protoc:
  - protoc-gen-go
  - protoc-gen-go-grpc

## Установка

1. Установите зависимости:
```bash
go mod download
```

2. Установите protoc и плагины:
```bash
# Установка protoc
brew install protobuf

# Установка Go плагинов
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Сборка и запуск

1. Сгенерируйте proto файлы:
```bash
make proto
```

2. Соберите проект:
```bash
make build
```

3. Запустите сервис:
```bash
make run
```

## API

Сервис предоставляет gRPC API для обработки новых сообщений:

- `HandleNewMessage` - потоковый RPC метод для обработки новых сообщений 