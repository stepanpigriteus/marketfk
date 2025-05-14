#!/bin/bash

# Базовая директория проекта
BASE_DIR="marketflow"

# Функция для создания папки и файла
create_file() {
  mkdir -p "$(dirname "$1")"
  touch "$1"
}

# Список файлов для создания
FILES=(
  "$BASE_DIR/cmd/marketflow/main.go"
  "$BASE_DIR/cmd/testgen/main.go"
  "$BASE_DIR/internal/domain/model/market.go"
  "$BASE_DIR/internal/domain/model/exchange.go"
  "$BASE_DIR/internal/domain/service/market.go"
  "$BASE_DIR/internal/application/usecase/price_fetcher.go"
  "$BASE_DIR/internal/application/usecase/price_aggregator.go"
  "$BASE_DIR/internal/application/usecase/mode_manager.go"
  "$BASE_DIR/internal/application/port/in/http.go"
  "$BASE_DIR/internal/application/port/in/exchange.go"
  "$BASE_DIR/internal/application/port/out/storage.go"
  "$BASE_DIR/internal/application/port/out/cache.go"
  "$BASE_DIR/internal/application/port/out/exchange.go"
  "$BASE_DIR/internal/adapter/in/http/handler/price.go"
  "$BASE_DIR/internal/adapter/in/http/handler/mode.go"
  "$BASE_DIR/internal/adapter/in/http/handler/health.go"
  "$BASE_DIR/internal/adapter/in/http/server.go"
  "$BASE_DIR/internal/adapter/in/http/router.go"
  "$BASE_DIR/internal/adapter/in/exchange/live/connector.go"
  "$BASE_DIR/internal/adapter/in/exchange/test/connector.go"
  "$BASE_DIR/internal/adapter/out/storage/postgres/repository.go"
  "$BASE_DIR/internal/adapter/out/storage/postgres/mapper.go"
  "$BASE_DIR/internal/adapter/out/cache/redis/cache.go"
  "$BASE_DIR/internal/adapter/out/cache/redis/mapper.go"
  "$BASE_DIR/internal/adapter/out/exchange/live/client.go"
  "$BASE_DIR/internal/adapter/out/exchange/test/client.go"
  "$BASE_DIR/pkg/config/config.go"
  "$BASE_DIR/pkg/config/loader.go"
  "$BASE_DIR/pkg/concurrency/fan_in.go"
  "$BASE_DIR/pkg/concurrency/fan_out.go"
  "$BASE_DIR/pkg/concurrency/worker_pool.go"
  "$BASE_DIR/pkg/logger/logger.go"
  "$BASE_DIR/configs/app.yaml"
  "$BASE_DIR/deployments/docker-compose.yml"
  "$BASE_DIR/deployments/Dockerfile"
)

# Создание файлов
for file in "${FILES[@]}"; do
  create_file "$file"
done

echo "Структура проекта успешно создана в папке '$BASE_DIR'"
