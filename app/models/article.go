package models

import (
	"api/config"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

// ResponseArticleData is frontに返す形
type ResponseArticleData struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title"`
	Pictures  []Picture      `json:"pictures"`
	Text      string         `gorm:"text" json:"text"`
	Category  int            `json:"category"`
	Artists   []ArtistInfo   `gorm:"many2many:article_artist_infos;" json:"artists"`
	CreatedAt time.Time      `json:"createdat"`
	UpdatedAt time.Time      `json:"updatedat"`
	DeletedAt gorm.DeletedAt `json:"deletedat"`
}

// RequestArticleData is フロントから受け取る形
type RequestArticleData struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title"`
	Pictures  []Picture      `json:"pictures"`
	Text      string         `gorm:"text" json:"text"`
	Category  int            `json:"category"`
	ArtistIds []int          `json:"artist_ids"`
	CreatedAt time.Time      `json:"createdat"`
	UpdatedAt time.Time      `json:"updatedat"`
	DeletedAt gorm.DeletedAt `json:"deletedat"`
}

// Picture is 画像の受け取りと送信
type Picture struct {
	Src   string `json:"src"`
	Title string `json:"title"`
}

func migrateArticle() {
	// DbConnection.AutoMigrate(&Article{})
}

const defaultPicture string = "https://www.pakutaso.com/shared/img/thumb/penfan_KP_2783_TP_V.jpg"

// S3の画像URL
var s3ImageURL string = "https://" + config.Env.BucketName + ".s3-ap-northeast-1.amazonaws.com/thumb/"

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

	if requestArticleData.Pictures != nil && requestArticleData.Pictures[0].Src != "" {
		err := UploadToS3(requestArticleData.Pictures[0].Src, "jpeg", strconv.FormatInt(int64(article.ID), 10))
		if err != nil {
			log.Println(err)
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
	articleData := Article{
		ID:       requestArticleData.ID,
		Title:    requestArticleData.Title,
		Text:     requestArticleData.Text,
		Category: requestArticleData.Category,
		Artists:  artistInfos,
	}

	result := DbConnection.Updates(articleData)
	if result.Error != nil {
		fmt.Printf("update error: %v", result.Error)
	}

	// 写真アップロード 画像なし、デフォルト画像、S3URLの場合はアップロードしない。(base64のみ)
	if requestArticleData.Pictures != nil &&
		requestArticleData.Pictures[0].Src != defaultPicture &&
		strings.Contains(requestArticleData.Pictures[0].Src, s3ImageURL) == false {
		err := UploadToS3(requestArticleData.Pictures[0].Src, "jpeg", strconv.FormatInt(int64(article.ID), 10))
		if err != nil {
			log.Printf("error after upload: %v", err)
		}
	}

	return GetAdminArticle(int(article.ID))
}

// UploadToS3 is s3に画像アップロード
func UploadToS3(imageBase64 string, fileExtension string, filename string) error {
	// 環境変数からS3Credential周りの設定を取得
	bucketName := config.Env.BucketName
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(config.Env.S3AK, config.Env.S3SK, ""),
		Region:      aws.String("ap-northeast-1"),
	}))

	uploader := s3manager.NewUploader(sess)

	b64data := imageBase64[strings.IndexByte(imageBase64, ',')+1:]
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Fatalf("upload decode error: %v", err)
	}
	wb := new(bytes.Buffer)
	wb.Write(data)

	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String("thumb/" + filename + "." + fileExtension),
		Body:        wb,
		ContentType: aws.String("image/" + fileExtension),
	})

	if err != nil {
		fmt.Sprintf("upload error: %v", res)
		if err, ok := err.(awserr.Error); ok && err.Code() == request.CanceledErrorCode {
			log.Fatalf("upload answer error: %v", err)
		} else {
			return fmt.Errorf("Upload Failed %d", 400)
		}
	}

	return nil
}

// checkS3KeyExists is 画像がS3に存在するかチェックする
func checkS3KeyExists(objectKey string) bool {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(config.Env.S3AK, config.Env.S3SK, ""),
		Region:      aws.String("ap-northeast-1"),
	}))

	bucketName := config.Env.BucketName

	// S3 clientを作成
	svc := s3.New(sess)
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("thumb/" + objectKey + ".jpeg"),
	}

	_, err := svc.HeadObject(input)
	if err != nil {
		return false
	}
	return true
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
	DbConnection.First(&article, id)

	// アーティスト情報取得
	DbConnection.Model(&article).Association("Artists").Find(&artistInfos)

	src := ""
	IDStr := strconv.FormatInt(int64(article.ID), 10)
	if checkS3KeyExists(IDStr) {
		src = s3ImageURL + IDStr + ".jpeg"
	}
	picture := Picture{
		Src:   src,
		Title: "thumbnail",
	}
	var pictures []Picture
	pictures = append(pictures, picture)

	responseArticleData := ResponseArticleData{
		ID:       article.ID,
		Pictures: pictures,
		Title:    article.Title,
		Text:     article.Text,
		Category: article.Category,
		Artists:  artistInfos,
	}

	return responseArticleData
}

// GetAdminArticle is 引数のIDに合致した記事を返す
func GetAdminArticle(id int) RequestArticleData {
	// 関連するアーティストを取得
	var article Article
	var artistInfos []ArtistInfo
	DbConnection.First(&article, id)
	DbConnection.Model(&article).Association("Artists").Find(&artistInfos)

	src := defaultPicture
	IDStr := strconv.FormatInt(int64(article.ID), 10)
	if checkS3KeyExists(IDStr) {
		src = s3ImageURL + IDStr + ".jpeg"
	}
	picture := Picture{
		Src:   src,
		Title: "thumbnail",
	}
	var pictures []Picture
	pictures = append(pictures, picture)

	// レスポンス用データを形成
	var artistData []int
	for _, artistInfo := range artistInfos {
		artistData = append(artistData, int(artistInfo.ID))
	}
	requestArticleData := RequestArticleData{
		ID:        article.ID,
		Pictures:  pictures,
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
