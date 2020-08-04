package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/config"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/pkg/database"
	internal "gitlab.silkrode.com.tw/team_golang/KM/chatbot/pkg/log"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/service"
)

func main() {
	config, err := config.InitConfig("./config")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to init config."))
	}

	logger, err := internal.InitLogger(config)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to init logger."))
	}

	db, err := database.InitDB(config)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to init db."))
	}

	s := service.InitServiceController(config, logger, db)
	s.MainService()

	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	err = e.Start(":" + config.Server.Port)
	if err == nil {
		logger.InfoMsg("listening port on: ", config.Server.Port)
	}
}
