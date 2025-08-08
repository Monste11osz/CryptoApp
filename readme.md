# 📊 Crypto Price Service

Сервис для получения и управления списком криптовалют с их актуальной ценой.

---

## 🚀 Старт

1. Клонируйте репозиторий:
```bash
git clone https://github.com/Monste11osz/CryptoApp.git
cd CryptoApp
```

💡 Если проект у вас в .zip просто распакуйте и перейдите в папку:

1. Распакуйте проект:
```bash
unzip testYTask.zip
cd testYTask
```

2. Запустите через Docker:
```bash
docker compose up --build
```

3. Примечание по подключению к БД

   Если в конфигурационном файле (config/conf.json) или в переменных окружения указан localhost или 127.0.0.1 в качестве хоста базы данных, замените его на db — имя сервиса PostgreSQL из docker-compose.yml.
```json
{
  "DB_HOST": "db",
  "DB_PORT": 5432,
  "DB_USER": "admin",
  "DB_PASSWORD": "123123",
  "DB_NAME": "diplomdb"
}

```

4. После запуска сервис будет доступен по адресу:
```
http://localhost:8080
```

---

## 📦 Запуск в деталях

### Требования
- **Docker** и **Docker Compose**
- **Go 1.23+** (если хотите запускать локально без контейнера)
- **PostgreSQL 14+**

### Запуск с Docker
```bash
docker compose up --build
```

### Запуск локально (без Docker)
1. Установите зависимости:
```bash
go mod download
```
2. Запустите сервер:
```bash
go run ./cmd/main.go
```

---

## 📡 API эндпоинты

| Метод  | Путь                | Описание                                 |
|--------|---------------------|------------------------------------------|
| POST   | `/currency/price`   | Получить цену криптовалюты               |
| POST   | `/currency/add`     | Добавить криптовалюту в отслеживание     |
| DELETE | `/currency/remove`  | Удалить криптовалюту из отслеживания     |

---

## 📥 Примеры запросов и ответов

### 🔍 POST `/currency/price`

**Запрос:**
```json
{
  "coin": "bitcoin",
  "timestamp": 1754603100
}
```

**Успешный ответ (200):**
```json
{
  "status": "OK",
  "data": {
    "coin": "bitcoin",
    "price": 117200,
    "currency": "USD",
    "timestamp": 1754603017
  }
}
```

**Цена не найдена (404):**
```json
{
  "status": "NotFound",
  "message": "Price not found"
}
```

**Неверный запрос (400):**
```json
{
  "status": "ERROR",
  "message": "Invalid request"
}
```

**Ошибка сервиса при получении цены монеты (500):**
```json

{
"status": "ERROR",
"message": "Error while receiving data"
}
```

---

### ➕ POST `/currency/add`

**Запрос:**
```json
{
  "name_coin": "bitcoin"
}
```

**Успешный ответ:**
```json
{
  "status": "OK",
  "message": "Coin 'bitcoin' added",
  "data": {}
}
```

**Монета не найдена:**
```json
{
  "status": "ERROR",
  "message": "Coin not found"
}
```

**Некорректные данные:**
```json
{
  "status": "ERROR",
  "message": "Incorrect input data"
}
```
**Ошибка сервиса при добавлении монеты (500):**
```json
{
  "status": "ERROR",
  "message": "Service error while adding"
}

```

---

### ❌ DELETE `/currency/remove`

**Запрос:**
```json
{
  "name_coin": "bitcoin"
}
```

**Успешный ответ:**
```json
{
  "status": "OK",
  "message": "Coin deleted",
  "data": {}
}
```

**Монета отсутствует в списке:**
```json
{
  "status": "ERROR",
  "message": "This coin is not on the list"
}
```

**Некорректные данные:**
```json
{
  "status": "ERROR",
  "message": "Incorrect input data"
}
```
**Ошибка сервиса при удалении монеты (500):**
```json
{
  "status": "ERROR",
  "message": "Service error while deleting"
}
```


