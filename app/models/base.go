package models

import (
	"api/config"
	"log"

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
	dsn := "host=postgres user=" + config.Env.DbUserName + " password=" + config.Env.DbPassword + " dbname=" + config.Env.DbName + " port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	DbConnection, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	migrateAdminUser()
	migrateArticle()
	migrateCategory()
	migrateYoutube()
	migrateArtistInfo()
}
