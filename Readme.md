# MarketFlow

**MarketFlow** — это высоконагруженное приложение, разработанное для обработки рыночных данных в реальном времени с использованием гексагональной архитектуры. Оно поддерживает два режима работы: **Live Mode** (работа с живыми биржами) и **Test Mode** (генерация синтетических данных). Приложение предоставляет REST API для получения информации о ценах и статистике.

---

## Порядок запуска
- docker-compose build
- docker-compose run loader (загружаем образы бирж)
- docker-compose up

## 📐 Архитектура

## MarketFlow использует **Hexagonal Architecture (Ports and Adapters)**:

- **Domain Layer** – бизнес-логика и модели.
- **Application Layer** – use-case'ы, связывающие бизнес-логику с адаптерами.
- **Adapters**:
  - Web Adapter – REST API (HTTP).
  - Storage Adapter – PostgreSQL.
  - Cache Adapter – Redis.
  - Exchange Adapter – Подключение к биржам и генератору.



## Структура проекта (Гексагональная архитектура)

```
├── cmd
│   ├── marketfuck
│   │   └── main.go
│   └── testgen
│       └── testgen.go
├── configs
│   └── app.yaml
├── deployments
│   ├── docker-compose.yml
│   ├── Dockerfile
│   └── sql
│       └── init.sql
├── go.mod
├── go.sum
├── internal
│   ├── adapter
│   │   ├── in
│   │   │   ├── exchange
│   │   │   │   └── live
│   │   │   │       └── connector.go
│   │   │   └── http
│   │   │       ├── handler
│   │   │       │   ├── all_handlers.go
│   │   │       │   ├── health_handler.go
│   │   │       │   ├── mode_handler.go
│   │   │       │   └── price_handler.go
│   │   │       ├── middleware
│   │   │       │   └── logger._middleware.go
│   │   │       ├── router
│   │   │       │   └── router.go
│   │   │       └── server.go
│   │   └── out_impl_for_port_out
│   │       ├── cache
│   │       │   └── redis
│   │       │       ├── cache.go
│   │       │       └── mapper.go
│   │       ├── exchange
│   │       │   ├── live
│   │       │   │   └── client.go
│   │       │   └── test
│   │       │       └── client.go
│   │       └── storage
│   │           └── postgres
│   │               ├── connectDB.go
│   │               ├── health_repo.go
│   │               └── pricePost.go
│   ├── application
│   │   ├── port
│   │   │   ├── in
│   │   │   │   ├── all_services.go
│   │   │   │   ├── exchange.go
│   │   │   │   ├── health_service.go
│   │   │   │   ├── mode_service.go
│   │   │   │   └── price_service.go
│   │   │   ├── logger.go
│   │   │   └── out
│   │   │       ├── cache.go
│   │   │       ├── exchange.go
│   │   │       └── storage.go
│   │   └── usecase_impl_for_port_in
│   │       ├── health_service.go
│   │       ├── mode_manager.go
│   │       ├── price_aggregator.go
│   │       └── price_fetcher.go
│   └── domain
│       ├── model
│       │   ├── exchange.go
│       │   ├── health.go
│       │   └── market.go
│       └── service
│           ├── market.go
│           └── mode.go
├── pkg
│   ├── concurrency
│   │   ├── fan_in.go
│   │   ├── fan_out.go
│   │   ├── gen_aggr.go
│   │   └── worker_pool.go
│   ├── config
│   │   ├── config.go
│   │   └── loader.go
│   ├── errors
│   │   └── errors.go
│   ├── logger
│   │   └── logger.go
│   ├── runner
│   │   └── runner.go
│   └── utils
│       ├── priceNameValid.go
│       └── serialize.go
├── Readme.md
├── setup_structure.sh
├── sources
│   ├── exchange1
│   │   └── exchange1_amd64.tar
│   ├── exchange2
│   │   └── exchange2_amd64.tar
│   └── exchange3
│       └── exchange3_amd64.tar
└── work.md
```


## Endpoints

GET /prices/latest/{symbol}– Получите последнюю цену за данный символ.   +++ -->

GET /prices/latest/{exchange}/{symbol}– Получите последнюю цену за данный символ от конкретной биржи. +++

GET /prices/highest/{symbol}– Получите самую высокую цену за определенный период. +++

GET /prices/highest/{exchange}/{symbol}– Получите самую высокую цену за определенный период от конкретной биржи. 

GET /prices/highest/{symbol}?period={duration} – Get the highest price within the last {duration} (e.g., the last 1s,  3s, 5s, 10s, 30s, 1m, 3m, 5m).

GET /prices/highest/{exchange}/{symbol}?period={duration} – Get the highest price within the last {duration} from a specific exchange.

GET /prices/lowest/{symbol}– Получите самую низкую цену за определенный период.

GET /prices/lowest/{exchange}/{symbol}– Получите самую низкую цену за определенный период от конкретной биржи.

GET /prices/lowest/{symbol}?period={duration}– Получите самую низкую цену в течение последнего {продления}.

GET /prices/lowest/{exchange}/{symbol}?period={duration} – Get the lowest price within the last {duration} from a specific exchange.

GET /prices/average/{symbol}– Получите среднюю цену за период.

GET /prices/average/{exchange}/{symbol}– Получите среднюю цену за определенный период с конкретной биржи.

GET /prices/average/{exchange}/{symbol}?period={duration} – Get the average price within the last {duration} from a specific exchange -


API режима данных

POST /mode/test – Switch to Test Mode (use generated data).

POST /mode/live – Switch to Live Mode (fetch data from provided programs).


Состояние системы

GET /health- Возвращает состояние системы (например, соединения, доступность Redis).


