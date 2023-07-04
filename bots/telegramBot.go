package bots

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type telegram struct {
	token   string
	updates map[string]interface{}
}

type chat struct {
	id                  int
	firstName, lastName string
	chatType            string
}

type message struct {
	room      chat
	messageId int32
	text      string
	author    author
}

type telegramUpdate struct {
	updateId int32
	date     int32
	message  message
}

type author struct {
	firstName string
	lastName  string
	id        int32
	isBot     bool
	username  string
}

func (bot *telegram) Initialize(token string) {
	bot.token = token
}

func (bot *telegram) SendMessage(message string) {
	fmt.Printf("enviando mensagem: %s\n", message)
}
func (bot *telegram) PullUpdates() {
	updates, err := sendAction(bot.token, "getUpdates")
	if err != nil {
		panic("coud not get updates")
	}
	bot.updates = updates
}
func (bot *telegram) Println() {
	fmt.Printf("%T", bot)
}

func sendAction(token, action string) (map[string]interface{}, error) {

	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, action)
	resp, err := http.Get(url)
	if err != nil {
		return map[string]interface{}{}, err
	}

	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{}, err
	}
	fmt.Println(string(res))
	var m map[string]interface{}
	err = json.Unmarshal(res, &m)
	if err != nil {
		return map[string]interface{}{}, err
	}
	fmt.Println(m)
	for _, up := range m["result"].([]interface{}) {
		var update telegramUpdate
		for key, value := range up.(map[string]interface{}) {
			fmt.Println(key)
			switch key {
			case "message":
				{
					v := value.(map[string]interface{})
					id, ok := v["message_id"].(float64)
					if !ok {
						panic("uepa1")
					}
					text, ok := v["text"].(string)
					if !ok {
						panic("uepa2")
					}
					authorM, ok := v["from"].(map[string]interface{})
					if !ok {
						panic("uepa3")
					}
					update.message = message{
						messageId: int32(id),
						text:      text,
						author: author{
							id:        int32(authorM["id"].(float64)),
							firstName: authorM["first_name"].(string),
							lastName:  authorM["last_name"].(string),
							isBot:     authorM["is_bot"].(bool),
							username:  authorM["username"].(string),
						},
					}
				}
			case "update_id":
				{
					update.updateId = int32(value.(float64))
				}
			case "date":
				{
					update.date = value.(int32)
				}
			}
			fmt.Println(update)
		}
	}
	return m, nil
}
