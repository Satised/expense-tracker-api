# expense-tracker-api

Простой REST API для учёта расходов. Go + net/http, хранение в памяти.

## Запуск

Сервер: http://localhost:8080

## Проверка

Здоровье сервиса:

curl.exe http://localhost:8080/healthz


Создать расход:

curl.exe -X POST http://localhost:8080/expenses -H "Content-Type: application/json" -d "@body.json"


Список расходов:

curl.exe http://localhost:8080/expenses

Проверка валидации (ждём 400):

curl.exe -X POST http://localhost:8080/expenses -H "Content-Type: application/json" -d "@bad.json"

Флаг `-i` покажет код ответа и заголовки:

curl.exe -i http://localhost:8080/expenses

## Ручки

| Метод | Путь | Что делает | Код |
|---|---|---|---|
| GET | /healthz | статус сервиса | 200 |
| POST | /expenses | создать расход | 201 / 400 |
| GET | /expenses | список расходов | 200 |
