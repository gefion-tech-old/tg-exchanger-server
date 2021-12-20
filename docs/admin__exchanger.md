## Exchanger

### Create

- **[POST]** `/api/v1/admin/exchanger` — Создать запись нового обменника

***Header***

```json
{
    "Authorization": "Bearer <token>"
}
```

***Request***

```json
{
    "created_by": "I0HuKc",
    "name": "1obmen",
    "url": "https://1obmen.net/request-exportxml.xml"
}
```

***Response***

```json
{
    "id": 2,
    "name": "test",
    "url": "https://1obmen.net/request-exportxml.xml",
    "created_by": "I0HuKc",
    "created_at": "2021-12-17T18:31:25.690833Z",
    "updated_at": "2021-12-17T18:31:25.690833Z"
}
```


<hr>

### Update

- **[PUT]** `/api/v1/admin/exchanger/:id` — Обновить запись обменника

***Header***

```json
{
    "Authorization": "Bearer <token>"
}
```

***Request***

```json
{   
    "name": "1obmen",
    "url": "https://1obmen.net/request-exportxml.xml"
}
```

***Response***

```json
{
    "id": 2,
    "name": "test",
    "url": "https://1obmen.net/request-exportxml.xml",
    "created_by": "I0HuKc",
    "created_at": "2021-12-17T18:31:25.690833Z",
    "updated_at": "2021-12-17T18:31:25.690833Z"
}
```

<hr>

### Get

- **[GET]** `/api/v1/admin/exchanger/:id` — Получить одну запись конкретного обменника

***Response***

```json
{
    "id": 2,
    "name": "test",
    "url": "https://1obmen.net/request-exportxml.xml",
    "created_by": "I0HuKc",
    "created_at": "2021-12-17T18:31:25.690833Z",
    "updated_at": "2021-12-17T18:31:25.690833Z"
}
```

<hr>

### Get Slice of Exchangers

- **[GET]** `/api/v1/admin/exchanger?page=1&limit=15` — Получить нужную выборку записей

Принцип работы `page` идентичен `limit` остальным страницам.

<hr>

### Delete

***Header***

- **[DELETE]** `/api/v1/admin/exchanger/:id` — Обновить запись обменника

```json
{
    "Authorization": "Bearer <token>"
}
```

***Response***

```json
{
    "id": 2,
    "name": "test",
    "url": "https://1obmen.net/request-exportxml.xml",
    "created_by": "I0HuKc",
    "created_at": "2021-12-17T18:31:25.690833Z",
    "updated_at": "2021-12-17T18:31:25.690833Z"
}
```

<hr>

### Get File 

- **[GET]** `/api/v1/admin/exchanger/document` — Получить документ с данными обменников

```json
{
    "Authorization": "Bearer <token>"
}
```

***Response***

```json
{
    "file": "/var/www/html/tmp/file/2021-12-17T12:58:04.05367541.xlsx"
}
```