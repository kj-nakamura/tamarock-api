package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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

// SecretValue secret managerから取得
type SecretValue struct {
	YOUTUBE_API_KEY string
	API_KEY         string
	API_SECRET      string
	DB_PASSWORD     string
	DB_HOST         string
	S3AK            string
	S3SK            string
	BUCKET_NAME     string
}

var Env EnvValue
var Config ConfigValue

func init() {
	if err := envconfig.Process("", &Env); err != nil {
		log.Fatalf("[ERROR] Failed to process env: %s", err.Error())
	}

	if Env.Env == "prod" {
		var secretValue SecretValue
		secretStr := getSecret()

		json.Unmarshal([]byte(secretStr), &secretValue)
		Env.YoutubeApiKey = secretValue.YOUTUBE_API_KEY
		Env.ApiKey = secretValue.API_KEY
		Env.ApiSecret = secretValue.API_SECRET
		Env.DbPassword = secretValue.DB_PASSWORD
		Env.DbHost = secretValue.DB_HOST
		Env.S3AK = secretValue.S3AK
		Env.S3SK = secretValue.S3SK
		Env.BucketName = secretValue.BUCKET_NAME
	}

	Config.LogFile = "tamarock.log"
	Config.SQLDriver = "mysql"
	Config.Port = 5000
}

func getSecret() string {
	secretName := "tamarock/prod/"
	region := "ap-northeast-1"

	svc := secretsmanager.New(session.New(), aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return err.Error()
	}

	var secretString, decodedBinarySecret string

	if result.SecretString != nil {
		secretString = *result.SecretString

		return secretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			return "Base64 Decode Error"
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])

		return decodedBinarySecret
	}
}
