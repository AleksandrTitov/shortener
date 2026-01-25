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
curl -i -b 'id_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjkzNzgyMjAsIlVzZXJJRCI6IjdmYTlhYWZlLTFiM2QtNGY3Ni04YjRmLWIxMmM3ODE3Yjk3OCJ9.C65SJdHisSezNvjm-CboRMmQ5wpAjRm1sXTv18FgkME' \
-X POST -H "Content-Type: text/plain" -d "http://tst1.ru/" http://localhost:8081/
```