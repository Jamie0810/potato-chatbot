package potato

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/config"
)

func InitChatBotClient(apiURL string, botToken string, meta map[string]string) Client {
	client := PotatoChatBot{
		apiURL: fmt.Sprintf("%s/%s", apiURL, botToken),
		token:  botToken,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	Potato = &client
	return &client
}

func (bot *PotatoChatBot) SetWebhook(w Webhook) (err error) {
	body, err := json.Marshal(w)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, bot.apiURL+"/setWebhook", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := bot.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf(resp.Status)
	}
	return nil
}

func (bot *PotatoChatBot) SendMessages(msg []byte) (err error) {
	req, err := http.NewRequest("POST", bot.apiURL+"/sendTextMessage", bytes.NewReader(msg))
	if err != nil {
		return errors.Wrap(err, "Failed to send messages.")
	}
	req.Header.Add("Content-Type", "application/json")

	err = bot.DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (bot *PotatoChatBot) BroadcastToGroups(config config.Config, msgInfo MessageInfo) (err error) {
	// send msg to user chat rooms
	for _, chatID := range msgInfo.ChatIDs {
		Msg := SendTextMessage{
			ChatType: UserChat,
			ChatID:   chatID,
			Text:     msgInfo.Msg,
		}

		body, err := json.Marshal(Msg)
		if err != nil {
			return err
		}
		err = bot.SendMessages(body)
		if err != nil {
			return err
		}
	}

	// send msg to groups
	for _, chatID := range msgInfo.ChatIDs {
		Msg := SendTextMessage{
			ChatType: StandardGroupChat,
			ChatID:   chatID,
			Text:     msgInfo.Msg,
		}

		body, err := json.Marshal(Msg)
		if err != nil {
			return err
		}
		err = bot.SendMessages(body)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bot *PotatoChatBot) DoRequest(req *http.Request) (err error) {
	resp, err := bot.client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("status:%s -- %s", resp.Status, string(b))
	}

	return nil
}
