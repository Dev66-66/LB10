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

---

## Промпт 3 — Тесты для Go-сервиса

**Дата:** 2026-04-02
**Промпт:** Напиши тесты для Go-сервиса внутри go-service/. Юнит-тесты для хранилища: TestGetAll_EmptyReturnsSliceNotNil, TestCreate_ValidWorkout, TestCreate_EmptyNameRejected, TestCreate_WhitespaceNameRejected, TestGetByID_NotFound. Интеграционные тесты для HTTP-хендлеров: POST /workouts с валидным телом → 201, POST /workouts с пустым name → 400, GET /workouts → 200 массив не null, GET /workouts/:id несуществующий → 404, middleware логгера проверить поля key=value. Только стандартная библиотека + testify, реальное хранилище.
**Результат:** Добавлена валидация имени в store.Create (возвращает ErrInvalidName), обновлён обработчик. Написаны 5 юнит-тестов в workout_store_test.go и 5 интеграционных тестов в workout_handler_test.go (тест middleware захватывает os.Stdout через os.Pipe). Добавлен testify. Все 10 тестов проходят (go test ./...). 1 коммит, push выполнен.

---

## Промпт 4 — gRPC-сервер в Go-сервисе

**Дата:** 2026-04-02
**Промпт:** Реализуй gRPC-сервер в Go-сервисе. Создай proto/workout.proto в корне: сервис WorkoutService с методами GetWorkout, ListWorkouts, CreateWorkout; сообщение Workout с полями id/name/type/duration/difficulty/calories_burned/created_at. Сгенерируй Go-код из proto, добавь команду protoc в README. Реализуй gRPC-сервер в internal/grpc/workout_grpc_server.go. gRPC на порту 50051, HTTP на 8080, оба запускаются одновременно в main.go, переиспользуют одно in-memory хранилище.
**Результат:** Установлен protoc v27.3. Создан proto/workout.proto. Сгенерированы pb/workout.pb.go и pb/workout_grpc.pb.go. Реализован WorkoutGRPCServer с методами GetWorkout/ListWorkouts/CreateWorkout, коды gRPC-статусов (NotFound, InvalidArgument). App.New() принимает store как аргумент — оба сервера делят один экземпляр store. main.go запускает gRPC в горутине, HTTP блокирует. README обновлён командой protoc. Сборка и все тесты проходят. 2 коммита, push выполнен.

---

## Промпт 5 — JWT-аутентификация в Go-сервисе

**Дата:** 2026-04-03
**Промпт:** Добавь JWT-аутентификацию в Go-сервис. POST /auth/login принимает username/password, хардкодный пользователь admin/password123, возвращает {"token":"..."}. JWT: HS256, секрет из JWT_SECRET (по умолчанию "dev-secret"), TTL 24 часа, claims: username+exp. Все /workouts/* защитить JWT-middleware (Bearer-токен, 401 если невалидный). Написать 2 теста: валидный токен проходит, невалидный → 401. Библиотека: github.com/golang-jwt/jwt/v5.
**Результат:** Созданы auth_handler.go (POST /auth/login) и middleware/jwt.go (Bearer-валидация). RegisterRoutes изменён на gin.IRouter. app.go читает JWT_SECRET из env, монтирует /auth/login публично и /workouts/* через JWT middleware. Все 4 существующих теста обновлены с validBearerToken(). Добавлены TestJWTMiddleware_ValidToken_PassesThrough и TestJWTMiddleware_InvalidToken_Returns401. Все 12 тестов проходят. 2 коммита, push выполнен.

---

## Промпт 6 — Python-сервис на FastAPI

**Дата:** 2026-04-03
**Промпт:** Реализуй Python-сервис внутри python-service/ используя FastAPI. Архитектура: main.py, routers/workouts.py, routers/auth.py, routers/stats.py, services/go_workout_service.py (HTTP-клиент), services/grpc_workout_service.py (gRPC-клиент), models/workout.py, proto/ (сгенерированные стабы). Pydantic-модели с валидацией type/difficulty через Literal и name через field_validator. GoWorkoutService: asynchx, ConnectError→503, HTTPStatusError→проброс кода, JWT-аутентификация к Go. GrpcWorkoutService: grpc.aio. Эндпоинты: GET/POST /workouts, GET /workouts/{id}/grpc, GET /stats, POST /auth/token. Порт 8000.
**Результат:** Установлены fastapi 0.135.3, httpx 0.28.1, grpcio 1.80.0. Сгенерированы Python proto-стабы, исправлен импорт в workout_pb2_grpc.py. Реализованы все файлы сервиса. WorkoutCreate валидирует name (пустые/пробельные → ValueError), type и difficulty через Literal. GoWorkoutService кэширует JWT и переиспользует его. GrpcWorkoutService использует grpc.aio. README обновлён. 2 коммита, push выполнен.

---

## Промпт 7 — Тесты для Python-сервиса

**Дата:** 2026-04-03
**Промпт:** Напиши тесты для Python-сервиса используя pytest + respx. Файл: python-service/tests/test_workouts.py. Тесты: test_create_workout_valid (мок Go POST → 201), test_create_workout_empty_name (Pydantic ловит до HTTP), test_get_workouts (мок GET → список), test_go_service_unavailable (ConnectError → 503), test_stats_aggregation (мок списка → проверить агрегацию), test_jwt_forwarded (проверить заголовок Authorization). Используй pytest-asyncio для async-тестов.
**Результат:** Установлены pytest 9.0.2, pytest-asyncio 1.3.0, respx 0.22.0. Создан pytest.ini (asyncio_mode=auto, pythonpath=.). Написаны 6 тестов в tests/test_workouts.py. Используется `async with respx.mock() as mock:` (context manager, избегает проблем с декоратором + фикстурами). Фикстура reset_tokens сбрасывает кэш JWT всех сервисов перед каждым тестом. Фикстура client использует ASGITransport → не перехватывается respx. Все 6 тестов проходят. 1 коммит, push выполнен.

---

## Промпт 8 — Финализация проекта

**Дата:** 2026-04-03
**Промпт:** Финализируй проект. 1) Обнови README.md: данные студента, запуск Go-сервиса (go run ./cmd/server), запуск Python-сервиса (uvicorn main:app --reload), генерация proto, запуск тестов (go test ./... и pytest), переменные окружения (JWT_SECRET, GO_SERVICE_URL, GRPC_HOST), ASCII-диаграмма взаимодействия сервисов. 2) Убедись что PROMPT_LOG.md содержит записи для всех 8 промптов. 3) Запусти go vet ./... и исправь предупреждения. 4) Убедись что GetAll() возвращает [] а не null.
**Результат:** README.md полностью переписан: данные студента, ASCII-архитектурная диаграмма, таблицы эндпоинтов для обоих сервисов, таблица переменных окружения, инструкции запуска, тестирования и генерации proto. go vet ./... — без ошибок. GetAll() использует make([]models.Workout, 0, ...) → всегда возвращает JSON-массив []. Все тесты проходят: Go (12 тестов), Python (6 тестов). PROMPT_LOG.md содержит все 8 промптов. 1 коммит, push выполнен.

---

## Промпт 9 — Исправление замечаний code review

**Дата:** 2026-04-03
**Промпт:** Исправь все найденные проблемы и недочёты из code review. Оставляя атомарные коммиты.
**Результат:** Выполнено 4 атомарных коммита. (1) Рефакторинг Python-сервиса: единственный синглтон GoWorkoutService через services/dependencies.py, все роутеры используют Depends() вместо module-level инстансов, добавлен asyncio.Lock для token refresh (double-checked locking), lifespan закрывает httpx-клиент при завершении, тест reset_tokens упрощён до одного инстанса, env-переменные GO_SERVICE_URL/GO_USERNAME/GO_PASSWORD/GRPC_HOST читаются через os.getenv. (2) WorkoutResponse: type и difficulty — Literal вместо str; grpc_workout_service: bare except → grpc.aio.AioRpcError. (3) Go gRPC: CreateWorkout теперь валидирует type, difficulty, duration, calories_burned аналогично HTTP-хендлеру. (4) .gitignore: glob **.exe покрывает поддиректории; README: go run ./cmd/server, добавлены GO_USERNAME/GO_PASSWORD в таблицу env vars, удалена заметка о захардкоженных значениях; proto/__init__.py: добавлен поясняющий комментарий. Все тесты проходят (Go: 12, Python: 6). Push выполнен.
