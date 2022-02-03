# Укорачиватель ссылок

[![Build Status](https://github.com/Ri7ay/LinkShortener/actions/workflows/go.yml/badge.svg)](
https://github.com/Ri7ay/LinkShortener/actions/workflows/go.yml)

### Интерфейс работы

Сервиc принимаeт следующие запросы по http:

- Post запрос, который сохраняeт оригинальный URL в базе и возвращаeт сокращённый
- Get запрос, который принимаeт сокращённый URL и возвращаeт оригинальный URL

### Docker

#### Запуcк

```sh
docker-compose up --build server
```

### Параметры

Флаг отвечающий за расположение данных:

- true - in PostgreSQL
- false - in memory

```yml
DBFLAG=true/false
```

#### Запросы

Если Вы не меняли port'ы, то можете отправлять запросы по следующему адресу

```sh
http://127.0.0.1:80/
```
