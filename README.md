# expense-tracker-api

REST API для учёта личных расходов. Учебный проект на чистом `net/http`, без фреймворков.

## Стек

- Go 1.25
- `net/http` — HTTP-сервер и роутинг (метод в шаблоне маршрута, Go 1.22+)
- `shopspring/decimal` — денежные суммы без потери точности
- PostgreSQL 17 в Docker (подключение кода к базе в работе; сейчас данные хранятся в памяти)

## Структура

```
cmd/api/              точка входа: сборка зависимостей и запуск
internal/
  model/              структуры данных
  repository/         хранение (сейчас в памяти)
  service/            бизнес-правила и валидация
  handler/            HTTP: разбор запроса, коды ответов, JSON
docker-compose.yml    PostgreSQL для локальной разработки
```

Зависимости направлены в одну сторону: `handler → service → repository`. Сервис работает с хранилищем через интерфейс, объявленный в самом сервисе, поэтому не зависит от конкретной реализации — её подставляет `main`.

## База данных

PostgreSQL поднимается в Docker одной командой:

```
docker compose up -d
```

Проверить, что работает:

```
docker compose ps
```

Подключиться к базе вручную:

```
docker compose exec db psql -U expenses -d expenses
```

| Команда | Что делает |
|---|---|
| `docker compose stop` | остановить, данные сохраняются |
| `docker compose start` | запустить снова |
| `docker compose logs db` | логи базы |
| `docker compose down` | удалить контейнер, данные сохраняются в volume |
| `docker compose down -v` | удалить контейнер вместе с данными |

Логин, пароль и имя базы заданы в `docker-compose.yml` — значения учебные, только для локального запуска.

## Запуск

```
go run ./cmd/api
```

Сервер поднимется на http://localhost:8080

## Тесты

```
go test ./...
go test ./... -cover
```

Бизнес-правила покрыты table-driven тестами с подставным хранилищем.

## Ручки

| Метод | Путь | Что делает | Коды |
|---|---|---|---|
| GET | `/healthz` | статус сервиса | 200 |
| POST | `/expenses` | создать расход | 201, 400 |
| GET | `/expenses` | список расходов | 200 |

### Правила валидации

- `category` — обязательна, пробелы по краям обрезаются
- `amount` — строго больше нуля
- `id` и `created_at` назначает сервер, из запроса не принимаются

## Проверка вручную

Создай файл `body.json`:

```json
{
  "amount": 1500.50,
  "category": "food",
  "note": "обед"
}
```

Здоровье сервиса:

```
curl.exe http://localhost:8080/healthz
```

Создать расход:

```
curl.exe -X POST http://localhost:8080/expenses -H "Content-Type: application/json" -d "@body.json"
```

Список расходов:

```
curl.exe http://localhost:8080/expenses
```

Проверка валидации — ждём 400:

```
curl.exe -X POST http://localhost:8080/expenses -H "Content-Type: application/json" -d "{\"amount\": -5, \"category\": \"food\"}"
```

Флаг `-i` показывает код ответа и заголовки:

```
curl.exe -i http://localhost:8080/expenses
```

> Команды написаны для PowerShell, где слово `curl` занято встроенной командой — поэтому `curl.exe`. В Linux и macOS достаточно `curl`.

## Планы

- перевести репозиторий на PostgreSQL, миграции через golang-migrate
- фильтры по датам и категории, отчёт по категориям
- Dockerfile для самого сервиса
- структурированные логи с trace id, метрики
- GitHub Actions: lint, vet, test