# pill-reminder

A pill reminder project with integration to a Telegram bot and recursive notifications until the user confirms.

## System Design

https://miro.com/app/board/uXjVIjNVeLA=/?share_link_id=578147349731

## Before starting

Need to create a `.env` file at the root of the directory, actually, you need only BOT_TOKEN if you're going to usethe default docker configuration

```
BOT_TOKEN=<CODE>
MONGO_DB_NAME=bill-reminder
TIMEZONE=UTC
MONGO_URL=mongodb://localhost:27017/bill-reminder
REDIS_URL=localhost
REDIS_PORT=6379
REMINDER_DB=1
ASYNCQ_DB=0
MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=admin

```

## How to start

```
go mod tidy && make
```
