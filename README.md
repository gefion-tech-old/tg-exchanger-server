# TG Exchanger Server

## Table of Contents
- [REST API](#rest-api)
    - [Public](#public)
        - [Registration in Bot](#registration-in-bot)
        - [Auth](#auth)
    - [Private](#private)
- [Database](#database)
    - [Postgres](#postgres)
    - [Redis](#redis)
- [Tools](#tools)
    - [Migrations](#migrations)


## REST API

### Public

#### Registration in Bot

- [POST] `/api/v1/bot/registration` — Регистрация пользователя через интерфейс Telegram

**Request**

```json
{
    "chat_id": 564353,
    "username": "SomeUsername"
}
```

#### Registration in Admin Panel


Для регистрации в админке необходимо быть добавленым в список сотрудников и иметь начатый диалог с ботом

##### Step 1

Код подтверждения актуален только в течении **30 минут**.

- [POST] `/api/v1/admin/registration/code` — В ЛС пользователю отправиться код подтверждения. 

***Request***

```json
{
    "password": "4tfgefhey75uh",
    "username": "I0HuKc"
}
```

***Response***

```json
{}
```

##### Step 2

- [POST] `/api/v1/admin/registration` — Завершение регистрации пользователя

***Request***

```json
{
    "code": 588227
}
```

#### Auth in Admin Panel


- [POST] `/api/v1/admin/auth` — Войти в созданный аккаунт

***Request***

```json
{
    "password": "4tfgefhey75uh",
    "username": "I0HuKc"
}
```

***Response***

```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjRmYzdhNmY5LWUzNzgtNDI0ZS05N2FlLTgxODZiNDI0N2FiOSIsImF1dGhvcml6ZWQiOnRydWUsImNoYXRfaWQiOjM1NDYyMjMsImV4cCI6MTYzODQ1ODAxMywidXNlcm5hbWUiOiJJMEh1S2MifQ.J90F-4a__q3uMkRWAS0K-IxXczT7t1rnPZqc1GAeDWU",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaGF0X2lkIjozNTQ2MjIzLCJleHAiOjE2MzkwNjE5MTMsInJlZnJlc2hfdXVpZCI6IjgxYmY4ZGI1LTY0MGItNDQ0Zi1iMDM0LWYwMWJjNjUwN2RiOCIsInVzZXJuYW1lIjoiSTBIdUtjIn0.YaxLs25XUbgSTDqSxwNoqdrQ9CNl40PoTznVVRe81z4"
}
```

## Migrations

**Накатить/Откатить миграцию**

```bash
migrate -path migrations -database "postgres://exchanger:qwerty@localhost:5432/exchanger_server_dev?sslmode=disable" up/down
```