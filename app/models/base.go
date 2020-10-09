package models

import (
	"log"

	// orm
	"api/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	TableNameArtists = "artist_infos"
)

// DbConnection is for using global
var DbConnection *gorm.DB

func init() {
	var err error
	DbConnection, err = gorm.Open(config.Config.SQLDriver, config.Env.DbUserName+":"+config.Env.DbPassword+"@tcp("+config.Env.DbHost+":3306)/"+config.Env.DbName+"?charset=utf8&parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	migrateAdminUser()
	migrateArtistInfo()
	migrateArticle()
}
