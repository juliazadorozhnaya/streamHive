APP_NAME=streamhive
PROTO_DIR=internal/delivery/grpc/pb
PROTO_FILE=$(PROTO_DIR)/streamhive.proto
MAIN=cmd/streamhive/main.go

.PHONY: build run clean proto

build:
	go build -o bin/$(APP_NAME) $(MAIN)

run:
	go run $(MAIN)

clean:
	rm -rf bin/

proto:
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(PROTO_DIR) \
	  --go-grpc_out=$(PROTO_DIR) \
	  $(PROTO_FILE)
