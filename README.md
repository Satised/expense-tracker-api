# expense-tracker-api

REST API для учёта личных расходов. Учебный проект, написан на чистом `net/http` без фреймворков.

## Стек

- Go 1.25
- `net/http` — HTTP-сервер и роутинг (метод в шаблоне маршрута, Go 1.22+)
- `shopspring/decimal` — денежные суммы без потери точности
- хранение в памяти (PostgreSQL в планах)

## Структура

```
cmd/api/              точка входа: сборка зависимостей и запуск
internal/
  model/              структуры данных
  repository/         хранение (сейчас в памяти)
  service/            бизнес-правила и валидация
  handler/            HTTP: разбор запроса, коды ответов, JSON
```

Зависимости направлены в одну сторону: `handler → service → repository`. Сервис работает с хранилищем через интерфейс, объявленный в самом сервисе, поэтому не зависит от конкретной реализации.

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

> Команды написаны для PowerShell, поэтому `curl.exe` — в PowerShell слово `curl` занято встроенной командой. В Linux и macOS достаточно `curl`.

## Планы

- PostgreSQL + миграции (golang-migrate)
- фильтры по датам и категории, отчёт по категориям
- Docker и docker-compose
- структурированные логи с trace id, метрики
- GitHub Actions: lint, vet, test