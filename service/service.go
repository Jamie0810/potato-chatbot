package service

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/config"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/pkg/log"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/potato"
	"golang.org/x/sync/errgroup"
)

type ServiceController struct {
	config config.Config
	logger *log.Logger
	db     *gorm.DB
	potato potato.PotatoChatBot
	client potato.Client
}

func InitServiceController(config config.Config, logger *log.Logger, db *gorm.DB) *ServiceController {
	client := newClient(config)

	svc := &ServiceController{
		config: config,
		logger: logger,
		db:     db,
		client: client,
	}
	return svc
}

func newClient(config config.Config) potato.Client {
	client := potato.InitChatBotClient(
		config.PotatoChatBot.APIURL,
		config.PotatoChatBot.BotToken,
		nil,
	)
	return client
}

// MainService main service
func (svc *ServiceController) MainService() {
	//Service1: Provide instructions by receiving commands from client
	g := errgroup.Group{}
	g.Go(func() error {
		svc.webhookHandling(svc.client)
		return nil
	})

	// Service2: Send notifications when an order was not able to create successfully
	err := svc.sendNotifications()
	if err != nil {
		errors.Wrap(err, "Failed to send notifications")
	}

	g.Wait()
}
