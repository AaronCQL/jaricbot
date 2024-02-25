# JaricBot

**J**ust **A**nother **R**ather **I**ntelligent **C**hat **Bot** - Telegram chat bot powered by Gemini.

## Setup

1. Install Go 1.21 or higher
2. Install [`air`](https://github.com/cosmtrek/air) (optional; for dev only): `go install github.com/cosmtrek/air@v1.49.0`
3. Clone this repository: `git clone https://github.com/AaronCQL/jaricbot.git`
4. Create a `.env` file in the root directory and fill in your Telegram bot API key and Gemini API key (refer to [`.env.example`](./.env.example))

> If you want the bot to reply to messages in groups, you will need to disable [privacy mode](https://core.telegram.org/bots#privacy-mode) via [@BotFather](https://t.me/BotFather).

## Deploying

```sh
# Build the bot and start it
make start
```

## Developing

```sh
# Run the bot using `air` which auto-restarts on file changes
make dev
```
