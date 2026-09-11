INFRA_DIR = ./infra

.PHONY: docker-up docker-down

docker-up:
	@echo "==> Subindo a infraestrutura..."
	docker-compose -f $(INFRA_DIR)/docker-compose.yml build --no-cache
	docker-compose -f $(INFRA_DIR)/docker-compose.yml up -d

docker-down:
	@echo "==> Parando a infraestrutura..."
	docker-compose -f $(INFRA_DIR)/docker-compose.yml down -v
