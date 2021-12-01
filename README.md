# TG Exchanger Server

## Table of Contents
- [REST API](#rest-api)
    - [Public](#public)
        - [Registration](#registration)
        - [Auth](#auth)
    - [Private](#private)
- [Database](#database)
    - [Postgres](#postgres)
    - [Redis](#redis)
- [Tools](#tools)
    - [Migrations](#migrations)


## Migrations

**Накатить/Откатить миграцию**

```
migrate -path migrations -database "postgres://exchanger:qwerty@localhost:5432/exchanger_server_dev?sslmode=disable" up/down
```