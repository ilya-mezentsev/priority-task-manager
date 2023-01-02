PROJECT_NAME := priority-task-manager

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

SOURCE_DIR := $(ROOT_DIR)/source

# Точка входа у приложений должна быть одна и
# расположена должна быть в одинаковом месте, относительно корня приложения
ENTRYPOINT := ./cmd/main.go

TASK_MANAGER_DIR := $(SOURCE_DIR)/manager
TASK_MANAGER_CONFIGS_DIR := $(TASK_MANAGER_DIR)/configs
TASK_MANAGER_WEB_CONFIGS_FILE := $(TASK_MANAGER_CONFIGS_DIR)/web.json
TASK_MANAGER_PERMISSION_CONFIGS_FILE := $(TASK_MANAGER_CONFIGS_DIR)/permission.json
TASK_MANAGER_QUEUE_PRIORITIES_CONFIGS_FILE := $(TASK_MANAGER_CONFIGS_DIR)/queue_priorities.json

TASK_EXECUTOR_DIR := $(SOURCE_DIR)/executor
TASK_EXECUTOR_CONFIGS_DIR := $(TASK_EXECUTOR_DIR)/configs
TASK_EXECUTOR_WORKERS_POOL_CONFIGS_FILE := $(TASK_EXECUTOR_CONFIGS_DIR)/workers_pool.json
TASK_EXECUTOR_REDIS_CONFIGS_FILE := $(TASK_EXECUTOR_CONFIGS_DIR)/redis.json

TASK_STATUS_DIR := $(SOURCE_DIR)/status

TASK_METRICS_DIR := $(SOURCE_DIR)/metrics
TASK_METRICS_CONFIGS_DIR := $(TASK_METRICS_DIR)/configs
TASK_METRICS_WEB_CONFIGS_FILE := $(TASK_METRICS_CONFIGS_DIR)/web_metrics.json

SHARED_DIR := $(SOURCE_DIR)/shared
SHARED_CONFIGS_DIR := $(SHARED_DIR)/configs
DB_CONFIGS_FILE := $(SHARED_CONFIGS_DIR)/db.json
RABBITMQ_CONFIGS_FILE := $(SHARED_CONFIGS_DIR)/rabbitmq.json
BASIC_AUTH_CONFIGS_FILE := $(SHARED_CONFIGS_DIR)/basic_auth.json
MERGED_CONFIGS_FILE := $(SHARED_CONFIGS_DIR)/main.json

CONTAINERS_DIR := $(SOURCE_DIR)/containers
DOCKER_COMPOSE_ENTRYPOINT := $(CONTAINERS_DIR)/docker-compose.yaml

SCRIPTS_DIR := $(ROOT_DIR)/scripts
JSON_MERGE_SCRIPT := $(SCRIPTS_DIR)/json-merge.py

GOMODCACHE = $(SHARED_DIR)/pkg/mod

run: configs-merge
	cd $(APP_DIR) && GOMODCACHE=$(GOMODCACHE) go run $(ENTRYPOINT) $(MERGED_CONFIGS_FILE)

tidy:
	cd $(APP_DIR) && GOMODCACHE=$(GOMODCACHE) go mod tidy

test:
	cd $(APP_DIR) && GOMODCACHE=$(GOMODCACHE) go test -cover -p 1 ./... | { grep -v "no test files"; true; }

manager-run:
	@$(MAKE) --no-print-directory APP_DIR=$(TASK_MANAGER_DIR) run

manager-tidy:
	@$(MAKE) --no-print-directory APP_DIR=$(TASK_MANAGER_DIR) tidy

manager-test:
	@$(MAKE) --no-print-directory APP_DIR=$(TASK_MANAGER_DIR) test

executor-run:
	@$(MAKE) --no-print-directory APP_DIR=$(TASK_EXECUTOR_DIR) run

executor-tidy:
	@$(MAKE) --no-print-directory APP_DIR=$(TASK_EXECUTOR_DIR) tidy

status-run:
	@$(MAKE) --no-print-directory APP_DIR=$(TASK_STATUS_DIR) run

status-tidy:
	@$(MAKE) --no-print-directory APP_DIR=$(TASK_STATUS_DIR) tidy

metrics-run:
	@$(MAKE) --no-print-directory APP_DIR=$(TASK_METRICS_DIR) run

metrics-tidy:
	@$(MAKE) --no-print-directory APP_DIR=$(TASK_METRICS_DIR) tidy

all-tidy: manager-tidy executor-tidy status-tidy metrics-tidy

containers-run:
	 COMPOSE_PROJECT_NAME=$(PROJECT_NAME) docker-compose -f $(DOCKER_COMPOSE_ENTRYPOINT) up

containers-build:
	COMPOSE_PROJECT_NAME=$(PROJECT_NAME) docker-compose -f $(DOCKER_COMPOSE_ENTRYPOINT) build

containers-clean:
	 COMPOSE_PROJECT_NAME=$(PROJECT_NAME) docker-compose -f $(DOCKER_COMPOSE_ENTRYPOINT) down

configs-merge:
	# Мерджим конфиги в один файл
	python3 $(JSON_MERGE_SCRIPT) \
		$(SHARED_CONFIGS_DIR) \
		$(TASK_MANAGER_WEB_CONFIGS_FILE) \
		$(TASK_MANAGER_PERMISSION_CONFIGS_FILE) \
		$(TASK_MANAGER_QUEUE_PRIORITIES_CONFIGS_FILE) \
		$(TASK_EXECUTOR_WORKERS_POOL_CONFIGS_FILE) \
		$(TASK_EXECUTOR_REDIS_CONFIGS_FILE) \
		$(TASK_METRICS_WEB_CONFIGS_FILE) \
		$(DB_CONFIGS_FILE) \
		$(RABBITMQ_CONFIGS_FILE) \
		$(BASIC_AUTH_CONFIGS_FILE)
