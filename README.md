# TG Exchanger Server

## Table of Contents
- [REST API](#rest-api)    
    - [Bot API Routs](#)
        - [User](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user.md)
            - [Registration](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user.md#registration)
        - [Bills](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user_bills.md)
            - [Create](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user_bills.md#create)
            - [Delete](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user_bills.md#delete)
            - [Get All](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user_bills.md#get-all)
    - [Admin API Routs](#)
        - [User](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__user.md)
            - [Registration Step One](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__user.md#registration-step-one)
            - [Registration Step Two](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__user.md#registration-step-two)     
            - [Auth](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__user.md#auth)
            - [Logout](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__user.md#logout)
            - [Update Access Token](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__user.md#update-access-token)        
        - [Bot Message](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__bot_messages.md)
            - [Create](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__bot_messages.md#create)
            - [Update](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__bot_messages.md#update)
            - [Get](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__bot_messages.md#get)
            - [Get All](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__bot_messages.md#get-all)
            - [Delete](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__bot_messages.md#delete)            
- [Database](#database)
    - [Postgres](#postgres)
    - [Redis](#redis)
- [Tools](#tools)
    - [Migrations](#migrations)
        - [Create Migration](#create-migration)
        - [Up/Down Migration](#up/down-migration)

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
