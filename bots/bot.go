package bots

type Bot interface {
	Initialize(string)
	SendMessage(string)
	PullUpdates()
	Println()
	Run()
}

func GetBot(arg string) Bot {
	switch arg {
	case "telegram":
		bot := telegram{}
		return &bot
	}

	bot := telegram{}
	return &bot
}
