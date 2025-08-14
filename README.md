# pill-reminder

A pill reminder project with integration to a Telegram bot and recursive notifications until the user confirms.

## System Design

https://miro.com/app/board/uXjVIjNVeLA=/?share_link_id=578147349731

## Before starting

Need to create a `.env` file at the root of the directory you need only BOT_TOKEN if you're going to use the default docker configuration

```
BOT_TOKEN=<CODE>
```

## How to start

```
docker compose up -d

go mod tidy && make
```

## Deployment

run deploy-prod.sh
```
