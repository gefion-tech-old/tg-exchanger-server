# TG Exchanger Server

## Table of Contents
- [Base Info](#base-info)
- [REST API](#rest-api)    
    - [Bot API Routs](#)
        - [User](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user.md)
            - [Registration](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user.md#registration)
        - [Bills](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user_bills.md)
            - [Create](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user_bills.md#create)
            - [Reject](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user_bills.md#reject)
            - [Delete](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user_bills.md#delete)
            - [Get Bill](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/bot__user_bills.md#get-bill)
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
            - [Get Slice of Messages](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__bot_messages.md#get-slice-of-messages)
            - [Delete](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__bot_messages.md#delete)       
        - [Notification](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__notification.md)   
            - [Notification Type](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__notification.md#notification-type) 
            - [Notification Status](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__notification.md#notification-status)
            - [Create](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__notification.md#create) 
                - [854](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__notification.md#854)
                - [855](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__notification.md#855) 
            - [Get Slice of Notifications](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__notification.md#get-slice-of-notifications)
            - [Update Status](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__notification.md#update-status) 
            - [Delete](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__notification.md#delete) 
        - [Exchanger](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__exchanger.md)
            - [Create](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__exchanger.md#create)
            - [Update](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__exchanger.md#update)
            - [Get](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__exchanger.md#get)
            - [Get Slice of Exchangers](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__exchanger.md#get-slice-of-messages)
            - [Delete](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__exchanger.md#delete)
            - [Get File](https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__exchanger.md#get-file)
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
