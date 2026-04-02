# Prompt Log — Лабораторная работа №10

Журнал выполненных промптов по лабораторной работе.

---

## Промпт 1 — Инициализация структуры проекта

**Дата:** 2026-04-02
**Промпт:** Инициализируй структуру проекта для Лабораторной работы №10, Вариант 1. Репозиторий уже создан: https://github.com/Dev66-66/LB10. Выполни: git remote add origin https://github.com/Dev66-66/LB10.git. Данные студента для README.md: ФИО: Фомичев Ярослав Николаевич, Группа: 221131, Вариант: 1, Номер лабораторной: 10. Создай следующую структуру: README.md с данными студента, описанием проекта, инструкцией по запуску обоих сервисов; PROMPT_LOG.md (пустой шаблон с заголовком); директорию go-service/ с Go-модулем; директорию python-service/ с requirements.txt. Предметная область: трекер тренировок. Поля модели Workout: id, name, type, duration, difficulty, calories_burned, created_at.
**Результат:** Добавлен remote origin. Созданы README.md с данными студента и инструкциями запуска, PROMPT_LOG.md, .gitignore, go-service/go.mod (модуль github.com/Dev66-66/LB10/go-service), python-service/requirements.txt. Сделан коммит "Initial project structure with README and module setup", выполнен push в main.

---

## Промпт 2 — Реализация Go-сервиса (Gin)

**Дата:** 2026-04-02
**Промпт:** Реализуй Go-сервис (Gin) внутри go-service/ с чистой архитектурой. Структура: cmd/server/main.go, internal/app/app.go, internal/handlers/http/workout_handler.go, internal/store/workout_store.go, internal/models/workout.go, internal/middleware/logger.go. Требования: модель Workout с полями id/name/type/duration/difficulty/calories_burned/created_at, in-memory хранилище с sync.RWMutex, GetAll() возвращает пустой слайс, эндпоинты GET/POST/GET/:id/DELETE/:id /workouts, валидация name через strings.TrimSpace, gin.New() + явный gin.Recovery(), middleware логгера key=value, порт 8080.
**Результат:** Реализованы все 6 файлов Go-сервиса. Модель с enum-типами WorkoutType/WorkoutDifficulty и методами IsValid(). Потокобезопасное in-memory хранилище на sync.RWMutex. Middleware логгера с форматом key=value. Обработчики с валидацией (name, type, difficulty, duration, calories_burned). App инициализирует gin.New() + Recovery() + Logger(). Добавлен gin v1.12.0. Сборка проверена (go build ./...). Два коммита, push выполнен.
