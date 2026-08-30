# Media Processing Platform

Асинхронная платформа обработки изображений на Go: REST API принимает задачи, RabbitMQ распределяет их между воркерами, результат возвращается клиенту. Построена по принципу event-driven архитектуры с полным разделением API-шлюза и вычислительного воркера.

## Что это и зачем

Типичная задача медиа-пайплайна: пользователь загружает изображение и выбирает фильтр, сервис обрабатывает его асинхронно и отдаёт результат. Такая схема развязывает сервисы и не блокирует API тяжёлыми вычислениями:

- API (`server`) только принимает задачу и кладёт её в очередь;
- фоновый воркер (`processor`) забирает задачу, обрабатывает и сохраняет результат;
- клиент периодически опрашивает статус и забирает готовый результат.

Это позволяет масштабировать воркеры горизонтально независимо от API и переживать пиковые нагрузки без потери задач.

## Возможности

- Регистрация, вход и выход из аккаунта (bcrypt + токен-сессии)
- Создание задачи обработки изображения с выбором фильтра
- Асинхронная обработка через очередь RabbitMQ (protobuf-сообщения)
- Плагинные фильтры: негатив, размытие, отражение по горизонтали, резкость
- Получение статуса и результата задачи
- Метрики воркера: время обработки и количество задач (Prometheus + Grafana)
- Интеграционные тесты полного жизненного цикла
- CI/CD через GitHub Actions

## Архитектура

```
            POST /task
┌──────────┐           ┌────────────────┐   RabbitMQ    ┌──────────────┐
│  Client  │──────────>│  server (API)  │──────────────>│  processor   │
└──────────┘           │  :8000         │  task_queue   │  worker :2112│
       ▲               └────────────────┘               └──────┬───────┘
       │  GET /status/{id}                                     │
       │  GET /result/{id}                                     ▼
       │                                                  PostgreSQL
       └───────────────────────────────────────────────────────┘
                      ┌─────────┐   ┌──────────────────────────┐
                      │  Redis  │   │ Prometheus / Grafana     │
                      │ сессии  │   │ метрики воркера          │
                      └─────────┘   └──────────────────────────┘
```

Поток задачи:

1. `server` сохраняет задачу в PostgreSQL (статус `in_progress`) и публикует protobuf-сообщение в очередь `task_queue`;
2. `processor` потребляет сообщение, декодирует изображение, применяет фильтр;
3. воркер сохраняет результат (base64 data-URI) и статус `ready`/`failed`;
4. клиент опрашивает `GET /status/{id}`, затем получает результат через `GET /result/{id}`.

Оба сервиса построены по чистой архитектуре с разделением слоёв `delivery → service → repository → infrastructure`.

## Компоненты

| Компонент | Роль |
|---|---|
| `server` | HTTP API: аутентификация, создание задач, статусы и результаты |
| `processor` | Фоновый воркер: обработка изображений, метрики |
| `pkg/proto` | Общая protobuf-схема сообщений между сервисами |
| PostgreSQL | Хранение пользователей, задач, результатов |
| Redis | Хранение токен-сессий (TTL) |
| RabbitMQ | Очередь задач, развязка сервисов |
| Prometheus / Grafana | Сбор и визуализация метрик воркера |
| tests | Модульные/Интеграционные/E2E тесты (pytest) |

## Ключевые решения

- **Асинхронность через брокер сообщений** - API не выполняет тяжёлую обработку, воркер масштабируется независимо.
- **protobuf вместо JSON в шине** - компактная и быстрая сериализация бинарных данных между сервисами.
- **Плагинная система фильтров** - единый интерфейс `Filter` + реестр; новый фильтр добавляется одной строкой без правки ядра.
- **Ручной ack/nack в консьюмере** - гарантирует обработку сообщения и предсказуемое поведение при ошибках.
- **Устойчивый к обрывам коннект к RabbitMQ** - экспоненциальный бэкофф + фоновый reconnect-цикл.
- **Чистая архитектура** - изоляция слоёв, внедрение зависимостей через конструкторы, compile-time проверка интерфейсов репозиториев.
- **Multi-stage Docker-сборка** - минимальный runtime-образ (debian-slim) без инструментов сборки.

## Технологический стек

- **Язык:** Go 1.26 (горутины, каналы, context)
- **HTTP:** chi v5 (роутер, мидлвари), swaggo/swag (Swagger), slog-chi
- **БД:** PostgreSQL 18 + GORM v2 (JSONB-сериализация, auto-migration)
- **Кэш:** Redis 7 (сессии с TTL)
- **Брокер:** RabbitMQ 3 (AMQP), `amqp091-go` (persistent-доставка, ack/nack)
- **Сериализация:** Protocol Buffers (protoc-gen-go, google.protobuf.Struct)
- **Обработка изображений:** `image`, `image/color`, `disintegration/imaging` (base64 ↔ PNG)
- **Безопасность:** bcrypt (хеширование паролей)
- **Метрики:** Prometheus client_golang (histogram, counter), promhttp
- **Логирование:** log/slog (text/JSON по окружению)
- **Конфигурация:** cleanenv (YAML + переопределение env-переменными)
- **Инфраструктура:** Docker (multi-stage), docker-compose, GitHub Actions
- **Тесты:** pytest, requests (интеграционные, e2e-сценарии)

## Структура репозитория

```
media-processing-platform/
├── .github/workflows/ci.yml        # CI/CD (GitHub Actions)
├── deployments/                    # docker-compose и инфраструктура
│   ├── docker-compose.yml          # весь стек
│   ├── prometheus/prometheus.yml   # конфиг Prometheus
│   └── grafana/provisioning/       # авто-настройка Grafana
├── pkg/proto/                      # общая protobuf-схема
├── server/                         # HTTP API
│   ├── cmd/main.go
│   ├── config/config.yaml
│   ├── docs/                       # сгенерированный Swagger
│   └── internal/
│       ├── delivery/http/          # роуты, хендлеры, мидлвари, DTO
│       ├── service/                # бизнес-логика
│       ├── repository/             # интерфейсы + реализации (postgres/redis)
│       ├── domain/                 # доменные модели
│       └── infrastructure/         # postgres, redis, rabbitmq
├── processor/                      # воркер
│   ├── cmd/main.go
│   ├── config/config.yaml
│   └── internal/
│       ├── delivery/http/metrics/  # эндпоинт /metrics
│       ├── service/                # обработка задач
│       ├── filter/                 # интерфейс Filter + фильтры
│       ├── repository/             # интерфейс + реализация
│       ├── metrics/                # метрики Prometheus
│       └── infrastructure/         # rabbitmq, postgres
├── tests/tests.py                  # интеграционные тесты
├── static/sigma.png                # тестовое изображение
├── Makefile
├── .env.example
└── go.mod
```

## Что реализовано

**API (`server`)**

- `POST /register` - регистрация (bcrypt-хеширование пароля)
- `POST /login` - вход, возвращает `{"token": "<uuid-сессии>"}`
- `POST /logout` - выход, инвалидация сессии
- `POST /task` - создание задачи `{"filter": {"name": "...", "parameters": {...}}, "image": "<base64>"}`
- `GET /status/{task_id}` - статус: `in_progress` / `ready` / `failed`
- `GET /result/{task_id}` - результат (изображение, content-type из data-URI)
- `GET /swagger/*` - документация API

Все роуты задач защищены `Bearer-middleware`; `user_id`/`session_id` пробрасываются в контекст запроса.

**Воркер (`processor`)**

- Консьюмер очереди с ручным ack/nack и переподключением
- Декодирование base64 → `image.Image`, применение фильтра, обратная PNG-кодировка
- Плагинные фильтры: `negative`, `blur`, `flip_x`, `sharpen`
- Метрики: `processor_task_duration_seconds` (гистограмма по фильтру), `processor_tasks_total` (счётчик по фильтру/статусу)
- Эндпоинт `GET /metrics` на порту 2112

## Быстрый старт

Требования: Docker, Docker Compose, Make, Go 1.26 (для локальной разработки).

```bash
# 1. Создать .env из шаблона
cp .env.example .env

# 2. Собрать образы
make docker-compose-build

# 3. �-апустить весь стек (server, processor, postgres, rabbitmq, redis, prometheus, grafana)
make docker-compose-up

# 4. Остановить
make docker-compose-down
```

После старта:

| Сервис | Адрес |
|---|---|
| API | http://localhost:8000 |
| Swagger | http://localhost:8000/swagger/index.html |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3000 |
| RabbitMQ Management | http://localhost:15672 |

## Как пользоваться (инструкция)

```bash
# Регистрация
curl -X POST http://localhost:8000/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo","password":"secret"}'

# Вход - получите токен
curl -X POST http://localhost:8000/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo","password":"secret"}'
# => {"token":"<uuid>"}

# Создать задачу (изображение в base64, например из файла)
IMG=$(base64 -w0 static/sigma.png)
curl -X POST http://localhost:8000/task \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"filter\":{\"name\":\"negative\"},\"image\":\"$IMG\"}"
# => {"task_id":"<uuid>"}

# Получить статус
curl http://localhost:8000/status/<task_id> -H "Authorization: Bearer $TOKEN"
# => {"status":"ready"}

# Получить результат (сохранить изображение)
curl http://localhost:8000/result/<task_id> -H "Authorization: Bearer $TOKEN" \
  -o result.png
```

## Тестирование

Тесты поднимают изолированный стек (профиль `test` в compose) и проверяют весь жизненный цикл: регистрация → вход → создание задачи → ожидание `ready` → получение результата, плюс негативные сценарии (404, 401).

```bash
make docker-compose-test
```

CI (GitHub Actions) выполняет `go build`/`go vet`/`go test` и прогоняет интеграционные тесты; пароли для инфраструктуры передаются через GitHub Secrets.

## Решение проблем

**Тест падает: «task is still in progress!»**
�-адача не обрабатывается воркером. Проверьте: запущен ли `processor` (`docker compose ps`), читает ли он очередь (`docker logs deployments-processor-1`), и совпадает ли имя очереди у сервера и воркера (`RABBITMQ_QUEUE_NAME`).

**API не стартует: «connection refused» к RabbitMQ/Postgres**
Скорее всего, не совпадают пароли между `.env` и контейнерами. RabbitMQ берёт креды из `RABBITMQ_DEFAULT_USER/PASS` (см. `deployments/docker-compose.yml`) - они должны совпадать с `RABBITMQ_USER/RABBITMQ_PASSWORD` в `.env`.

**`/metrics` пустой**
Убедитесь, что воркер пересобран с метриками (`make docker-compose-build`) и порт 2112 опубликован (`PROCESSOR_PORT`).

**Grafana не видит Prometheus**
Проверьте датасорс в `deployments/grafana/provisioning/datasources/` - URL должен указывать на `http://prometheus:9090`.

## Лицензия

Pet-проект автора. Исходный код открыт.
