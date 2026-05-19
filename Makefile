.PHONY: proto generate build test docker-up docker-down

PROTO_DIR := proto
PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

PROTOC ?= $(shell command -v protoc 2>/dev/null || echo ./.tools/protoc-bin)

proto:
	@test -x "$(PROTOC)" || { echo "protoc is required (install or run: curl protoc to .tools/)"; exit 1; }
	$(PROTOC) \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTO_DIR)/gen --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_DIR)/gen --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

generate: proto
	cd gateway && go generate ./...

build:
	cd user-service && go build -o ../bin/user-service ./cmd
	cd event-service && go build -o ../bin/event-service ./cmd
	cd ticket-service && go build -o ../bin/ticket-service ./cmd
	cd gateway && go build -o ../bin/gateway ./cmd

test:
	go work sync
	cd user-service && go test ./...
	cd event-service && go test ./...
	cd ticket-service && go test ./...
	cd gateway && go test ./...

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v
