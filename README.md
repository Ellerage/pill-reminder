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

```
# Local build
make build_prod

# Stop
sudo systemctl stop pill-reminder

# Upload
scp ./pill-reminder <ssh>

# Update
sudo systemctl daemon-reexec
sudo systemctl daemon-reload

# Start
sudo systemctl enable pill-reminder
sudo systemctl start pill-reminder

# Status Check
sudo systemctl status pill-reminder

# Logs
journalctl -u pill-reminder -n 50
```
