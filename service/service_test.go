package service

import (
	"context"
	"testing"

	"gitlab.silkrode.com.tw/golang/mq"
	"gitlab.silkrode.com.tw/golang/mq/inside/pub"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/config"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/potato"
)

var cnf config.Config

// const CredentialsFilePath = "/Users/jamie/Desktop/go-projects/src/kbc/km/chatbot/prod_pubsub_credentials.json"
const CredentialsFilePath = "/Users/jamie/Desktop/go-projects/src/kbc/km/chatbot/pubsub_credentials.json"

func TestMain(m *testing.M) {
	cnf, _ = config.InitConfig("../config")
	m.Run()
}

func TestPub(t *testing.T) {
	ctx := context.Background()
	pubInstance, err := mq.Init(ctx, cnf.Pubsub.TopicID, CredentialsFilePath, cnf.Pubsub.ProjectID, mq.InitPub())
	if err != nil {
		t.Fatalf("failed to create pubInstance: %v", err)
	}

	// order := "📣Congrats! An order just created failed. 🍎😳"
	order := "Jamie主题：风控异常预警⚠️ \n商户：KB000001 王公司"

	pubInstance.Publisher().Options(
		pub.SetErrorHook(func(err error, requestID string) {
			t.Errorf("failed to publish requestID: %v ; message: %v", requestID, err)
		}),
	).Publish([]byte(order), "requestID")
}

func TestWebhook(t *testing.T) {
	client := potato.InitChatBotClient(
		cnf.PotatoChatBot.APIURL,
		cnf.PotatoChatBot.BotToken,
		nil,
	)

	// set webhook to recieve messages from potato chat
	_ = client.SetWebhook(potato.Webhook{
		URL: "localhost:3000" + "/method",
	})
}

func TestBroadcastToGroups(t *testing.T) {
	client := potato.InitChatBotClient(
		cnf.PotatoChatBot.APIURL,
		cnf.PotatoChatBot.BotToken,
		nil,
	)

	text := "訂單編號:123\n商戶:abc"
	msg := potato.MessageInfo{
		Msg:     text,
		ChatIDs: []int64{12105639},
	}

	err := client.BroadcastToGroups(cnf, msg)
	if err != nil {
		t.Fatal(err)
	}
}
