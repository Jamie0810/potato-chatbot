package service

import (
	"context"

	"gitlab.silkrode.com.tw/golang/mq"
	"gitlab.silkrode.com.tw/golang/mq/inside/sub"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/model"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/potato"

	"github.com/pkg/errors"
)

func (svc *ServiceController) sendNotifications() error {
	// subscribe to messages from the topic
	ctx := context.Background()
	subInstance, err := mq.Init(
		ctx,
		svc.config.Pubsub.TopicID,
		svc.config.Pubsub.CredentialsFilePath,
		svc.config.Pubsub.ProjectID,
		mq.InitSub(),
	)
	if err != nil {
		svc.logger.Error().Msgf("%+v\n", errors.Wrap(err, "Failed to create subInstance"))
		return err
	}

	subInstance.Subscriber().
		Options(
			sub.SyncMode(),
			sub.SetErrorHook(func(err error, msgData string, msgID string) {
				svc.logger.Error().Msgf("%+v\n", err, "Failed to pull messages")
				return
			}),
		).
		Subscribe(func(ctx context.Context, msg []byte, msgId string) error {
			svc.logger.InfoMsg("Received a message from the topic. Message ID: ", msgId)
			svc.logger.InfoMsg("Message: ", string(msg))
			// get chatroom ids from db
			chatRooms, err := svc.getChatRoomIDs()
			if err != nil {
				return err
			}

			if err := svc.invokeChatBot(string(msg), chatRooms); err != nil {
				return err
			}

			svc.logger.InfoMsg("The message has been sent to groups successfully!")
			return nil
		})
	return nil
}

func (svc *ServiceController) getChatRoomIDs() ([]int64, error) {
	var chatRooms []model.ChatRoom
	if err := svc.db.Select("chat_id").Find(&chatRooms).Error; err != nil {
		return nil, errors.Wrap(err, "Failed to select chat ids")
	}
	var ids []int64
	for _, v := range chatRooms {
		ids = append(ids, v.ChatID)
	}

	return ids, nil
}

func (svc *ServiceController) invokeChatBot(text string, chatroomIDs []int64) error {
	if text == "" {
		return errors.New("no message to notify")
	}

	msg := potato.MessageInfo{
		Msg:     text,
		ChatIDs: chatroomIDs,
	}

	if err := svc.client.BroadcastToGroups(svc.config, msg); err != nil {
		return errors.Wrap(err, "Failed to broadcast the message to groups")
	}

	return nil
}
