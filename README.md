# expense-tracker-api

REST API для учёта личных расходов. Учебный проект на чистом `net/http`, без фреймворков.

## Стек

- Go 1.25
- `net/http` — HTTP-сервер и роутинг (метод в шаблоне маршрута, Go 1.22+)
- PostgreSQL 17 в Docker, драйвер `pgx` с пулом соединений
- миграции через golang-migrate
- `shopspring/decimal` — денежные суммы без потери точности

## Структура

```
cmd/api/              точка входа: подключение к базе, сборка зависимостей, запуск
internal/
  model/              структуры данных
  repository/         хранение: PostgreSQL и вариант в памяти
  service/            бизнес-правила и валидация
  handler/            HTTP: разбор запроса, коды ответов, JSON
migrations/           миграции схемы базы
docker-compose.yml    PostgreSQL для локальной разработки
```

Зависимости направлены в одну сторону: `handler → service → repository`. Сервис работает с хранилищем через интерфейс, объявленный в самом сервисе, поэтому не зависит от конкретной реализации — её подставляет `main`. Переход с хранения в памяти на PostgreSQL поменял одну строку в `main` и не затронул сервис.

Через все слои передаётся `context.Context` из HTTP-запроса: если клиент отключился, запрос к базе тоже отменяется. Ошибки оборачиваются с указанием слоя, клиенту при внутренней ошибке уходит только `500`, подробности пишутся в лог.

## Запуск

1. Поднять базу:

```
docker compose up -d
```

2. Задать адрес базы (PowerShell):

```
$env:DB_URL = "postgres://expenses:expenses@localhost:5432/expenses?sslmode=disable"
```

3. Применить миграции:

```
migrate -path migrations -database $env:DB_URL up
```

4. Запустить сервер:

```
go run ./cmd/api
```

Сервер поднимется на http://localhost:8080. Адрес базы читается из переменной окружения `DB_URL`, в коде пароля нет.

## База данных

| Команда | Что делает |
|---|---|
| `docker compose ps` | проверить, работает ли |
| `docker compose stop` | остановить, данные сохраняются |
| `docker compose start` | запустить снова |
| `docker compose logs db` | логи базы |
| `docker compose down` | удалить контейнер, данные сохраняются в volume |
| `docker compose down -v` | удалить контейнер вместе с данными |

Подключиться к базе вручную:

```
docker compose exec db psql -U expenses -d expenses
```

Логин, пароль и имя базы заданы в `docker-compose.yml` — значения учебные, только для локального запуска.

## Миграции

Схема базы описана в папке `migrations/` и применяется через [golang-migrate](https://github.com/golang-migrate/migrate). Руками структуру базы не меняем — только новыми миграциями.

Установка:

```
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

| Команда | Что делает |
|---|---|
| `migrate -path migrations -database $env:DB_URL up` | применить все миграции |
| `migrate -path migrations -database $env:DB_URL down 1` | откатить последнюю |
| `migrate create -ext sql -dir migrations -seq название` | создать новую |

### Таблица `expenses`

| Колонка | Тип | Ограничения |
|---|---|---|
| `id` | `BIGSERIAL` | первичный ключ |
| `amount` | `NUMERIC(14, 2)` | обязательна, больше нуля |
| `category` | `TEXT` | обязательна |
| `note` | `TEXT` | по умолчанию пустая строка |
| `created_at` | `TIMESTAMPTZ` | по умолчанию текущее время |

Индекс по `created_at` — для выборок за период.

## Тесты

```
go test ./...
go test ./... -cover
```

Бизнес-правила покрыты table-driven тестами с подставным хранилищем — база для них не нужна.

## Ручки

| Метод | Путь | Что делает | Коды |
|---|---|---|---|
| GET | `/healthz` | статус сервиса | 200 |
| POST | `/expenses` | создать расход | 201, 400, 500 |
| GET | `/expenses` | список расходов, новые сверху | 200, 500 |

### Правила валидации

- `category` — обязательна, пробелы по краям обрезаются
- `amount` — строго больше нуля; то же правило дублируется ограничением `CHECK` в базе
- `id` и `created_at` назначает база, из запроса не принимаются

## Проверка вручную

Создай файл `body.json`:

```json
{
  "amount": 1500.50,
  "category": "food",
  "note": "обед"
}
```

| Что проверяем | Команда |
|---|---|
| здоровье сервиса | `curl.exe http://localhost:8080/healthz` |
| создать расход | `curl.exe -X POST http://localhost:8080/expenses -H "Content-Type: application/json" -d "@body.json"` |
| список | `curl.exe http://localhost:8080/expenses` |
| код и заголовки | `curl.exe -i http://localhost:8080/expenses` |

Проверка валидации — ждём 400:

```
curl.exe -X POST http://localhost:8080/expenses -H "Content-Type: application/json" -d "{\"amount\": -5, \"category\": \"food\"}"
```

> Команды написаны для PowerShell, где слово `curl` занято встроенной командой — поэтому `curl.exe`. В Linux и macOS достаточно `curl`.

## Планы

- получение, изменение и удаление расхода по id
- фильтры по датам и категории, отчёт по категориям, постраничная выдача
- структурированные логи с trace id, middleware, graceful shutdown
- Dockerfile для самого сервиса
- GitHub Actions: lint, vet, test