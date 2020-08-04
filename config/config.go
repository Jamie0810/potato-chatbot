package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server        Server        `mapstructure:"server"`
	Pubsub        Pubsub        `mapstructure:"pubsub"`
	Database      Database      `mapstructure:"database"`
	PotatoChatBot PotatoChatBot `mapstructure:"potato_chat_bot"`
	Log           Log           `mapstructure:"log"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type Server struct {
	Port string `mapstructure:"port"`
}
type Pubsub struct {
	ProjectID           string `mapstructure:"project_id"`
	TopicID             string `mapstructure:"topic_id"`
	CredentialsFilePath string `mapstructure:"credentials_file_path"`
}

type Database struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Port         uint   `mapstructure:"port"`
	Driver       string `mapstructure:"driver"`
	InstanceName string `mapstructure:"instance_name"`
}

type PotatoChatBot struct {
	APIURL     string `mapstructure:"api_url"`
	WebhookURL string `mapstructure:"webhook_url"`
	BotToken   string `mapstructure:"bot_token"`
}

func InitConfig(configPath string) (config Config, err error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	/* default */
	v.SetDefault("log_level", "INFO")
	v.SetDefault("log_format", "console")

	defaultPath := `./configs`

	if configPath == "" {
		configPath = defaultPath
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AddConfigPath(configPath)

	files, _ := ioutil.ReadDir(configPath)
	index := 0

	for _, file := range files {
		if filepath.Ext("./"+file.Name()) != ".yaml" && filepath.Ext("./"+file.Name()) != ".yml" {
			continue
		}

		v.SetConfigName(file.Name())
		var err error
		if index == 0 {
			err = v.ReadInConfig()
		} else {
			err = v.MergeInConfig()
		}
		if err == nil {
			index++
		}
	}

	if err = v.Unmarshal(&config); err != nil {
		return
	}

	return
}
