# Order Processing Service

#### Микросервис для обработки онлайн заказов с использованием:
- Go (серверная часть)
- PostgreSQL (хранение данных)
- Kafka (обработка событий)
- Docker (развертывание)
- Chi (HTTP роутинг)

## Особенности
- Прием заказов через Kafka
- Кэширование заказов в памяти
- REST API для получения данных о заказах
- Веб-интерфейс для просмотра заказов

## Запуск проекта

### Требования:
- Docker
- Docker Compose

### Стандартный запуск:
```bash
# Клонировать репозиторий
git clone https://github.com/SANEKNAYMCHIK/order-service.git
cd order-service

# Создание .env файла по примеру
cp .env.example .env

# Запустить все сервисы, кроме генератора заказов
docker-compose up -d --build postgres zookeeper kafka order-service

# Запустить только генератор заказов
docker-compose up -d --build producer

# Запустить все сервисы
docker-compose up -d --build

# Остановить сервисы
docker-compose down
```
### Взаимодействие с сервисом:
```bash
# Открыть веб-интерфейс:
http://localhost:8080

# Получить данные о заказе по order_uid
http://localhost:8080/order/{order_uid}
```
### Структура API
GET /order/{uid} - Получить заказ по ID