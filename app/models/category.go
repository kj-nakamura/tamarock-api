package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// Category is table
type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Articles  []Article      `gorm:"foreignKey:Category" json:"articles"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// UpdateCategoryData is request data
type UpdateCategoryData struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func migrateCategory() {
	DbConnection.AutoMigrate(&Category{})
}

// CreateCategory is カテゴリ作成
func CreateCategory(r *http.Request) Category {
	// リクエストをjsonに変える
	var category Category
	dec := json.NewDecoder(r.Body)
	for err := dec.Decode(&category); err != nil && err != io.EOF; {
		log.Println("category ERROR: " + err.Error())
	}

	// カテゴリを保存
	categoryData := map[string]interface{}{
		"Name":      category.Name,
		"CreatedAt": time.Now(),
		"UpdatedAt": time.Now(),
	}
	result := DbConnection.Model(&category).Create(categoryData)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	DbConnection.Last(&category)
	return category
}

// UpdateCategory is カテゴリ更新
func UpdateCategory(r *http.Request, id int) Category {
	// リクエストをjsonに変える
	var updateCategoryData UpdateCategoryData
	dec := json.NewDecoder(r.Body)
	for err := dec.Decode(&updateCategoryData); err != nil && err != io.EOF; {
		log.Println("article ERROR: " + err.Error())
	}

	// カテゴリを保存
	var category Category
	result := DbConnection.Model(&category).Where("id = ?", uint(id)).Update("name", updateCategoryData.Name)
	if result.Error != nil {
		fmt.Printf("update error: %s", result.Error)
	}

	return category
}

// DeleteCategory is カテゴリ削除
func DeleteCategory(id int) {
	var category Category
	DbConnection.First(&category, id)

	DbConnection.Delete(&category)
}

// GetAdminCategory is 引数のIDに合致した記事を返す
func GetAdminCategory(id int) UpdateCategoryData {
	// 関連する記事を取得
	var category Category
	var articles []Article
	DbConnection.First(&category, id)
	DbConnection.Model(&category).Association("Articles").Find(&articles)

	requestCategoryData := UpdateCategoryData{
		ID:   category.ID,
		Name: category.Name,
	}

	return requestCategoryData
}

// GetCategories is カテゴリを複数返す
func GetCategories(start int, end int, order string, sort string, query string) []Category {
	var categories []Category

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
		DbConnection.Order(createdOrder).Offset(start).Limit(limit).Where("name LIKE?", "%"+query+"%").Find(&categories)
	} else {
		DbConnection.Find(&categories)
	}

	return categories
}

// CountCategory is 全カテゴリ数を取得
func CountCategory(query string) int {
	var categories []Category
	DbConnection.Where("name LIKE?", "%"+query+"%").Find(&categories)

	return len(categories)
}
