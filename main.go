package main

import (
	"bot-teste/bots"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func run(runner bots.Bot) {
	for {
		runner.PullUpdates()
		time.Sleep(10 * 1000000)
		break
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Coud not read env file")
	}
	argsWithoutProg := os.Args[1:]
	token := os.Getenv("token")
	exc := bots.GetBot(argsWithoutProg[0])
	exc.Initialize(token)
	run(exc)
	//initialize(token, "getUpdates")
	//initialize(token, "getUpdates")
}
