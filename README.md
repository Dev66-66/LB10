# Лабораторная работа №10

| Поле | Значение |
|------|----------|
| ФИО | Фомичев Ярослав Николаевич |
| Группа | 221131 |
| Вариант | 1 |
| Номер лабораторной | 10 |

## Описание проекта

Трекер тренировок — система из двух микросервисов для управления данными о тренировках.

### Модель Workout

| Поле | Тип | Описание |
|------|-----|----------|
| id | uuid/int | Уникальный идентификатор |
| name | string | Название тренировки |
| type | enum | Тип: `cardio`, `strength`, `flexibility` |
| duration | int | Длительность в минутах |
| difficulty | enum | Сложность: `easy`, `medium`, `hard` |
| calories_burned | int | Сожжённые калории |
| created_at | datetime | Дата и время создания |

---

## Сервисы

### go-service (Go REST API + gRPC)

**Стек:** Go, Gin, gRPC, in-memory store

**Запуск:**

```bash
cd go-service
go mod download
go run ./cmd/server/main.go
```

- HTTP REST API: `http://localhost:8080`
- gRPC сервер: `localhost:50051`

#### Генерация кода из proto

Требования: [protoc ≥ 27](https://github.com/protocolbuffers/protobuf/releases), плагины:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Команда генерации (выполнять из корня репозитория):

```bash
protoc \
  --proto_path=proto \
  --go_out=go-service/internal/grpc/pb \
  --go_opt=paths=source_relative \
  --go-grpc_out=go-service/internal/grpc/pb \
  --go-grpc_opt=paths=source_relative \
  proto/workout.proto
```

---

### python-service (Python REST API)

**Стек:** Python, FastAPI, SQLAlchemy, PostgreSQL

**Запуск:**

```bash
cd python-service
python -m venv venv
source venv/bin/activate   # Windows: venv\Scripts\activate
pip install -r requirements.txt
cp .env.example .env       # заполните переменные окружения
uvicorn main:app --reload --port 8000
```

API будет доступен на `http://localhost:8000`

---

## Требования

- Go 1.21+
- Python 3.11+
- PostgreSQL 15+
