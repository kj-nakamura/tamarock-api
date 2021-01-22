package controllers

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"syscall"

	"api/app/auth"
	"api/app/models"
	"api/config"

	"github.com/gorilla/mux"
	"golang.org/x/sys/unix"
)

type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type HealthCheck struct {
	Status int
	Result string
}

func APIError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JSONError{Error: errMessage, Code: code})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonError)
}

// responseJSON JSON形式に変換する
func responseJSON(w http.ResponseWriter, value interface{}) {
	js, err := json.Marshal(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// getID URLのIDを取得する
func getID(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	return strconv.Atoi(vars["id"])
}

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

// healthCheckHandler is ALBによるヘルスチェック用
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ping := HealthCheck{http.StatusOK, "ok"}

	res, err := json.Marshal(ping)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func listenCtrl(network string, address string, c syscall.RawConn) error {
	var err error
	c.Control(func(s uintptr) {
		err = unix.SetsockoptInt(int(s), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1) // portをbindできる設定
		if err != nil {
			return
		}
	})
	return err
}

func StartWebServer() error {
	r := mux.NewRouter()

	// image
	var dir string
	flag.StringVar(&dir, "dir", "./static", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	// web
	// health check
	r.HandleFunc("/api/health-check", healthCheckHandler).Methods("GET")
	// artist
	r.HandleFunc("/api/search", searchArtistHandler).Methods("GET")
	r.HandleFunc("/api/artist/infos", getArtistInfosHandler).Methods("GET")
	r.HandleFunc("/api/artist/infos/count", getArtistInfosCountHandler).Methods("GET")
	r.HandleFunc("/api/artist/info/{id}", getArtistInfoHandler).Methods("GET")
	r.HandleFunc("/api/artist/{id}", getArtistHandler).Methods("GET")

	// article
	r.HandleFunc("/api/articles", getArticlesHandler).Methods("GET")
	r.HandleFunc("/api/articles/count", getArticleCountHandler).Methods("GET")
	r.HandleFunc("/api/articles/{id}", getArticleHandler).Methods("GET")

	// admin
	// artist
	r.HandleFunc("/api/admin/artists", auth.TokenVerifyMiddleWare(getAdminArtistsHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/artists/{id}", auth.TokenVerifyMiddleWare(getAdminArtistHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/artists", auth.TokenVerifyMiddleWare(createArtistHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/admin/artists/{id}", auth.TokenVerifyMiddleWare(updateArtistHandler)).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/admin/artists/{id}", auth.TokenVerifyMiddleWare(deleteArtistHandler)).Methods("DELETE", "OPTIONS")

	// article
	r.HandleFunc("/api/admin/articles", auth.TokenVerifyMiddleWare(getAdminArticlesHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/articles/{id}", auth.TokenVerifyMiddleWare(getAdminArticleHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/articles", auth.TokenVerifyMiddleWare(createAdminArticleHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/admin/articles/{id}", auth.TokenVerifyMiddleWare(updateAdminArticleHandler)).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/admin/articles/{id}", auth.TokenVerifyMiddleWare(deleteAdminArticleHandler)).Methods("DELETE", "OPTIONS")

	// category
	r.HandleFunc("/api/admin/categories", auth.TokenVerifyMiddleWare(getAdminCategoriesHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/categories/{id}", auth.TokenVerifyMiddleWare(getAdminCategoryHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/categories", auth.TokenVerifyMiddleWare(createAdminCategoryHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/admin/categories/{id}", auth.TokenVerifyMiddleWare(updateAdminCategoryHandler)).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/admin/categories/{id}", auth.TokenVerifyMiddleWare(deleteAdminCategoryHandler)).Methods("DELETE", "OPTIONS")

	// auth
	r.HandleFunc("/api/admin/login", auth.Login).Methods("POST", "OPTIONS")
	// r.HandleFunc("/api/admin/signup", auth.Signup).Methods("POST")

	r.HandleFunc("/health-check/", healthCheckHandler)
	http.Handle("/", r)

	lc := net.ListenConfig{
		Control: listenCtrl, //portのbindを許可する設定を入れる
	}

	listener, err := lc.Listen(context.Background(), "tcp4", fmt.Sprintf(":%d", config.Config.Port))
	if err != nil {
		panic(err)
	}

	return http.Serve(listener, nil)
}
