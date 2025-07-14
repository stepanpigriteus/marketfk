

# MarketFlow - InPRogress

Система обработки и агрегации биржевой информации.



## Структура проекта (Гексагональная архитектура)

```
marketflow/
├── cmd/                   # Точки входа в приложение (исполняемые файлы)
│   ├── marketflow/        # Основной бинарник приложения
│   │   └── main.go        # Инициализация и запуск приложения, выбор профиля, DI-контейнер
│   └── testgen/           # Генератор тестовых данных
│       └── main.go        # Логика генерации тестовых данных для эмуляции биржи
│
├── internal/              # Внутренний код приложения, не предназначенный для импорта
│   ├── domain/            # Доменный слой - бизнес-логика и правила
│   │   ├── model/         # Доменные модели - основные сущности бизнес-логики
│   │   │   ├── market.go  # Модели рыночных данных (Price, Pair, Order, etc.)
│   │   │   └── exchange.go # Модель биржи (Exchange, ExchangeInfo, etc.)
│   │   └── service/       # Доменные сервисы - основная бизнес-логика
│   │       └── market.go  # Сервис обработки рыночных данных (расчеты, анализ)
│   │
│   ├── application/       # Прикладной слой - реализует сценарии использования
│   │   ├── usecase/       # Реализация сценариев использования (бизнес-операций)
│   │   │   ├── price_fetcher.go   # Получение цен с бирж, реализация входного порта
│   │   │   ├── price_aggregator.go # Агрегация цен, реализация входного порта
│   │   │   └── mode_manager.go    # Управление режимами, реализация входного порта
│   │   └── port/          # Порты для адаптеров (интерфейсы)
│   │       ├── in/        # Входные порты (интерфейсы для входящих адаптеров)
│   │       │   ├── price_service.go  # Интерфейс PriceService
│   │       │   ├── mode_service.go   # Интерфейс ModeService
│   │       │   ├── health_service.go # Интерфейс HealthService
│   │       │   ├── http.go    # Интерфейсы для REST API
│   │       │   └── exchange.go # Интерфейсы для обработки данных с бирж
│   │       └── out/      # Выходные порты (интерфейсы для исходящих адаптеров)
│   │           ├── storage.go  # Интерфейс для хранилища (Repository)
│   │           ├── cache.go    # Интерфейс для кеша (CacheClient)
│   │           └── exchange.go # Интерфейс для получения данных с бирж (ExchangeClient)
│   │
│   └── adapter/          # Адаптеры - соединяют внешний мир с приложением
│       ├── in/           # Входящие адаптеры - преобразуют внешние запросы в вызовы приложения
│       │   ├── http/     # HTTP адаптер - REST API
│       │   │   ├── handler/ # Обработчики запросов
│       │   │   │   ├── price.go   # Обработчики для цен, использует PriceService
│       │   │   │   ├── mode.go    # Обработчики для режимов, использует ModeService
│       │   │   │   └── health.go  # Обработчик для проверки здоровья, использует HealthService
│       │   │   ├── server.go # HTTP сервер (конфигурация, middleware)
│       │   │   └── router.go # Роутер (маршрутизация запросов)
│       │   └── exchange/    # Адаптеры для бирж
│       │       ├── live/    # Адаптер для живых данных
│       │       │   └── connector.go # Коннектор к биржам (подписка на данные)
│       │       └── test/    # Адаптер для тестовых данных
│       │           └── connector.go # Коннектор к тест-генератору
│       └── out/          # Исходящие адаптеры - преобразуют вызовы приложения во внешние
│           ├── storage/  # Адаптер хранилища
│           │   └── postgres/ # PostgreSQL адаптер
│           │       ├── repository.go # Репозиторий (реализация порта storage.go)
│           │       └── mapper.go    # Маппер данных (DTO <-> Domain)
│           ├── cache/    # Адаптер кеша
│           │   └── redis/  # Redis адаптер
│           │       ├── cache.go  # Реализация кеша (реализация порта cache.go)
│           │       └── mapper.go # Маппер данных (DTO <-> Domain)
│           └── exchange/ # Адаптер для получения данных с бирж
│               ├── live/ # Живой обмен
│               │   └── client.go # Клиент для подключения (реализация порта exchange.go)
│               └── test/ # Тестовый обмен
│                   └── client.go # Клиент для подключения (реализация порта exchange.go)
│
├── pkg/                 # Общие пакеты, которые могут быть использованы другими проектами
│   ├── config/          # Конфигурация
│   │   ├── config.go    # Модель конфигурации (структуры для настроек)
│   │   └── loader.go    # Загрузчик конфигурации (из файлов, env)
│   ├── concurrency/     # Утилиты для конкурентности
│   │   ├── fan_in.go    # Fan-In шаблон (объединение потоков данных)
│   │   ├── fan_out.go   # Fan-Out шаблон (распределение потоков данных)
│   │   └── worker_pool.go # Worker Pool шаблон (пул обработчиков)
│   │   └── gen_aggr.go  # создание каналов, запуск воркеров и т.д.Извлечение данных из fan-in 
│   └── logger/          # Логирование
│       └── logger.go    # Настройка логгера (с уровнями, форматами)
│
├── deployments/         # Конфигурация для развертывания
│   ├── docker-compose.yml # Docker Compose для запуска всех компонентов
│   ├── Dockerfile       # Dockerfile для сборки образа приложения
│
│
├── test/                # Тесты
│   ├── integration/     # Интеграционные тесты
│   └── e2e/             # End-to-End тесты
├── setup_structure.sh
├── sources             # Исходники бирж
│   ├── exchange1
│   │   └── exchange1_amd64.tar
│   ├── exchange2
│   │   └── exchange2_amd64.tar
│   └── exchange3
│       └── exchange3_amd64.tar
├── .gitignore           # Файлы, игнорируемые Git
├── Makefile             # Makefile для автоматизации задач
├── README.md            # Описание проекта
├── go.mod               # Go модули
└── go.sum               # Хеши модулей
```



<!-- 
GET /prices/latest/{symbol}– Получите последнюю цену за данный символ.   +++ -->
<!-- 
GET /prices/latest/{exchange}/{symbol}– Получите последнюю цену за данный символ от конкретной биржи. +++ -->
<!-- 
GET /prices/highest/{symbol}– Получите самую высокую цену за определенный период. +++ -->
<!-- 
GET /prices/highest/{exchange}/{symbol}– Получите самую высокую цену за определенный период от конкретной биржи.  -->

<!-- GET /prices/highest/{symbol}?period={duration} – Get the highest price within the last {duration} (e.g., the last 1s,  3s, 5s, 10s, 30s, 1m, 3m, 5m). -->

<!-- GET /prices/highest/{exchange}/{symbol}?period={duration} – Get the highest price within the last {duration} from a specific exchange. -->

<!-- GET /prices/lowest/{symbol}– Получите самую низкую цену за определенный период. -->

<!-- GET /prices/lowest/{exchange}/{symbol}– Получите самую низкую цену за определенный период от конкретной биржи. -->

<!-- GET /prices/lowest/{symbol}?period={duration}– Получите самую низкую цену в течение последнего {продления}. -->

<!-- GET /prices/lowest/{exchange}/{symbol}?period={duration} – Get the lowest price within the last {duration} from a specific exchange. -->

<!-- GET /prices/average/{symbol}– Получите среднюю цену за период. -->

<!-- GET /prices/average/{exchange}/{symbol}– Получите среднюю цену за определенный период с конкретной биржи.

GET /prices/average/{exchange}/{symbol}?period={duration} – Get the average price within the last {duration} from a specific exchange -->


API режима данных

POST /mode/test – Switch to Test Mode (use generated data).

POST /mode/live – Switch to Live Mode (fetch data from provided programs).


Состояние системы

GET /health- Возвращает состояние системы (например, соединения, доступность Redis).