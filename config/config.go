package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type EnvValue struct {
	Env           string `default:"dev"`
	YoutubeApiKey string `split_words:"true"`
	ApiKey        string `split_words:"true"`
	ApiSecret     string `split_words:"true"`
	DbName        string `split_words:"true" default:"tamarock"`
	DbHost        string `split_words:"true" default:"127.0.0.1"`
	DbUserName    string `split_words:"true" default:"root"`
	DbPassword    string `split_words:"true" default:"password"`
}

type ConfigValue struct {
	LogFile   string
	SQLDriver string
	Port      int
}

var Env EnvValue
var Config ConfigValue

func init() {
	if err := envconfig.Process("", &Env); err != nil {
		log.Fatalf("[ERROR] Failed to process env: %s", err.Error())
	}

	Config.LogFile = "tamarock.log"
	Config.SQLDriver = "mysql"
	Config.Port = 5000
}
