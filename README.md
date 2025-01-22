# Сервис сохранения метрик
 
## Основные команды
- `make start` - запуск проекта
- `make stop` - остановка проекта
- `make down` - остановка и удаление всех артефактов проекта
- `make sh` - провалиться в терминал контейнера
---

## Запуск приложения
1. `make start` - запустить сборку образов и запуск контейнеров
2. `make sh` - провалиться в go-контейнер.
3. `go run -tags='no_mysql no_sqlite3 no_ydb no_clickhouse no_libsql no_vertica no_mssql' ./cmd/migration_tool` - выполнить миграции
4. `curl -X POST "http://localhost:8889/v0.1/update" \
     -H "Content-Type: application/json" \
     -H "User-Agent: curl/8.5.0" \
     -H "Accept: */*" \
     -d '{
           "id": "PollCount",
           "type": "counter",
           "delta": 5
         }'` - записать метрику в бд
5. `curl -X GET http://localhost:8889/v0.1/value/counter/PollCount -H "User-Agent: curl/8.5.0" -H "Accept: */*"` - получить метрику методом GET
6. `curl -X POST "http://localhost:8889/v0.1/value" \
     -H "Content-Type: application/json" \
     -H "User-Agent: curl/8.5.0" \
     -H "Accept: */*" \
     -d '{
           "id": "PollCount",
           "type": "counter"
         }'` - получить метрику методом POST

## Дополнительно
1. TODO 
