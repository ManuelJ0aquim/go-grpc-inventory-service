# Nome do módulo Go (extraído do go.mod)
MODULE_NAME = github.com/ManuelJ0aquim/go-grpc-inventory-service

# Caminhos principais
PROTO_DIR = api/proto
PB_OUT_DIR = internal/infrastructure/grpc/pb
CMD_DIR = cmd/server/main.go

.PHONY: all help proto clean run docker-up docker-down tidy test

all: proto tidy build

## help: Exibe todos os comandos disponíveis no Makefile
help:
	@echo "Comandos disponíveis:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

## proto: Compila os arquivos .proto para Go (gRPC e Protobuf)
proto:
	@echo "==> Gerando códigos Go a partir do Protobuf..."
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(PB_OUT_DIR) --go_opt=module=$(MODULE_NAME)/$(PB_OUT_DIR) \
		--go-grpc_out=$(PB_OUT_DIR) --go-grpc_opt=module=$(MODULE_NAME)/$(PB_OUT_DIR) \
		$(PROTO_DIR)/inventory/v1/inventory.proto
	@echo "==> Código gRPC gerado com sucesso em $(PB_OUT_DIR)!"

## run: Sobe o banco Postgres via Docker e inicia o servidor gRPC
## run: Sobe o banco Postgres via Docker e inicia o servidor gRPC
run: docker-up
	@echo "==> Aguardando o PostgreSQL inicializar..."
	@until docker exec inventory_postgres pg_isready -U postgres > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "==> PostgreSQL pronto! Iniciando o servidor gRPC..."
	go run $(CMD_DIR)

## build: Compila o binário da aplicação na pasta bin/
build: proto
	@echo "==> Compilando o binário..."
	go build -o bin/inventory-service $(CMD_DIR)

## docker-up: Sobe o container do PostgreSQL em background
docker-up:
	@echo "==> Subindo a infraestrutura do PostgreSQL..."
	docker-compose up -d

## docker-down: Para e remove os containers e volumes da infraestrutura
docker-down:
	@echo "==> Parando a infraestrutura..."
	docker-compose down -v

## tidy: Organiza e baixa as dependências do Go
tidy:
	@echo "==> Atualizando dependências Go..."
	go mod tidy

## clean: Limpa o binário gerado e os arquivos compilados do Protobuf
clean:
	@echo "==> Limpando arquivos gerados..."
	rm -rf bin/
	rm -rf $(PB_OUT_DIR)/inventory