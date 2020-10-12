package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Article is table
type Article struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	Title     string     `json:"title"`
	Text      string     `gorm:"text" json:"text"`
	Category  int        `json:"category"`
	CreatedAt time.Time  `json:"createdat"`
	UpdatedAt time.Time  `json:"updatedat"`
	DeletedAt *time.Time `json:"deletedat"`
}

func migrateArticle() {
	DbConnection.AutoMigrate(&Article{})
	DbConnection.Model(&Article{}).ModifyColumn("text", "text")
}

func CreateArticle(r *http.Request) Article {
	var article Article
	dec := json.NewDecoder(r.Body)
	for err := dec.Decode(&article); err != nil && err != io.EOF; {
		log.Println("ERROR: " + err.Error())
	}

	result := DbConnection.Create(&article)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return article
}

func UpdateArticle(r *http.Request, id int) Article {
	var article Article

	dec := json.NewDecoder(r.Body)
	for err := dec.Decode(&article); err != nil && err != io.EOF; {
		log.Println("ERROR: " + err.Error())
	}

	result := DbConnection.Table("articles").Where("id = ?", id).Update(map[string]interface{}{
		"title":    article.Title,
		"text":     article.Text,
		"category": article.Category,
	})

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return article
}

func DeleteArticle(id int) {
	var article Article

	DbConnection.Delete(&article, id)
}

// GetArticle is 引数のIDに合致した記事を返す
func GetArticle(id int) Article {
	var article Article

	DbConnection.First(&article, id)

	return article
}

// GetArticles is 記事を複数返す
func GetArticles() []Article {
	var articles []Article
	DbConnection.Order("created_at desc").Find(&articles)

	return articles
}
