# TG Exchanger Server

## Table of Contents
- [REST API](#rest-api)    
        - [Bot API Routs](#bot-api-routs)
            - [Registration in Bot](#registration-in-bot)
            - [Create User Bill](#create-user-bill)
            - [Delete User Bill](#delete-user-bill)
            - [Get All User Bills](#get-all-user-bills)
        - [Admin API Routs](#admin-api-routs)
            - [Public](#public)
                - [Registration in Admin Panel](#registration-in-admin-panel)
                - [Auth in Admin Panel](#auth-in-admin-panel)       
                - [Update Access Token](#update-access-token)
            - [Private](#private)
                - [Logout from Admin Panel](#logout-from-admin-panel)  
- [Database](#database)
    - [Postgres](#postgres)
    - [Redis](#redis)
- [Tools](#tools)
    - [Migrations](#migrations)
        - [Create Migration](#create-migration)
        - [Up/Down Migration](#up/down-migration)


## REST API


### Bot API Routs

#### Registration in Bot

- **[POST]** `/api/v1/bot/registration` — Регистрация пользователя через интерфейс Telegram

***Request***

```json
{
    "chat_id": 564353,
    "username": "I0HuKc"
}
```

#### Create User Bill 

- **[POST]** `/api/v1/bot/user/bill` **[201]** — Создать новый пользовательский счет

***Request***

```json
{
    "chat_id": 524164407,
    "bill": "535949490410854"
}
```

***Response***

```json
{
    "bill": "535949490410854",
    "chat_id": 524164407,
    "created_at": "2021-12-04T14:10:04.12226Z",
    "id": 19
}
```

#### Delete User Bill

- **[DELETE]** `/api/v1/bot/user/bill` **[200]** — Удалить новый пользовательский счет

***Request***

```json
{
    "chat_id": 524164407,
    "bill": "535949490410854"
}
```

***Response***

```json
{}
```

#### Get All User Bills

- **[GET]** `/api/v1/bot/user/<chat_id>/bills` **[200]** — Получить список всех пользовательских счетов

***Response***

```json
{
    "bills": [
        {
            "id": 19,
            "chat_id": 524164407,
            "bill": "535949490410854",
            "created_at": "2021-12-04T14:10:04.12226Z"
        }
    ]
}
```

### Admin API Routs

#### Public

#### Registration in Admin Panel


Для регистрации в админке необходимо быть добавленым в список сотрудников и иметь начатый диалог с ботом

##### Step 1

Код подтверждения актуален только в течении **30 минут**.

- **[POST]** `/api/v1/admin/registration/code` — В ЛС пользователю отправиться код подтверждения. 

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

- **[POST]** `/api/v1/admin/registration` — Завершение регистрации пользователя

***Request***

```json
{
    "code": 588227
}
```

#### Auth in Admin Panel


- **[POST]** `/api/v1/admin/auth` — Войти в созданный аккаунт

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

#### Update Access Token

- **[POST]** `/api/v1/admin/auth` — Обновить сессию

***Request***

```json
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaGF0X2lkIjo1MjQxNjQ0MDcsImV4cCI6MTYzOTEzMDMwNiwicmVmcmVzaF91dWlkIjoiZWM0NDExNmYtZmRkZS00ZWE2LWE2OTItYTVhNmI0ZTBmMjliIiwidXNlcm5hbWUiOiJJMEh1S2MifQ.Bm_E6IIW4k4FK7GDecwJOvxFVw248bkUTl8SZ5uJGVQ"
}
```

***Response***

```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjRmYzdhNmY5LWUzNzgtNDI0ZS05N2FlLTgxODZiNDI0N2FiOSIsImF1dGhvcml6ZWQiOnRydWUsImNoYXRfaWQiOjM1NDYyMjMsImV4cCI6MTYzODQ1ODAxMywidXNlcm5hbWUiOiJJMEh1S2MifQ.J90F-4a__q3uMkRWAS0K-IxXczT7t1rnPZqc1GAeDWU",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaGF0X2lkIjozNTQ2MjIzLCJleHAiOjE2MzkwNjE5MTMsInJlZnJlc2hfdXVpZCI6IjgxYmY4ZGI1LTY0MGItNDQ0Zi1iMDM0LWYwMWJjNjUwN2RiOCIsInVzZXJuYW1lIjoiSTBIdUtjIn0.YaxLs25XUbgSTDqSxwNoqdrQ9CNl40PoTznVVRe81z4"
}
```

#### Private

#### Logout from Admin Panel 

При выходе из хранилища пользовательских сессий удаляется `Access Token`, т.е еще раз авторизовываться с ним нельзя.

- **[POST]** `/api/v1/admin/logout` — Выйти из аккаунта в панели администратора.

***Header***

```json
{
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImZiYjIwNDA0LTFlNWQtNGFiOS04MDIzLWY3NjQ2ZTEzYWU3MSIsImF1dGhvcml6ZWQiOnRydWUsImNoYXRfaWQiOjc2NTYyNTYsImV4cCI6MTYzODAyMDA0Mywicm9sZSI6MX0.Lo14QcF_MmGwUFHiSo452aTtsRzV4sZ0a-HdGyWzyds"
}
```

***Response***

```json
{
    "message": "successfully logged out"
}
```

## Migrations

### Create Migration

```
migrate create -ext sql -dir migrations migration_name
```

### Up/Down Migration

```
migrate -path migrations -database "postgres://exchanger:qwerty@localhost:5432/exchanger_server_dev?sslmode=disable" up/down
```
