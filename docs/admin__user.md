Для регистрации в админке необходимо быть добавленым в список сотрудников и иметь начатый диалог с ботом

### Registration Step One

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

### Registration Step Two

- **[POST]** `/api/v1/admin/registration` — Завершение регистрации пользователя

***Request***

```json
{
    "code": 588227
}
```

### Auth


- **[POST]** `/api/v1/admin/auth` — Войти в созданный аккаунт

***Request***

```json
{
    "password": "4tfgefhey75uh",
    "username": "I0HuKc"
}
```

***Response***

Время жизни `access_token` — **15 минут**

Время жизни `refresh_token` — **7 дней**

```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjRmYzdhNmY5LWUzNzgtNDI0ZS05N2FlLTgxODZiNDI0N2FiOSIsImF1dGhvcml6ZWQiOnRydWUsImNoYXRfaWQiOjM1NDYyMjMsImV4cCI6MTYzODQ1ODAxMywidXNlcm5hbWUiOiJJMEh1S2MifQ.J90F-4a__q3uMkRWAS0K-IxXczT7t1rnPZqc1GAeDWU",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaGF0X2lkIjozNTQ2MjIzLCJleHAiOjE2MzkwNjE5MTMsInJlZnJlc2hfdXVpZCI6IjgxYmY4ZGI1LTY0MGItNDQ0Zi1iMDM0LWYwMWJjNjUwN2RiOCIsInVzZXJuYW1lIjoiSTBIdUtjIn0.YaxLs25XUbgSTDqSxwNoqdrQ9CNl40PoTznVVRe81z4"
}
```

### Update Access Token

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

### Logout

При выходе из хранилища пользовательских сессий удаляется `Access Token`, т.е еще раз авторизовываться с ним нельзя.

- **[POST]** `/api/v1/admin/logout` — Выйти из аккаунта в панели администратора.

***Header***

```json
{
    "Authorization": "Bearer <token>"
}
```

***Response***

```json
{
    "message": "successfully logged out"
}
```