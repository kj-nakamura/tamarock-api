package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// Article is table
type Article struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title"`
	Text      string         `gorm:"text" json:"text"`
	Category  int            `json:"category"`
	Artists   []ArtistInfo   `gorm:"many2many:article_artist_infos;" json:"artists"`
	CreatedAt time.Time      `json:"createdat"`
	UpdatedAt time.Time      `json:"updatedat"`
	DeletedAt gorm.DeletedAt `json:"deletedat"`
}

type RequestArticleData struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title"`
	Text      string         `gorm:"text" json:"text"`
	Category  int            `json:"category"`
	ArtistIds []int          `json:"artist_ids"`
	CreatedAt time.Time      `json:"createdat"`
	UpdatedAt time.Time      `json:"updatedat"`
	DeletedAt gorm.DeletedAt `json:"deletedat"`
}

func migrateArticle() {
	DbConnection.AutoMigrate(&Article{})
}

// CreateArticle is 記事を作成する
func CreateArticle(r *http.Request) Article {
	// リクエストをjsonに変える
	var requestArticleData RequestArticleData
	dec := json.NewDecoder(r.Body)
	for err := dec.Decode(&requestArticleData); err != nil && err != io.EOF; {
		log.Println("article ERROR: " + err.Error())
	}

	// リクエストからアーティスト情報を取得
	var artistInfos []ArtistInfo
	if len(requestArticleData.ArtistIds) > 0 {
		DbConnection.Where(requestArticleData.ArtistIds).Find(&artistInfos)
	}

	// 記事を保存
	article := Article{
		ID:       requestArticleData.ID,
		Title:    requestArticleData.Title,
		Text:     requestArticleData.Text,
		Category: requestArticleData.Category,
		Artists:  artistInfos,
	}
	result := DbConnection.Create(&article)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return article
}

func UpdateArticle(r *http.Request, id int) Article {
	// リクエストをjsonに変える
	var requestArticleData RequestArticleData
	dec := json.NewDecoder(r.Body)
	for err := dec.Decode(&requestArticleData); err != nil && err != io.EOF; {
		log.Println("article ERROR: " + err.Error())
	}

	// 連携を一度解除する
	var article Article
	var artistInfos []ArtistInfo
	DbConnection.First(&article, id)
	DbConnection.Model(&article).Association("Artists").Find(&artistInfos)
	DbConnection.Model(&article).Association("Artists").Delete(artistInfos)

	// リクエストからアーティスト情報を取得
	DbConnection.Where(requestArticleData.ArtistIds).Find(&artistInfos)

	// 記事を保存
	articleData := Article{
		ID:       requestArticleData.ID,
		Title:    requestArticleData.Title,
		Text:     requestArticleData.Text,
		Category: requestArticleData.Category,
		Artists:  artistInfos,
	}
	result := DbConnection.Updates(articleData)
	if result.Error != nil {
		fmt.Printf("update error: %s", result.Error)
	}

	return article
}

func DeleteArticle(id int) {
	var article Article

	var artistInfos []ArtistInfo
	DbConnection.First(&article, id)
	DbConnection.Model(&article).Association("Artists").Find(&artistInfos)
	DbConnection.Model(&article).Association("Artists").Delete(artistInfos)

	DbConnection.Delete(&article)
}

// GetArticle is 引数のIDに合致した記事を返す
func GetArticle(id int) Article {
	// 記事を取得
	var article Article
	var artistInfos []ArtistInfo
	DbConnection.First(&article, id)

	// アーティスト情報取得
	DbConnection.Model(&article).Association("Artists").Find(&artistInfos)
	article.Artists = artistInfos

	return article
}

// GetAdminArticle is 引数のIDに合致した記事を返す
func GetAdminArticle(id int) RequestArticleData {
	// 関連するアーティストを取得
	var article Article
	var artistInfos []ArtistInfo
	DbConnection.First(&article, id)
	DbConnection.Model(&article).Association("Artists").Find(&artistInfos)

	// レスポンス用データを形成
	var artistData []int
	for _, artistInfo := range artistInfos {
		artistData = append(artistData, int(artistInfo.ID))
	}
	requestArticleData := RequestArticleData{
		ID:        article.ID,
		Title:     article.Title,
		Text:      article.Text,
		Category:  article.Category,
		ArtistIds: artistData,
	}

	return requestArticleData
}

// GetArticles is 記事を複数返す
func GetArticles(start int, end int, order string, sort string, query string, column string) []Article {
	var articles []Article

	if end > 0 {
		sortColumn := sort
		if sort == "" {
			sortColumn = "id"
		}
		createdOrder := sortColumn + " asc"
		if order == "DESC" {
			createdOrder = sortColumn + " desc"
		}
		limit := end - start
		if column == "" {
			DbConnection.Order(createdOrder).Offset(start).Limit(limit).Where("title LIKE?", "%"+query+"%").Find(&articles)
		} else {
			DbConnection.Order(createdOrder).Offset(start).Limit(limit).Where(column+" = ?", query).Find(&articles)
		}
	} else {
		DbConnection.Find(&articles)
	}

	return articles
}

// CountArticle is 全記事数を取得
func CountArticle(query string) int {
	var articles []Article
	DbConnection.Where("title LIKE?", "%"+query+"%").Find(&articles)

	return len(articles)
}
