package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"getartist/app/models"

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

func signup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("successfully called signup"))
}

func responseByJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
	return
}

// Token 作成関数
func CreateToken(user models.User) (string, error) {
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

	fmt.Println("-----------------------------")
	fmt.Println("tokenString:", tokenString)

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
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

	password := user.Password
	fmt.Println("password: ", password)

	// 認証キー(Email)のユーザー情報をDBから取得
	if result := models.DbConnection.Where("email = ?", user.Email).Find(&user); result.Error != nil {
		if result.RecordNotFound() {
			error.Message = "ユーザが存在しません。"
			http.Error(w, error.Message, http.StatusBadRequest)
		} else {
			log.Fatal(result.Error)
		}
	}

	// hasedPassword := user.Password
	hashdPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("hasedPassword: ", hashdPassword)

	err = bcrypt.CompareHashAndPassword([]byte(hashdPassword), []byte(password))

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
