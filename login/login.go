package login

import (
	//"crypto/md5"
	//"encoding/hex"
	//"fmt"
	//"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
    "time"
    "github.com/dgrijalva/jwt-go"
	"encoding/json"
	"database/sql"
	"omega/database"
)

// Input is a ...
type Input struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Token is a ...
type Token struct {
	Token string 
}

// Login is a ...
var Login = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var data Input
	data.Username = r.FormValue("username")
	data.Password = r.FormValue("password")

	var role string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT role FROM users WHERE username=$1 and password=$2;`
	rows,err := db.Query(sqlStatement, data.Username,data.Password)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
 
	for rows.Next() {
	err = rows.Scan(&role)
		 if err != nil {
			 panic(err)
		 }
	 }
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	if role != "" {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = data.Username
		claims["role"] = role
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		 json.NewEncoder(w).Encode( map[string]string{
			"token": t,
			"role": role,
		})

	}else{
		w.Write([]byte("Wrong Username or Password"))
	} 
})

// Logout is a ...
var Logout = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Already logout"))
})