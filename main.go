package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// Payload ...
type Payload struct {
	data string
}

// Error ...
type Error struct {
	Message string `json:"message"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler).Methods("GET")
	r.HandleFunc("/protected", TokenVerify(protectedHandler)).Methods("GET")

	log.Println("listen on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func respondWithError(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("main handler invoked")
	//
	payload := Payload{data: "payload"}
	token, err := GenerateToken(payload)
	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte(token))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("protected handler invoked")
}

// TokenVerify func ..
func TokenVerify(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("token verify invoked")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObject Error
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) == 2 {
			authToken := bearerToken[1]
			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Error")
				}
				return []byte("secret"), nil
			})
			if error != nil {
				errorObject.Message = error.Error()
				respondWithError(w, http.StatusUnauthorized, errorObject)
				return

			}
			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				errorObject.Message = error.Error()
				respondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Message = "Invalid Token"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}

	})
}

func GenerateToken(payload Payload) (string, error) {
	var err error
	secret := "secret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "service-manager",
	})
	// spew.Dump(token)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(tokenString)
	return tokenString, nil
}
