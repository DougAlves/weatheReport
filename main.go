package main

import (
	"bot-teste/bots"
	"os"

	"github.com/joho/godotenv"
)

func initEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Coud not read env file")
	}

}

func main() {
	initEnv()
	argsWithoutProg := os.Args[1:]
	token := os.Getenv("token")
	bot := bots.GetBot(argsWithoutProg[0])
	bot.Initialize(token)
	bot.Run()
	//initialize(token, "getUpdates")
	//initialize(token, "getUpdates")
}
