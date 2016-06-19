package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
	"time"
)

func ValidateRequest(pass http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			pass(w, r)
		} else {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}
}

func TokenAuth(pass http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)

		// Check to see if the Authorization Header has two parts.  A "Bearer" followed by something
		if len(auth) != 2 || auth[0] != "Bearer" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		payload := auth[1]

		// Attempt to Validate what was sent in the second array element of the Authorization header
		if !ValidateToken(payload) {
			http.Error(w, "Authorization failed", http.StatusUnauthorized)
			return
		}

	}
}

func ValidateToken(tkn string) bool {

	tokenToBeValidated := tkn
	token, err := jwt.Parse(tokenToBeValidated, func(token *jwt.Token) (interface{}, error) {

		// Check What Algorithm Signed the Token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Get the key used to sign the token in generateToken
		return GetSigningKey(), nil
	})

	if err != nil {
		fmt.Println("ValidateToken::Error on parsing::", err)
	}
	if err == nil && token.Valid {
		return true
	} else {
		return false
	}
}

func GetSigningKey() []byte {
	// Read File, Database Call, Environment Variable, whereever we want to store the key
	mySigningKey := []byte(os.Getenv("SuperSecretKey"))

	if mySigningKey == nil {
		mySigningKey = []byte("SecretLike")
	}

	return mySigningKey
}

func generateToken(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	payload := r.Form["payload"]

	if payload == nil {
		http.Error(w, "Bad syntax", http.StatusBadRequest)
		return
	}

	mySigningKey := GetSigningKey()

	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["payload"] = payload
	token.Claims["nbf"] = time.Now()
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return
	}
	fmt.Fprintf(w, "%s\n", tokenString)
}
