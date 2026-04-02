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

### go-service (Go REST API)

**Стек:** Go, Gin, GORM, PostgreSQL

**Запуск:**

```bash
cd go-service
cp .env.example .env   # заполните переменные окружения
go mod download
go run ./cmd/main.go
```

API будет доступен на `http://localhost:8080`

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
