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

.PHONY: client
client:
	go run client/cmd/*.go

.PHONY: proto
proto:
	protoc -I pkg/grpc/proto \
	  --go_out=./.. \
	  --go-grpc_out=./.. \
	  pkg/grpc/proto/*.proto