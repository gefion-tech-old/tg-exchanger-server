## Bot Messages

Сообщения которые использует бот для общения в TG с пользователями.

### Create

- **[POST]** `/api/v1/admin/message` — Создать новое сообщение для бота.

***Header***

```json
{
    "Authorization": "Bearer <token>"
}
```

***Request***

Поле `created_by` следует брать из `access_token`

```json
{
    "connector": "msg_connector",
    "message_text": "some message text here",
    "created_by": "I0HuKc"
}
```

***Response***

```json
{
    "id": 3,
    "connector": "text_connector",
    "message_text": "some message text here",
    "created_by": "I0HuKc",
    "created_at": "2021-12-06T12:49:35.303698Z",
    "updated_at": "2021-12-06T12:49:35.303698Z"
}
```

<hr>

### Update

- **[PUT]** `/api/v1/admin/message/:id` — Обновить сообщение для бота.

```json
{
    "Authorization": "Bearer <token>"
}
```

***Request***

```json
{   
    "message_text": "some message text here",   
}
```

***Response***

```json
{
    "id": 3,
    "connector": "text_connector",
    "message_text": "some message text here",
    "created_by": "I0HuKc",
    "created_at": "2021-12-06T12:49:35.303698Z",
    "updated_at": "2021-12-06T12:49:35.303698Z"
}
```

<hr>

### Get

- **[GET]** `/api/v1/admin/message/:connector` — получить конкретное сообщение из БД

***Response***

```json
{
    "id": 3,
    "connector": "text_connector",
    "message_text": "some message text here",
    "created_by": "I0HuKc",
    "created_at": "2021-12-06T12:49:35.303698Z",
    "updated_at": "2021-12-06T12:49:35.303698Z"
}
```

<hr>

### Get Slice of Messages

Принцип работы с данным запросом и его параметрами идентичен [уведомлениям](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__notification.md#get-slice-of-notifications).

- **[GET]** `/api/v1/admin/messages?page=1&limit=15` — получить массив сообщений из БД

***Response***

```json
{
    "current_page": 1,
    "data": [],
    "last_page": 2,
    "limit": 15
}
```

<hr>

### Delete

- **[DELETE]** `/api/v1/admin/message/:id` — Удалить сообщение бота.

***Header***

```json
{
    "Authorization": "Bearer <token>"
}
```


***Response***

```json
{
    "id": 3,
    "connector": "text_connector",
    "message_text": "some message text here",
    "created_by": "I0HuKc",
    "created_at": "2021-12-06T12:49:35.303698Z",
    "updated_at": "2021-12-06T12:49:35.303698Z"
}
```