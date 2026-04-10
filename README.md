# SysPulse

SysPulse is a Go-based Discord bot that combines Discord gateway commands, HTTP interactions, host monitoring, and RSS delivery.

## Features

- Prefix commands over gateway: `!ping`, `!sys`, `!rss <feed>`
- Slash commands: `/ping`, `/sys`, `/rss` (registered at startup)
- Secure Discord interaction handling at `POST /interactions` using Ed25519 signature verification
- Health endpoint at `GET /health`
- System snapshot command (`sys`) with CPU, RAM, disk, and uptime
- Background CPU monitor with cooldown alerts (currently 80% threshold, 30s check interval, 2m cooldown)
- RSS support for `tech`, `news`, and `gaming`
- RSS polling cron job with in-memory dedupe and 5-minute feed cache for on-demand requests

## Tech Stack

- Go (`go.mod` currently declares `go 1.26.0`)
- `github.com/bwmarrin/discordgo`
- `github.com/mmcdole/gofeed`
- `github.com/robfig/cron/v3`
- `github.com/shirou/gopsutil/v3`

## Project Layout

- `main.go` - bootstrap, Discord session, HTTP server, worker startup/shutdown
- `bot/` - command registry, prefix handlers, slash command definitions, RSS send helpers
- `httpapi/` - interaction endpoint + signature verification, health endpoint
- `services/` - system stats, CPU worker, RSS fetch/cache/cron, worker manager
- `config/` - environment loading and validation

## Prerequisites

- A Discord application with bot token
- Discord application public key (for interaction signature verification)
- Bot invited to your server with required permissions
- Enable Message Content Intent if using prefix commands (`!`)
- Go toolchain compatible with `go.mod`

## Environment Variables

Set these before starting the app (via shell exports, `.env` + compose, or another secret manager):

```dotenv
PORT=8080
RSS_FEED_SIZE=10
RSS_CRON_SCHEDULE=0 5 * * *
PUBLIC_KEY=your_discord_public_key
APPLICATION_ID=your_application_id
GUILD_ID=your_guild_id
GUILD_CHANNEL_ID=your_alert_channel_id
CHANNEL_ID=your_default_channel_id
BOT_TOKEN=your_bot_token
```

Notes:

- `PORT`, `RSS_FEED_SIZE`, `PUBLIC_KEY`, `APPLICATION_ID`, `GUILD_ID`, `GUILD_CHANNEL_ID`, `CHANNEL_ID`, and `BOT_TOKEN` are required.
- `RSS_CRON_SCHEDULE` is optional; default is `0 5 * * *`.
- Cron timezone is fixed to `Asia/Kolkata` in current code.
- The app does not auto-load `.env` by itself in local runs; ensure env vars are exported in your shell when running `go run .` directly.

## Run Locally

```bash
go mod download
export PORT=8080
export RSS_FEED_SIZE=10
export PUBLIC_KEY=your_discord_public_key
export APPLICATION_ID=your_application_id
export GUILD_ID=your_guild_id
export GUILD_CHANNEL_ID=your_alert_channel_id
export CHANNEL_ID=your_default_channel_id
export BOT_TOKEN=your_bot_token
go run .
```

At startup, SysPulse:

1. Validates required env vars
2. Connects to Discord gateway
3. Registers slash commands for the configured guild
4. Starts HTTP server on `PORT` with `/interactions` and `/health`
5. Starts workers (CPU monitor + RSS cron)

## Build Binary

```bash
go build -o syspulse .
./syspulse
```

## Docker

Build and run with Docker:

```bash
docker build -t syspulse:latest .
docker run --rm -p 8080:8080 --env-file .env syspulse:latest
```

Use compose:

```bash
docker compose up --build -d
docker compose logs -f syspulse
```

The compose setup reads `.env`, exposes `${PORT:-8080}`, and includes a healthcheck against `/health`.

## Commands

- `ping` / `!ping` - liveness check
- `sys` / `!sys` - host system stats snapshot
- `rss <feed>` / `!rss <feed>` - fetch RSS (`tech`, `news`, `gaming`)

## Health Check

```bash
curl -i http://localhost:8080/health
```

Expected: HTTP `200 OK`.

## Interactions URL Setup

Configure your Discord Interactions Endpoint URL as:

```text
https://your-domain-or-tunnel/interactions
```

For local development, expose your local port with a tunnel and point Discord to that public URL.

## Troubleshooting

- `Missing required .env variable: ...`: the env var is not present in process environment
- `invalid request signature`: `PUBLIC_KEY` is incorrect or request is not from Discord
- Prefix commands not responding: verify Message Content Intent + `!` prefix usage
- Slash command changes not visible: restart the app so commands are re-registered

