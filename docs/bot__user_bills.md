## Bills

### Create

- **[POST]** `/api/v1/admin/bill` — Создать новый пользовательский счет

***Header***

```json
{
    "Authorization": "Bearer <token>"
}
```

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

<hr>

### Reject

После выполнения запроса на указанный `chat_id` бот отправит уведолмение.

- **[POST]** `/api/v1/admin/bill/reject` — отклонить верификацию карты

***Header***

```json
{
    "Authorization": "Bearer <token>"
}
```

***Request***

```json
{
    "chat_id": 524164407,
    "bill": "535949490410854",
    "reason": "какая-то причина тут."
}
```


### Delete

- **[DELETE]** `/api/v1/bot/user/bill` — Удалить новый пользовательский счет

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

<hr>

### Get Bill 

- **[GET]** `/api/v1/bot/user/bill/:id` — Получить пользовательский счет

***Response***

```json
{
    "chat_id": 524164407,
    "bill": "535949490410854"
}
```

<hr>

### Get All

- **[GET]** `/api/v1/bot/user/<chat_id>/bills` — Получить список всех пользовательских счетов

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