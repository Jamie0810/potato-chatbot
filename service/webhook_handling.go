package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/model"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/potato"
)

type Response struct {
	Messages struct {
		Text string `json:"text"`
		Chat struct {
			ID   int64  `json:"id"`
			Name string `json:"title"`
		} `json:"chat"`
		User struct {
			ID   int64  `json:"id"`
			Name string `json:"username"`
		} `json:"from"`
		CreatedAt int64 `json:"date"`
	} `json:"message"`
}

func (svc *ServiceController) webhookHandling(client potato.Client) error {
	// Setup a HTTP server listening for webhooks
	err := svc.setWebhookServer(client)
	if err != nil {
		return errors.Wrap(err, "newWebhookServer")
	}
	return nil
}

func (svc *ServiceController) setWebhookServer(client potato.Client) error {
	// Set a webhook to recieve messages from client
	err := client.SetWebhook(potato.Webhook{
		URL: svc.config.PotatoChatBot.WebhookURL,
	})
	if err != nil {
		errors.Wrap(err, "Failed to set the webhook.")
		return err
	}

	http.HandleFunc("/", svc.getMessagesFromClient(client))
	err = http.ListenAndServe((":" + svc.config.Server.Port), nil)
	if err != nil {
		errors.Wrap(err, "ListenAndServe")
	}
	return nil
}

func (svc *ServiceController) getMessagesFromClient(client potato.Client) func(w http.ResponseWriter, resp *http.Request) {
	return func(w http.ResponseWriter, resp *http.Request) {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errors.Wrap(err, "Can not read messages sent by client")
		}

		var response []Response
		if err := json.Unmarshal(body, &response); err != nil {
			errors.Wrap(err, "Unmarshal err")
		}

		if len(response) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chatRoom := model.ChatRoom{
			ChatID:    response[0].Messages.Chat.ID,
			GroupName: response[0].Messages.Chat.Name,
			UserID:    response[0].Messages.User.ID,
			UserName:  response[0].Messages.User.Name,
		}

		if string([]rune(response[0].Messages.Text)[0]) != "/" {
			return
		}

		// add or remove chat ids to db with the command 'subscribe' or 'unsubscribe' accordingly
		text := strings.Split(response[0].Messages.Text, "/")[1]
		if text == "subscribe" {
			if svc.db.Where("chat_id = ?", response[0].Messages.Chat.ID).First(&chatRoom).Error != nil {
				svc.db.Create(&chatRoom)
				svc.logger.InfoMsg(response[0].Messages.User.Name, " has subscribed to alert notifications")
			}
		}

		if text == "unsubscribe" {
			svc.db.Where("chat_id = ?", response[0].Messages.Chat.ID).Delete(&chatRoom)
			svc.logger.InfoMsg(response[0].Messages.User.Name, " has unsubscribed from alert notifications")
		}

		// send instructions
		if text == "help" {
			instructions := "/subscribe: 訂閱預警服務\n/unsubscribe: 取消訂閱預警服務"

			// send instructions to user chat rooms
			userChatMsg := potato.SendTextMessage{
				ChatType: potato.UserChat,
				ChatID:   response[0].Messages.Chat.ID,
				Text:     instructions,
				Markdown: false,
			}

			body, err = json.Marshal(userChatMsg)
			if err != nil {
				errors.Wrap(err, "Marshal")
			}

			err = client.SendMessages(body)
			if err != nil {
				errors.Wrap(err, "SendMessages Err")
			}

			// send instructions to groups
			groupChatMsg := potato.SendTextMessage{
				ChatType: potato.StandardGroupChat,
				ChatID:   response[0].Messages.Chat.ID,
				Text:     instructions,
				Markdown: false,
			}

			body, err = json.Marshal(groupChatMsg)
			if err != nil {
				errors.Wrap(err, "Marshal")
			}

			err = client.SendMessages(body)
			if err != nil {
				errors.Wrap(err, "SendMessages Err")
			}
		}
	}
}
