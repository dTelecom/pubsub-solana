APP_NAME = pubsub
BUILD_DIR = ./bin

.DEFAULT_GOAL := build

build: ## Собрать бинарник
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go
