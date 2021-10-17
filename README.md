
### Приложение для вывода иерархии сотрудников
```bash
Swagger route /swagger/index.html
Поиск сотрудника осуществляется по id
Route /api/:id (например, http://localhost:9000/api/6167f1ce3e6afbabae66b4a5)
```


### Запуск mongo в docker

```bash
./docker run --name=mongo_db -e MONGODB_DATABASE=admin -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=123 -p 27019:27017 -d mongo
Скрипт с данными БД /scripts/db employees data.json 
Данные вставляются при запуске приложения. Чтобы отменить загрузку данных при старте, нужно убрать вызов функции loadDbData в файле pkg/repository/mongo.go.
```
