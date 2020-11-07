package models

import (
	"log"

	// orm
	"api/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DbConnection is for using global
var DbConnection *gorm.DB

func init() {
	var err error
	dsn := config.Env.DbUserName + ":" + config.Env.DbPassword + "@tcp(" + config.Env.DbHost + ":3306)/" + config.Env.DbName + "?charset=utf8&parseTime=true"
	DbConnection, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// DbConnection, err = gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		log.Fatalln(err)
	}

	// migrateAdminUser()
	// migrateArtistInfo()
	// migrateArticle()
	// migrateCategory()
}
