#!/bin/bash

# Функция для создания директорий и файлов
create_structure() {
    # Создание директорий
    mkdir -p \
    cmd/marketflow \
    config \
    internal/domain/model \
    internal/domain/service \
    internal/application/usecase \
    internal/application/port/in \
    internal/application/port/out \
    internal/adapter/in/api/handler \
    internal/adapter/in/api/middleware \
    internal/adapter/out/cache \
    internal/adapter/out/storage \
    internal/adapter/out/exchange/live \
    internal/adapter/out/exchange/test \
    pkg/logger \
    pkg/utils \
    scripts

    # Создание файлов внутри соответствующих папок
    touch cmd/marketflow/main.go
    touch config/config.go
    touch config/config.yaml
    touch internal/domain/model/price.go
    touch internal/domain/model/statistics.go
    touch internal/domain/service/market_service.go
    touch internal/application/usecase/price_usecase.go
    touch internal/application/usecase/mode_usecase.go
    touch internal/application/port/in/api_port.go
    touch internal/application/port/out/cache_port.go
    touch internal/application/port/out/storage_port.go
    touch internal/application/port/out/exchange_port.go
    touch internal/adapter/in/api/handler/price_handler.go
    touch internal/adapter/in/api/handler/mode_handler.go
    touch internal/adapter/in/api/handler/health_handler.go
    touch internal/adapter/in/api/middleware/logging.go
    touch internal/adapter/in/api/middleware/recovery.go
    touch internal/adapter/in/api/router.go
    touch internal/adapter/out/cache/redis_adapter.go
    touch internal/adapter/out/storage/postgres_adapter.go
    touch internal/adapter/out/exchange/live/client.go
    touch internal/adapter/out/exchange/live/exchange_adapter.go
    touch internal/adapter/out/exchange/test/generator.go
    touch pkg/logger/logger.go
    touch pkg/utils/shutdown.go
    touch pkg/utils/helpers.go
    touch scripts/start_exchanges.sh
    touch scripts/init_db.sql
}

# Запуск функции
create_structure

# Вывод сообщения об успешном создании
echo "Структура папок и файлов успешно создана!"
