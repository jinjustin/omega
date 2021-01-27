package authentication

import (
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

//GetUsername is a function that use to get username from token
func GetUsername(r *http.Request) string{
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	jwtToken := authHeader[1]
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		panic(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return name
}

//GetUserRole is a function that use to get role of user from token
func GetUserRole(r *http.Request) string{
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	jwtToken := authHeader[1]
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		panic(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	return role
}

//GetUsers is a API that use to get username by token.
var GetUsers = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	name := GetUsername(r)
	w.Write([]byte(name))
})

//GetRole is a API that use to get role by token.
var GetRole = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	role := GetUserRole(r)
	w.Write([]byte(role))
})