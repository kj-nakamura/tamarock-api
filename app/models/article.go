package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// Article is table
type Article struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `json:"title"`
	Text        string         `gorm:"text" json:"text"`
	Category    int            `json:"category"`
	Artists     []ArtistInfo   `gorm:"many2many:article_artist_infos;" json:"artists"`
	PublishedAt time.Time      `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

// ResponseArticleData is frontに返す形
type ResponseArticleData struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `json:"title"`
	Pictures    []Picture      `json:"pictures"`
	Text        string         `gorm:"text" json:"text"`
	Category    int            `json:"category"`
	Artists     []ArtistInfo   `gorm:"many2many:article_artist_infos;" json:"artists"`
	PublishedAt time.Time      `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

// RequestArticleData is フロントから受け取る形
type RequestArticleData struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `json:"title"`
	Pictures    []Picture      `json:"pictures"`
	Text        string         `gorm:"text" json:"text"`
	Category    int            `json:"category"`
	ArtistIds   []int          `json:"artist_ids"`
	PublishedAt string         `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

// Picture is 画像の受け取りと送信
type Picture struct {
	Src   string `json:"src"`
	Title string `json:"title"`
}

func migrateArticle() {
	DbConnection.AutoMigrate(&Article{})
}

const defaultPicture string = "https://www.pakutaso.com/shared/img/thumb/penfan_KP_2783_TP_V.jpg"

// S3の画像URL(CloudFront利用)
var imageURL string = "https://static-prod.tamarock.jp/thumb/"

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
	t, err := time.Parse("2006-01-02", requestArticleData.PublishedAt)
	if err != nil {
		log.Fatalf("time parse error: %v", err)
	}

	// 記事を保存
	article := Article{
		ID:          requestArticleData.ID,
		Title:       requestArticleData.Title,
		Text:        requestArticleData.Text,
		Category:    requestArticleData.Category,
		Artists:     artistInfos,
		PublishedAt: t,
	}
	result := DbConnection.Create(&article)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	if requestArticleData.Pictures != nil && requestArticleData.Pictures[0].Src != "" {
		err := uploadImageToLocal(requestArticleData.Pictures[0].Src, "jpeg", strconv.FormatInt(int64(article.ID), 10))
		if err != nil {
			log.Printf("local upload error: %v", err)
		}
	}

	return article
}

// UpdateArticle is 記事更新
func UpdateArticle(r *http.Request, id int) RequestArticleData {
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

	// 一度、全関連アーティストを削除
	DbConnection.Model(&article).Association("Artists").Find(&artistInfos)
	DbConnection.Model(&article).Association("Artists").Delete(artistInfos)

	// リクエストからアーティスト情報を取得
	if len(requestArticleData.ArtistIds) > 0 {
		DbConnection.Where(requestArticleData.ArtistIds).Find(&artistInfos)
	} else {
		var emptyArtistInfo []ArtistInfo
		artistInfos = emptyArtistInfo
	}

	// 記事を保存
	t, err := time.Parse("2006-01-02", requestArticleData.PublishedAt)

	if err != nil {
		log.Fatalf("time parse error: %v", err)
	}
	articleData := Article{
		ID:          requestArticleData.ID,
		Title:       requestArticleData.Title,
		Text:        requestArticleData.Text,
		Category:    requestArticleData.Category,
		PublishedAt: t,
		Artists:     artistInfos,
	}

	result := DbConnection.Updates(articleData)
	if result.Error != nil {
		fmt.Printf("update error: %v", result.Error)
	}

	// 写真アップロード 画像なし、デフォルト画像、S3URLの場合はアップロードしない。(base64のみ)
	// フォームに画像がある&デフォルトの画像ではない
	if requestArticleData.Pictures != nil && requestArticleData.Pictures[0].Src != defaultPicture {
		// ローカル
		err := uploadImageToLocal(requestArticleData.Pictures[0].Src, "jpeg", strconv.FormatInt(int64(article.ID), 10))
		if err != nil {
			log.Printf("local upload error: %v", err)
		}
	}

	return GetAdminArticle(int(article.ID))
}

// uploadImageToLocal is ローカルに画像をアップロードする
func uploadImageToLocal(imageBase64 string, fileExtension string, fileDir string) error {
	if !strings.Contains(imageBase64, "data:") {
		return fmt.Errorf("not picture")
	}
	b64data := imageBase64[strings.IndexByte(imageBase64, ',')+1:]

	data := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data))

	filePath := "./static/" + fileDir
	// ディレクトリがなければ作成
	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		os.Mkdir("static", 0777)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, 0777)
	}
	file, err := os.Create(filePath + "/thumb." + fileExtension)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, data)
	if err != nil {
		return err
	}
	file.Close()

	return nil
}

// DeleteArticle is 記事を1つ削除
func DeleteArticle(id int) {
	var article Article

	var artistInfos []ArtistInfo
	DbConnection.First(&article, id)
	DbConnection.Model(&article).Association("Artists").Find(&artistInfos)
	DbConnection.Model(&article).Association("Artists").Delete(artistInfos)

	DbConnection.Delete(&article)
}

// GetArticle is 引数のIDに合致した記事を返す
func GetArticle(id int) ResponseArticleData {
	// 記事を取得
	var article Article
	var artistInfos []ArtistInfo

	// 今日以前に投稿されている記事を返す
	DbConnection.Where("published_at < ? OR published_at IS NULL", time.Now()).First(&article, id)

	// アーティスト情報取得
	DbConnection.Model(&article).Association("Artists").Find(&artistInfos)

	responseArticleData := ResponseArticleData{
		ID:          article.ID,
		Pictures:    addPicture(int64(article.ID), ""),
		Title:       article.Title,
		Text:        article.Text,
		Category:    article.Category,
		Artists:     artistInfos,
		PublishedAt: article.PublishedAt,
	}

	return responseArticleData
}

// GetAdminArticle is 引数のIDに合致した記事を返す
func GetAdminArticle(id int) RequestArticleData {
	// 関連するアーティストを取得
	var article Article
	var artistInfos []ArtistInfo

	DbConnection.Where("published_at < ? OR published_at IS NULL", time.Now()).First(&article, id)
	DbConnection.Model(&article).Association("Artists").Find(&artistInfos)

	// レスポンス用データを形成
	var artistInfoIDs []int
	var artistIDs []string

	for _, artistInfo := range artistInfos {
		artistIDs = append(artistIDs, artistInfo.ArtistId)
	}
	DbConnection.Where("artist_id IN ?", artistIDs).Find(&artistInfos)
	for _, artistInfo := range artistInfos {
		artistInfoIDs = append(artistInfoIDs, int(artistInfo.ID))
	}

	requestArticleData := RequestArticleData{
		ID:          article.ID,
		Pictures:    addPicture(int64(article.ID), defaultPicture),
		Title:       article.Title,
		Text:        article.Text,
		Category:    article.Category,
		ArtistIds:   artistInfoIDs,
		PublishedAt: article.PublishedAt.Format("2006-01-02"),
	}

	return requestArticleData
}

func GetAdminArticles(start int, end int, order string, sort string, query string, column string) []Article {
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

// GetArticles is 記事を複数返す
func GetArticles(start int, end int, order string, sort string, query string, column string) []ResponseArticleData {
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
			if query != "" {
				DbConnection.Where("published_at < ? OR published_at IS NULL", time.Now()).Order(createdOrder).Offset(start).Limit(limit).Where("title LIKE?", "%"+query+"%").Find(&articles)
			} else {
				DbConnection.Where("published_at < ? OR published_at IS NULL", time.Now()).Order(createdOrder).Offset(start).Limit(limit).Find(&articles)
			}
		} else {
			DbConnection.Where("published_at < ? OR published_at IS NULL", time.Now()).Order(createdOrder).Offset(start).Limit(limit).Where(column+" = ?", query).Find(&articles)
		}
	} else {
		DbConnection.Where("published_at < ? OR published_at IS NULL", time.Now()).Find(&articles)
	}

	var responseArticleDatas []ResponseArticleData
	for _, article := range articles {
		responseArticleData := ResponseArticleData{
			ID:          article.ID,
			Pictures:    addPicture(int64(article.ID), ""),
			Title:       article.Title,
			Text:        article.Text,
			Category:    article.Category,
			PublishedAt: article.PublishedAt,
		}

		responseArticleDatas = append(responseArticleDatas, responseArticleData)
	}

	return responseArticleDatas
}

// CountArticle is 全記事数を取得
func CountArticle(query string) int {
	var articles []Article
	DbConnection.Where("title LIKE?", "%"+query+"%").Find(&articles)

	return len(articles)
}

// private function

func getThumbnail(src string, articleID int64) string {
	IDStr := strconv.FormatInt(articleID, 10)
	// ローカル
	localFilePath := "./static/" + IDStr + "/thumb.jpeg"
	filePath := "http://localhost:5000/static/" + IDStr + "/thumb.jpeg"
	if _, err := os.Stat(localFilePath); !os.IsNotExist(err) {
		src = filePath
	}

	return src
}

func addPicture(articleID int64, defaultPicture string) []Picture {
	src := getThumbnail(defaultPicture, articleID)

	picture := Picture{
		Src:   src,
		Title: "thumbnail",
	}
	var pictures []Picture
	return append(pictures, picture)
}
