# SACAS - Smart Automated Code Analysis System

Система автоматического анализа кода с использованием AI.

## 🚀 Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- OpenAI API ключ

### Настройка

1. **Клонируйте репозиторий:**
   ```bash
   git clone <repository-url>
   cd SACAS
   ```

2. **Настройте переменные окружения:**
   ```bash
   cp .env.example .env
   ```
   
   Отредактируйте `.env` файл и добавьте ваш OpenAI API ключ:
   ```
   OPENAI_API_KEY=your_actual_api_key_here
   ```

3. **Запустите приложение:**
   ```bash
   docker compose up -d
   ```

4. **Проверьте что все работает:**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Health check: http://localhost:8080/health

5. **Запустите тесты (опционально):**
   ```bash
   cd backend
   make test
   ```

## 🏗️ Архитектура

### Backend (Go)
```
backend/
├── cmd/server/main.go           # Точка входа
├── internal/
│   ├── config/                  # Конфигурация
│   ├── models/                  # Модели данных
│   ├── repositories/            # Слой данных
│   ├── services/                # Бизнес-логика
│   ├── handlers/                # HTTP обработчики
│   └── database/                # Подключение к БД
```

### Frontend (React)
```
frontend/
├── src/
│   ├── app/                     # Основное приложение
│   ├── features/                # Функциональные компоненты
│   ├── pages/                   # Страницы
│   ├── shared/                  # Общие компоненты
│   └── widgets/                 # Виджеты
```

## 📝 API Endpoints

### RESTful API

- `POST /api/submissions` - Создать новую проверку кода
- `GET /api/submissions` - Получить список всех проверок
- `GET /api/submissions/:id` - Получить конкретную проверку
- `DELETE /api/submissions/:id` - Удалить проверку
- `GET /health` - Проверка состояния сервиса

### Deprecated (для обратной совместимости)
- `POST /api/submit` - Отправить код на проверку

## 🔧 Разработка

### Запуск в режиме разработки

```bash
# Только база данных
docker compose up postgres -d

# Backend (из папки backend)
go run cmd/server/main.go

# Frontend (из папки frontend)
npm start

# Запуск всего приложения
docker compose up -d
```

### Остановка

```bash
docker compose down

# С удалением volumes (очистка БД)
docker compose down -v
```

## 🧪 Тестирование

### Запуск тестов

**Используя Make (рекомендуется):**
```bash
cd backend

# Показать все доступные команды
make help

# Запустить все тесты
make test

# Запустить тесты с проверкой гонок (race conditions)
make test-race

# Запустить бенчмарки
make benchmark

# Генерация отчета о покрытии кода (HTML)
make coverage

# Показать покрытие в консоли
make coverage-text
```

**Используя go test напрямую:**
```bash
cd backend

# Запустить все тесты
go test ./...

# Запустить тесты с подробным выводом
go test -v ./...

# Запустить только unit тесты (без бенчмарков)
go test ./internal/handlers -run "^Test"

# Запустить конкретный тест
go test ./internal/handlers -run "TestSubmissionHandler_CreateSubmission_Success"

# Запустить только бенчмарки
go test ./internal/handlers -bench=.

# Запустить бенчмарки с информацией о памяти
go test ./internal/handlers -bench=. -benchmem

# Тесты с покрытием кода
go test -cover ./...

# Создать HTML отчет о покрытии
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html
```

### Структура тестов

```
backend/internal/handlers/
├── submission_test.go           # Unit тесты для API endpoints
└── benchmark_test.go           # Бенчмарки производительности
```

## 🛠️ Поддерживаемые языки

- Python (.py)
- C++ (.cpp)
- Java (.java)
- JavaScript (.js)
- Kotlin (.kt)

