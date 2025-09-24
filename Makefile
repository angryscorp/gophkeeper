SERVER_NAME  = gophkeeper-server
COMPOSE_FILE = server/docker-compose.yml
CLIENT_NAME  = gophkeeper-client
CLIENT_SRC   = client/cmd
CLIENT_BIN   = client/bin

.PHONY: start
start:
	docker compose -f $(COMPOSE_FILE) up --build -d

.PHONY: logs
logs:
	docker compose -f $(COMPOSE_FILE) logs -f $(SERVER_NAME)

.PHONY: stop
stop:
	docker compose -f $(COMPOSE_FILE) down

.PHONY: status
status:
	docker compose -f $(COMPOSE_FILE) ps

.PHONY: restart
restart: stop start

.PHONY: debug
debug:
	GOOS=darwin GOARCH=arm64 go build -o $(CLIENT_BIN)/$(CLIENT_NAME)-darwin-arm64 $(CLIENT_SRC)/*.go
	 ./client/bin/gophkeeper-client-darwin-arm64

.PHONY: proto
proto:
	protoc -I pkg/grpc/proto \
	  --go_out=./.. \
	  --go-grpc_out=./.. \
	  pkg/grpc/proto/*.proto

.PHONY: build
build:
	@echo "Building client binaries..."
	@mkdir -p $(CLIENT_BIN)

	# Linux amd64
	GOOS=linux GOARCH=amd64 go build -o $(CLIENT_BIN)/$(CLIENT_NAME)-linux-amd64 $(CLIENT_SRC)/*.go

	# Linux arm64
	GOOS=linux GOARCH=arm64 go build -o $(CLIENT_BIN)/$(CLIENT_NAME)-linux-arm64 $(CLIENT_SRC)/*.go

	# Windows amd64
	GOOS=windows GOARCH=amd64 go build -o $(CLIENT_BIN)/$(CLIENT_NAME)-windows-amd64.exe $(CLIENT_SRC)/*.go

	# Windows arm64
	GOOS=windows GOARCH=arm64 go build -o $(CLIENT_BIN)/$(CLIENT_NAME)-windows-arm64.exe $(CLIENT_SRC)/*.go

	# macOS amd64 (Intel)
	GOOS=darwin GOARCH=amd64 go build -o $(CLIENT_BIN)/$(CLIENT_NAME)-darwin-amd64 $(CLIENT_SRC)/*.go

	# macOS arm64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -o $(CLIENT_BIN)/$(CLIENT_NAME)-darwin-arm64 $(CLIENT_SRC)/*.go

	@echo "Build completed. Binaries are in $(CLIENT_BIN)/"

.PHONY: gen-keys
gen-keys:
	@openssl genpkey -algorithm Ed25519 -out server/keys/ed25519.key
	@openssl pkey -in server/keys/ed25519.key -pubout -out server/keys/ed25519.pub
	@echo "Keys generated"
