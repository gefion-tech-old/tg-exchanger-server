# TG Exchanger Server

## Table of Contents
- [Base Info](#base-info)
- [REST API](https://github.com/exchanger-bot/docs) 
- [Database](#database)
    - [Postgres](#postgres)
    - [Redis](#redis)
- [Tools](#tools)
    - [Migrations](#migrations)
        - [Create Migration](#create-migration)
        - [Up/Down Migration](#up/down-migration)


## Base Info


### Start Server

1. Необходимо поднять NSQ очередь использую docker-compose

```
sudo docker-compose up -d
```

2. Запустить сервер 

```
./server -prod true
```


## Tools

### Migrations

#### Create Migration

```
migrate create -ext sql -dir migrations migration_name
```

#### Up/Down Migration

```
migrate -path migrations -database "postgres://exchanger:qwerty@localhost:5432/exchanger_server_dev?sslmode=disable" up/down
```
