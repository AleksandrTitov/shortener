# sorter

Сервис для сокращения ссылок

### Запуск

```shell
  -a string
        Адрес сервера в формате <хост>:<порт> (default "localhost:8080")
  -b string
         HTTP адрес сервера в сокращенном URL в формате <http схема>://<хост>:<порт> (default "http://localhost:8080")
```

### Обращение
```shell
curl -i -X POST -H "Content-Type: text/plain" -d "https://tst1.ru/" http://localhost:8081/
```

```shell
curl -i -b 'id_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjkzODU4MjQsIlVzZXJJRCI6IjI2MjBhZGRmLTI2YWQtNDdmOC1iN2MxLWJjZGI5ZDliMzY5NyJ9.gldnIf5uvsuD8DJarZjR0MlgcGuBULlJIuOH4ussIHg' \
-X POST -H "Content-Type: text/plain" -d "https://tst2.ru/" http://localhost:8081/
```

```shell
curl -i -b 'id_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njk0Njc3MzcsIlVzZXJJRCI6Ijk3MDYzMGIwLWJjZDAtNDE2MS1hZTRmLTYwYmVhOGMwNzA0ZCJ9.uOqu6B2PuJU--yISmPIYb8I-0v6vhDriMsYnu-Dev2w' \
-X GET http://localhost:8081/api/user/urls
```

```shell
curl -X POST -H "Content-Type: application/json" \
  -d '[{"correlation_id": "111111","original_url": "https://practicum1.yandex.ru"},{"correlation_id": "111112","original_url": "https://practicum2.yandex.com"}]' \
  http://localhost:8081/api/shorten/batch
```