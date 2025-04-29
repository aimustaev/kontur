.PHONY: proto build run

proto:
	protoc --go_out=internal/generated --go_opt=paths=source_relative \
		--go-grpc_out=internal/generated --go-grpc_opt=paths=source_relative \
		proto/gateway.proto

build: proto
	go build -o bin/gateway main.go

run: build
	./bin/gateway 