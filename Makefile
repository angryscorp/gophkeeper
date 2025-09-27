SERVER_NAME  = gophkeeper-server
COMPOSE_FILE = server/docker-compose.yml

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
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -tags sqlcipher -o client/bin/debug-client-app client/cmd/*.go
	./client/bin/debug-client-app

.PHONY: proto
proto:
	protoc -I pkg/grpc/proto \
	  --go_out=./.. \
	  --go-grpc_out=./.. \
	  pkg/grpc/proto/*.proto

.PHONY: gen-keys
gen-keys:
	@openssl genpkey -algorithm Ed25519 -out server/keys/ed25519.key
	@openssl pkey -in server/keys/ed25519.key -pubout -out server/keys/ed25519.pub
	@echo "Keys generated in ./server/keys/"
