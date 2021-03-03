package controllers

import (
	"api/app/models"
	"fmt"
	"net/http"
	"strconv"
)

// web
func getArticlesHandler(w http.ResponseWriter, r *http.Request) {
	start, _ := strconv.Atoi(r.URL.Query().Get("_start"))
	end, _ := strconv.Atoi(r.URL.Query().Get("_end"))
	order := r.URL.Query().Get("_order")
	sort := r.URL.Query().Get("_sort")
	query := r.URL.Query().Get("q")
	column := r.URL.Query().Get("column")
	articles := models.GetArticles(start, end, order, sort, query, column)

	responseJSON(w, articles)
}

func getArticleCountHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	articleCount := models.CountArticle(query)

	responseJSON(w, articleCount)
}

func getArticleHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	article := models.GetArticle(ID)

	responseJSON(w, article)
}

// admin
func getAdminArticlesHandler(w http.ResponseWriter, r *http.Request) {
	start, _ := strconv.Atoi(r.URL.Query().Get("_start"))
	end, _ := strconv.Atoi(r.URL.Query().Get("_end"))
	order := r.URL.Query().Get("_order")
	sort := r.URL.Query().Get("_sort")
	query := r.URL.Query().Get("q")
	column := r.URL.Query().Get("column")
	articles := models.GetAdminArticles(start, end, order, sort, query, column)
	articleCount := models.CountArticle(query)

	w.Header().Set("X-Total-Count", strconv.Itoa(articleCount))
	responseJSON(w, articles)
}

func getAdminArticleHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	article := models.GetAdminArticle(ID)

	responseJSON(w, article)
}

func createAdminArticleHandler(w http.ResponseWriter, r *http.Request) {
	article := models.CreateArticle(r)

	responseJSON(w, article)
}

func updateAdminArticleHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	article := models.UpdateArticle(r, ID)

	responseJSON(w, article)
}

func deleteAdminArticleHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	models.DeleteArticle(ID)

	// 一覧を返す
	getAdminArticlesHandler(w, r)
}
