package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type EnvValue struct {
	Env           string `required:"true" default:"dev"`
	YoutubeApiKey string `split_words:"true"`
	ApiKey        string `split_words:"true"`
	ApiSecret     string `split_words:"true"`
	DbName        string `required:"true" split_words:"true"`
	DbHost        string `split_words:"true" default:"mysql"`
	DbUserName    string `required:"true" split_words:"true"`
	DbPassword    string `split_words:"true"`
	S3AK          string
	S3SK          string
	BucketName    string `split_words:"true"`
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
