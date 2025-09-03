# Golang Task API

REST API для управления задачами (TODO-лист), реализованное на Go с использованием Fiber и PostgreSQL.

## Технологии

- **Go 1.24** - основной язык программирования
- **Fiber** - веб-фреймворк
- **PostgreSQL** - база данных
- **pgx** - драйвер для PostgreSQL
- **Docker** - контейнеризация

## Структура проекта

```
golang-task-api/
├── cmd/
│   └── main.go                 # Точка входа приложения
├── internal/
│   ├── config/                 # Конфигурация
│   ├── db/                     # Подключение к БД
│   ├── logger/                 # Логирование
│   ├── middleware/             # Middleware
│   ├── migrations/             # Миграции БД
│   ├── models/                 # Модели данных
│   ├── repository/             # Репозитории
│   ├── server/                 # HTTP сервер
│   ├── service/                # Бизнес-логика
│   └── startup/                # Инициализация приложения
├── Dockerfile                  # Docker образ
├── docker-compose.yml          # Docker Compose
└── README.md                   # Документация
```

## API Endpoints

### Создание задачи
```http
POST /tasks
Content-Type: application/json

{
  "title": "Название задачи",
  "description": "Описание задачи"
}
```

**Ответ:**
```json
{
  "id": 1,
  "title": "Название задачи",
  "description": "Описание задачи",
  "status": "new",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### Получение списка всех задач
```http
GET /tasks
```

**Ответ:**
```json
[
  {
    "id": 1,
    "title": "Название задачи",
    "description": "Описание задачи",
    "status": "new",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
]
```

### Получение задачи по ID
```http
GET /tasks/{id}
```

### Обновление задачи
```http
PUT /tasks/{id}
Content-Type: application/json

{
  "title": "Обновленное название",
  "description": "Обновленное описание",
  "status": "in_progress"
}
```

### Удаление задачи
```http
DELETE /tasks/{id}
```

## Статусы задач

- `new` - новая задача
- `in_progress` - в процессе выполнения
- `done` - выполнена

## Запуск

### С Docker (рекомендуется)

1. Запустите приложение с помощью Docker Compose:
   ```bash
   docker-compose up --build
   ```

2. Приложение будет доступно по адресу `http://localhost:8080`

### Локально

1. Установите PostgreSQL и создайте базу данных:
   ```sql
   CREATE DATABASE golangtaskapi;
   ```

2. Установите зависимости:
   ```bash
   go mod download
   ```

3. Запустите миграции:
   ```bash
   # Установите golang-migrate
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   
   # Запустите миграции
   migrate -path internal/migrations -database "postgres://postgres:postgres@localhost:5432/golangtaskapi?sslmode=disable" up
   ```

4. Запустите приложение:
   ```bash
   go run cmd/main.go
   ```

## Конфигурация

Переменные окружения:

- `RUN_ADDRESS` - адрес сервера (по умолчанию: `localhost:8080`)
- `DATABASE_URI` - URI подключения к PostgreSQL (по умолчанию: `postgres://postgres:postgres@localhost:5432/golangtaskapi?sslmode=disable`)

## Тестирование

Запуск всех тестов:
```bash
go test ./... -v
```

Запуск тестов по модулям:
```bash
go test ./internal/service/... -v
go test ./internal/server/... -v
go test ./internal/middleware/... -v
```

## Структура базы данных

```sql
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
```

## Примеры использования

### Создание задачи
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Изучить Go", "description": "Изучить основы языка Go"}'
```

### Получение всех задач
```bash
curl http://localhost:8080/tasks
```

### Обновление задачи
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "in_progress"}'
```

### Удаление задачи
```bash
curl -X DELETE http://localhost:8080/tasks/1
```
