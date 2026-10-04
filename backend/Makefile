.PHONY: help fmt lint vet test build run clean all check

# Цветной вывод
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RED    := \033[0;31m
NC     := \033[0m

help: ## Показать команды
	@echo "$(GREEN)Доступные команды:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-10s$(NC) %s\n", $$1, $$2}'

fmt: ## Форматировать код (gofmt)
	@gofmt -w -s .
	@echo "$(GREEN)✅ Форматирование готово$(NC)"

imports: ## Организовать импорты (goimports)
	@goimports -w .
	@echo "$(GREEN)✅ Импорты организованы$(NC)"

vet: ## Проверить код (go vet)
	@go vet ./...
	@echo "$(GREEN)✅ go vet прошёл$(NC)"

lint: ## Линтить код (golangci-lint)
	@golangci-lint run
	@echo "$(GREEN)✅ Линт прошёл$(NC)"

test: ## Запустить тесты
	@go test -v ./...

build: ## Собрать бинарник
	@go build -o picmap-backend .
	@echo "$(GREEN)✅ Сборка готова$(NC)"

run: ## Запустить локально
	@go run main.go

clean: ## Удалить бинарник
	@rm -f picmap-backend
	@echo "$(GREEN)✅ Очищено$(NC)"

check: fmt imports vet ## Полная проверка (fmt + imports + vet)
	@echo "$(GREEN)✅ Проверка завершена$(NC)"

all: check build ## Полный цикл (check + build)
	@echo "$(GREEN)🚀 Всё готово!$(NC)"
