package models

import (
	"api/config"
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"

	// orm

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DbConnection is for using global
var DbConnection *gorm.DB

func init() {
	var err error
	// mysql
	// dsn := config.Env.DbUserName + ":" + config.Env.DbPassword + "@tcp(" + config.Env.DbHost + ":3306)/" + config.Env.DbName + "?charset=utf8&parseTime=true"
	// DbConnection, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// postgres
	if config.Env.Env == "prod" {
		//prod
		url := os.Getenv("DATABASE_URL")
		connection, err := pq.ParseURL(url)
		if err != nil {
			panic(err.Error())
		}
		connection += " sslmode=require"
		sqlDB, err := sql.Open("postgres", connection)
		if err != nil {
			log.Fatalln(err)
		}

		DbConnection, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		// dev
		dsn := "host=" + config.Env.DbHost + "user=" + config.Env.DbUserName + " password=" + config.Env.DbPassword + " dbname=" + config.Env.DbName + " port=5432 sslmode=disable TimeZone=Asia/Tokyo"
		DbConnection, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Fatalln(err)
	}

	migrateAdminUser()
	migrateArticle()
	migrateCategory()
	migrateYoutube()
	migrateArtistInfo()
}
