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

```bash
# Запуск всего приложения
docker compose up -d
```

### Остановка

```bash
```bash
docker compose down

# С удалением volumes (очистка БД)
docker compose down -v
```

## 🛠️ Поддерживаемые языки

- Python (.py)
- C++ (.cpp)
- Java (.java)
- JavaScript (.js)
- Kotlin (.kt)

