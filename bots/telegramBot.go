package bots

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type telegram struct {
	token      string
	updates    []TelegramUpdate
	topMessage Message
}

type Chat struct {
	Id        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"Last_name"`
	ChatType  string `json:"type"`
}

type Message struct {
	Room      Chat   `json:"chat"`
	MessageId uint64 `json:"message_id"`
	Text      string `json:"text"`
	Date      uint64 `json:"date"`
	Author    Author `json:"from"`
}

type TelegramUpdate struct {
	UpdateId uint64  `json:"update_id"`
	Message  Message `json:"message"`
}

type TelegramUpdateDTO struct {
	Ok      bool             `json:"ok"`
	Updates []TelegramUpdate `json:"result"`
}

type SentMessageReturn struct {
	Ok          bool    `json:"ok"`
	SentMessage Message `json:"result"`
}
type Author struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Id        uint64 `json:"id"`
	IsBot     bool   `json:"is_bot"`
	Username  string `json:"username"`
}

func (bot *telegram) Initialize(token string) {
	bot.token = token
}

func (bot *telegram) SendMessage(message string) {
	url, m := fmt.Sprintf("https://api.telegram.org/bot%s/%s", bot.token, "sendMessage"), new(bytes.Buffer)
	body := map[string]interface{}{
		"chat_id": bot.topMessage.Room.Id,
		"text":    message,
	}

	json.NewEncoder(m).Encode(body)
	resp, err := http.Post(url, "application/json", m)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var messageSent SentMessageReturn
	err = json.Unmarshal(res, &messageSent)
	if err != nil {
		panic(err)
	}

	fmt.Println(messageSent)

}
func (bot *telegram) PullUpdates() {
	updates, err := sendAction(bot.token, "getUpdates")
	if err != nil {
		panic(err)
	}
	bot.updates = updates
}

func (bot telegram) Run() {
	for {
		bot.PullUpdates()
		updatesChannel := updatesToChannel(bot.updates)
		done := treatMessage(updatesChannel, bot)
		for processedUpdate := range done {
			processedUpdates = append(processedUpdates, processedUpdate)
		}
	}
}

var processedUpdates []TelegramUpdate

func updatesToChannel(updates []TelegramUpdate) <-chan TelegramUpdate {
	out := make(chan TelegramUpdate)
	go func() {
		for _, update := range updates {
			if !updateProcessed(update) {
				out <- update
			}
		}
		close(out)
	}()
	return out
}

func treatMessage(upChan <-chan TelegramUpdate, bot telegram) <-chan TelegramUpdate {
	out := make(chan TelegramUpdate)
	go func() {
		for update := range upChan {
			if update.Message.Text == "/ping" {
				bot.topMessage = update.Message
				bot.SendMessage("pong")
			}
			out <- update
		}
		close(out)
	}()
	return out
}

func updateProcessed(update TelegramUpdate) bool {
	for _, up := range processedUpdates {
		if up.UpdateId == update.UpdateId {
			return true
		}
	}
	return false
}

func (bot *telegram) Println() {
	fmt.Printf("%T", bot)
}

func sendAction(token, action string) ([]TelegramUpdate, error) {

	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, action)
	resp, err := http.Get(url)
	if err != nil {
		return []TelegramUpdate{}, err
	}

	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return []TelegramUpdate{}, err
	}
	var updates TelegramUpdateDTO
	err = json.Unmarshal(res, &updates)
	if err != nil {
		return []TelegramUpdate{}, err
	}
	return updates.Updates, nil
}
