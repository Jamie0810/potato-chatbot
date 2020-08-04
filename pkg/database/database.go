package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/config"
	"gitlab.silkrode.com.tw/team_golang/KM/chatbot/model"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitDB(config config.Config) (db *gorm.DB, err error) {
	var connInfo string

	if config.Database.InstanceName == "" {
		connInfo = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local", config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.DBName)
	} else {
		connInfo = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?charset=utf8mb4&parseTime=true&loc=UTC&time_zone=UTC", config.Database.User, config.Database.Password, config.Database.InstanceName, config.Database.DBName)
	}

	db, err = gorm.Open("mysql", connInfo)
	if err != nil {
		return nil, err
	}

	chatbotRoom := model.ChatRoom{}
	db.AutoMigrate(&chatbotRoom)
	return db, nil
}
