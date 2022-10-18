package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kompocikdot/labelbuddy/src/discord"
)

func main() {
	godotenv.Load(".env")
	discord.RunBot(os.Getenv("BOT_TOKEN"))
}