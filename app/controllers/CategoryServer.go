package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"api/app/models"
)

// admin category
func getAdminCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	start, _ := strconv.Atoi(r.URL.Query().Get("_start"))
	end, _ := strconv.Atoi(r.URL.Query().Get("_end"))
	order := r.URL.Query().Get("_order")
	sort := r.URL.Query().Get("_sort")
	query := r.URL.Query().Get("q")
	categories := models.GetCategories(start, end, order, sort, query)
	categoryCount := models.CountCategory(query)

	w.Header().Set("X-Total-Count", strconv.Itoa(categoryCount))
	responseJSON(w, categories)
}

func getAdminCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	article := models.GetAdminCategory(ID)

	responseJSON(w, article)
}

func createAdminCategoryHandler(w http.ResponseWriter, r *http.Request) {
	category := models.CreateCategory(r)

	responseJSON(w, category)
}

func updateAdminCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	category := models.UpdateCategory(r, ID)

	responseJSON(w, category)
}

func deleteAdminCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	models.DeleteCategory(ID)

	// 一覧を返す
	getAdminCategoriesHandler(w, r)
}
