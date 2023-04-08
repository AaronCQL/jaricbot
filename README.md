# JaricBot

**J**ust **A**nother **R**ather **I**ntelligent **C**hat **Bot** - written in Go and powered by Telegram and OpenAI.

## Deploying

1. Create a `.env` file in the root directory (refer to `.env.example`), and fill in the respective values
2. Start the bot for production using `make start`

## Developing

1. Install [`air`](https://github.com/cosmtrek/air) (used to rerun code automatically on file changes)
2. Create a `.env` file in the root directory (refer to `.env.example`), and fill in the respective values
3. Start the bot with automatic code reload using `make dev`
