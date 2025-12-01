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
curl -X POST -H "Content-Type: text/plain" -d "https://tst.ru/" http://localhost:8080/
```