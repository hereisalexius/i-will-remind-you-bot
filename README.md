# i-will-remind-you-bot
---
## About
A Telegram bot that AGGRESSIVELY (each 5s till you dismiss it) sends you notifications about an event so you don't miss it.

Bot URL : https://t.me/i_will_remind_you_bot

## How to use
Please use commands below to manage your notification
- /start - obviously, to start using bot
- /set - you will be prompted to questionnaire. 
  - Question #1: What I should remind you? - Just type what you want to be reminded of. *(example: Turn off oven)*
  - Question #2: After how long to warn you? - Provide duration for time when you should be notified *(example: 30m)*
- /dismiss - to dismiss notification on any stage
- /ping - to check if notification was set
- /help - to see help

## How to build and run

### Prerequisites (for all installation methods)
- *Golang* should be installed
- Bot token should be assigned to system env variable with name `TELE_TOKEN`

### Run release bin file
1. `./build/i-will-remind-you-bot start`
### Local build and run
1. Just build using `go build`
2. Run with `./i-will-remind-you-bot start`
### Docker
1. Build image `docker build -t i-will-remind-you-bot . `
2. Run container `docker run -d --name i-will-remind-you-bot-container -e TELE_TOKEN=$TELE_TOKEN i-will-remind-you-bot`
3. Stop and remove container`docker stop i-will-remind-you-bot-container | docker rm i-will-remind-you-bot-container`

**For deployments on Heroku use `heroku` branch*
**Heroku is not for free anymore* 😥
