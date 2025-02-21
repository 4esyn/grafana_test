# Geo Service API

## Описание

Сервис для поиска и геокодирования адресов с JWT аутентификацией, мониторингом Prometheus и визуализацией в Grafana.

## Требования

- Docker
- Docker Compose

## Быстрый старт

```bash
# Клонирование репозитория
git clone <repository-url>
cd grafana

# Запуск сервиса
docker-compose up --build

# Сервис будет доступен:
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- Метрики Prometheus: http://localhost:8080/metrics
- Prometheus UI: http://localhost:9090
- Grafana UI: http://localhost:3000 (admin/admin)

# Остановка сервиса
docker-compose down
```

## Структура проекта

```markdown
grafana/
├── geoservice/ # Основной сервис
│ ├── docs/ # Swagger документация
│ ├── handlers/ # Обработчики HTTP запросов
│ ├── metrics/ # Prometheus метрики
│ ├── middleware/ # Middleware для JWT
│ ├── models/ # Модели данных
│ └── storage/ # Хранилище пользователей
├── grafana/ # Конфигурация Grafana
│ ├── provisioning/ # Автоматическая настройка
│ │ ├── dashboards/ # Конфигурация дашбордов
│ │ └── datasources/ # Конфигурация источников данных
│ └── dashboards/ # JSON файлы дашбордов
└── prometheus.yml # Конфигурация Prometheus
```

## Доступные метрики Prometheus

### Метрики эндпоинтов

- geoservice_endpoint_duration_seconds_bucket - время выполнения запросов к эндпоинтам
- geoservice_endpoint_requests_total - количество запросов к эндпоинтам

### Метрики внешних сервисов

- geoservice_external_api_duration_seconds_bucket - время обращения к внешнему API
- geoservice_cache_duration_seconds_bucket - время обращения к кэшу
- geoservice_db_duration_seconds_bucket - время обращения к БД

## Grafana

### Доступ

- URL: http://localhost:3000
- Логин: admin
- Пароль: admin

### Предустановленный дашборд

Дашборд "Geo Service Dashboard" содержит следующие панели:

1. Время ответа эндпоинтов
2. Количество запросов в секунду
3. Время ответа внешнего API

## Тестирование API

### Регистрация

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### Получение JWT токена

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### Поиск адреса

```bash
curl -X POST http://localhost:8080/api/address/search \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{"query":"Москва"}'
```

### Геокодирование адреса

```bash
curl -X POST http://localhost:8080/api/address/geocode \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{"lat":"55.7558","lng":"37.6173"}'
```

## Создание и публикация образа Grafana

```bash
# Создание образа
docker commit grafana your-dockerhub-username/grafana:latest

# Публикация образа
docker login
docker push your-dockerhub-username/grafana:latest
```

## Prometheus запросы

```bash
# Среднее время ответа эндпоинтов за 5 минут
rate(geoservice_endpoint_duration_seconds_sum[5m]) / rate(geoservice_endpoint_duration_seconds_count[5m])

# Количество запросов в секунду
rate(geoservice_endpoint_requests_total[5m])

# 95-й процентиль времени ответа внешнего API
histogram_quantile(0.95, sum(rate(geoservice_external_api_duration_seconds_bucket[5m])) by (le))
```

## Разработка

### Генерация Swagger документации

```bash
cd geoservice
swag init
```

### Запуск в режиме разработки

```bash
cd geoservice
go run main.go
```

## Полезные команды

```bash
# Просмотр логов
docker-compose logs -f

# Перезапуск отдельного сервиса
docker-compose restart grafana

# Проверка статуса сервисов
docker-compose ps

# Остановка сервисов
docker-compose down
```
