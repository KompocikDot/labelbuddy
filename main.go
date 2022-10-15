package main

import (

	// "github.com/kompocikdot/labelbuddy/src/restocks"
	"os"

	"github.com/joho/godotenv"
	"github.com/kompocikdot/labelbuddy/src/discord"
)

func main() {
	godotenv.Load(".env")
	discord.RunBot(os.Getenv("BOT_TOKEN"))
	// items := restocks.RetrieveSalesLinks()
	// fmt.Println(items)
}