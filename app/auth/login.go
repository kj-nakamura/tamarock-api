package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"api/app/models"

	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

func responseByJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
	return
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.AdminUser
	var error Error

	// r.body に何が帰ってくるか確認
	fmt.Println(r.Body)

	// https://golang.org/pkg/encoding/json/#NewDecoder
	// json.NewDecoder(r.Body).Decode(&user)
	dec := json.NewDecoder(r.Body)
	for err := dec.Decode(&user); err != nil && err != io.EOF; {
		log.Println("ERROR: " + err.Error())
	}

	if user.Email == "" {
		error.Message = "Email は必須です。"
		http.Error(w, error.Message, http.StatusBadRequest)
		return
	}

	if user.Password == "" {
		error.Message = "パスワードは必須です。"
		http.Error(w, error.Message, http.StatusBadRequest)
		return
	}

	// user に何が格納されているのか
	fmt.Println(user)

	// dump も出せる
	fmt.Println("---------------------")
	// spew.Dump(user)

	// パスワードのハッシュを生成
	// https://godoc.org/golang.org/x/crypto/bcrypt#GenerateFromPassword
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("パスワード: ", user.Password)
	fmt.Println("ハッシュ化されたパスワード", hash)

	user.Password = string(hash)
	fmt.Println("コンバート後のパスワード: ", user.Password)

	result := models.DbConnection.Create(&user)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	if err != nil {
		error.Message = "サーバーエラー"
		http.Error(w, error.Message, http.StatusBadRequest)
		return
	}

	// DB に登録できたらパスワードをからにしておく
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")

	// JSON 形式で結果を返却
	responseByJSON(w, user)
}

// Token 作成関数
func CreateToken(user models.AdminUser) (string, error) {
	var err error

	// 鍵となる文字列(多分なんでもいい)
	secret := "secret"

	// Token を作成
	// jwt -> JSON Web Token - JSON をセキュアにやり取りするための仕様
	// jwtの構造 -> {Base64 encoded Header}.{Base64 encoded Payload}.{Signature}
	// HS254 -> 証明生成用(https://ja.wikipedia.org/wiki/JSON_Web_Token)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "nakamura", // JWT の発行者が入る(文字列(__init__)は任意)
	})

	//Dumpを吐く
	spew.Dump(token)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.AdminUser
	var error Error
	var jwt JWT

	dec := json.NewDecoder(r.Body)
	for err := dec.Decode(&user); err != nil && err != io.EOF; {
		log.Println("ERROR: " + err.Error())
		return
	}

	if user.Email == "" {
		error.Message = "Email は必須です。"
		http.Error(w, error.Message, http.StatusBadRequest)
		return
	}

	if user.Password == "" {
		error.Message = "パスワードは、必須です。"
		http.Error(w, error.Message, http.StatusBadRequest)
	}

	// リクエストのパスワード
	password := user.Password
	hashdPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Fatal(err)
	}

	// 認証キー(Email)のユーザー情報をDBから取得
	if result := models.DbConnection.Where("email = ?", user.Email).Find(&user); result.Error != nil {
		if result.RecordNotFound() {
			error.Message = "ユーザが存在しません。"
			http.Error(w, error.Message, http.StatusBadRequest)
		} else {
			log.Fatal(result.Error)
		}
	}

	// リクエストのパスワードとDBから取得したパスワードを比較
	err = bcrypt.CompareHashAndPassword(hashdPassword, []byte(password))

	if err != nil {
		error.Message = "無効なパスワードです。"
		http.Error(w, error.Message, http.StatusBadRequest)
		return
	}

	token, err := CreateToken(user)

	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	jwt.Token = token

	responseByJSON(w, jwt)
}
