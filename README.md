# Лабораторная работа №10

| Поле | Значение |
|------|----------|
| ФИО | Фомичев Ярослав Николаевич |
| Группа | 221131 |
| Вариант | 1 |
| Номер лабораторной | 10 |

---

## Описание проекта

**Трекер тренировок** — система из двух микросервисов:

- **go-service** — основной REST API + gRPC-сервер (Go, Gin)
- **python-service** — прокси и агрегатор (Python, FastAPI), использует Go-сервис через HTTP и gRPC

### Модель Workout

| Поле | Тип | Описание |
|------|-----|----------|
| id | int | Уникальный идентификатор |
| name | string | Название тренировки |
| type | enum | `cardio` / `strength` / `flexibility` |
| duration | int | Длительность в минутах |
| difficulty | enum | `easy` / `medium` / `hard` |
| calories_burned | int | Сожжённые калории |
| created_at | datetime | Дата и время создания (UTC) |

---

## Архитектура

```
┌─────────────────────────────────────────────────────────┐
│                        Клиент                           │
└────────────┬──────────────────────────┬─────────────────┘
             │ HTTP :8000               │ HTTP :8000
             ▼                          ▼
┌────────────────────────┐   ┌──────────────────────────┐
│    python-service      │   │      go-service           │
│    FastAPI :8000       │   │      Gin REST :8080        │
│                        │   │      gRPC     :50051       │
│  POST /auth/token ─────┼──►│  POST /auth/login         │
│  GET  /workouts   ─────┼──►│  GET  /workouts  (JWT)    │
│  POST /workouts   ─────┼──►│  POST /workouts  (JWT)    │
│  GET  /workouts/{id}   │   │  GET  /workouts/:id (JWT) │
│    /grpc ──────────────┼──►│  gRPC GetWorkout          │
│  GET  /stats      ─────┼──►│  GET  /workouts  (JWT)    │
└────────────────────────┘   └──────────────────────────┘
        httpx (HTTP)                in-memory store
        grpc.aio (gRPC)             sync.RWMutex
```

### Эндпоинты Go-сервиса

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| POST | `/auth/login` | — | Получить JWT |
| GET | `/workouts` | Bearer | Список тренировок |
| POST | `/workouts` | Bearer | Создать тренировку |
| GET | `/workouts/:id` | Bearer | Получить по ID |
| DELETE | `/workouts/:id` | Bearer | Удалить |
| gRPC | `GetWorkout` | — | Получить по ID |
| gRPC | `ListWorkouts` | — | Список |
| gRPC | `CreateWorkout` | — | Создать |

### Эндпоинты Python-сервиса

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/auth/token` | Получить JWT через Go |
| GET | `/workouts` | Список (проксирует в Go) |
| POST | `/workouts` | Создать (проксирует в Go) |
| GET | `/workouts/{id}/grpc` | Получить через gRPC |
| GET | `/stats` | Агрегация по type и difficulty |

---

## Переменные окружения

| Переменная | Сервис | Значение по умолчанию | Описание |
|------------|--------|-----------------------|----------|
| `JWT_SECRET` | go-service | `dev-secret` | Секрет подписи JWT (HS256) |
| `GO_SERVICE_URL` | python-service | `http://localhost:8080` | Базовый URL Go HTTP API |
| `GRPC_HOST` | python-service | `localhost:50051` | Адрес Go gRPC сервера |

> В текущей реализации `GO_SERVICE_URL` и `GRPC_HOST` захардкожены в сервисных классах.
> Для смены адресов отредактируйте `services/go_workout_service.py` и `services/grpc_workout_service.py`.

---

## Запуск

### Требования

- Go 1.21+
- Python 3.11+

### Go-сервис

```bash
cd go-service
# Опционально: задать секрет JWT
export JWT_SECRET=my-secret
go run ./cmd/server/main.go
```

- REST API: `http://localhost:8080`
- gRPC: `localhost:50051`

### Python-сервис

> Go-сервис должен быть запущен перед стартом Python-сервиса.

```bash
cd python-service
pip install -r requirements.txt
uvicorn main:app --reload --port 8000
```

- API: `http://localhost:8000`
- Документация (Swagger): `http://localhost:8000/docs`

---

## Тесты

### Go

```bash
cd go-service
go test ./...
```

### Python

```bash
cd python-service
pytest
```

---

## Генерация proto-кода

### Go (из корня репозитория)

Установить плагины:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Сгенерировать:

```bash
protoc \
  --proto_path=proto \
  --go_out=go-service/internal/grpc/pb \
  --go_opt=paths=source_relative \
  --go-grpc_out=go-service/internal/grpc/pb \
  --go-grpc_opt=paths=source_relative \
  proto/workout.proto
```

### Python (из директории python-service)

```bash
python -m grpc_tools.protoc \
  --proto_path=../proto \
  --python_out=proto \
  --grpc_python_out=proto \
  ../proto/workout.proto
# Заменить в proto/workout_pb2_grpc.py:
#   import workout_pb2  →  from proto import workout_pb2
```
