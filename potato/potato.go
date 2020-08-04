package potato

import (
	"net/http"

	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/config"
)

const (
	// UserChat user chat
	UserChat ChatType = 1
	// StandardGroupChat standard group chat
	StandardGroupChat ChatType = 2
	// SuperGroupChat super group chat
	SuperGroupChat ChatType = 3
)

type ChatType int64

var Potato *PotatoChatBot

type PotatoChatBot struct {
	apiURL string
	token  string
	client *http.Client
}

type Client interface {
	SetWebhook(opts Webhook) (err error)
	BroadcastToGroups(config.Config, MessageInfo) (err error)
	SendMessages([]byte) (err error)
}
type Webhook struct {
	URL string `json:"url"`
}

type MessageInfo struct {
	Msg     string
	ChatIDs []int64
}

type SendTextMessage struct {
	ChatType         ChatType `json:"chat_type"`           // Required: Type for the target chat
	ChatID           int64    `json:"chat_id"`             // Required: Unique identifier for the target chat
	Text             string   `json:"text"`                // Required: Text of the message to be sent
	Markdown         bool     `json:"markdown"`            // Optional: Whether to use MarkDown rendering
	ReplyToMessageID int64    `json:"reply_to_message_id"` // Optional: If the message is a reply, then it'd be the ID of the original message
}

type Groups struct {
	Items struct {
		Channels    []Peer `json:"Channels"`
		Groups      []Peer `json:"Groups"`
		SuperGroups []Peer `json:"SuperGroups"`
	} `json:"result"`
}

type Peer struct {
	ID   int64  `json:"PeerID"`
	Name string `json:"PeerName"`
}
