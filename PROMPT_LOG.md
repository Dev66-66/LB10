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
