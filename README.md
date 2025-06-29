# pill-reminder

A pill reminder project with integration to a Telegram bot and recursive notifications until the user confirms.

## System Design

https://miro.com/app/board/uXjVIjNVeLA=/?share_link_id=578147349731

## Before starting

Need to create a `.env` file at the root of directory

```
BOT_TOKEN=<YOUR_BOT_TOKEN>
MONGO_URL=<URL_FOR_CONNECTION_TO_MONGO>
MONGO_DB_NAME=<DATABASE_NAME_IN_MONGO>
TIMEZONE=<BASE_TIMEZONE_FOR_MANAGING_NOTIFY_TIME>
```

## How to start

```
go mod tidy && make
```
