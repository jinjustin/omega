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

// Login is a ...
var Login = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var data Input
	data.Username = r.FormValue("username")
	data.Password = r.FormValue("password")

	var username string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT username FROM users WHERE username=$1 and password=$2;`
	rows,err := db.Query(sqlStatement, data.Username,data.Password)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
 
	for rows.Next() {
	 err = rows.Scan(&username)
		 if err != nil {
			 panic(err)
		 }
	 }
	 err = rows.Err()
	 if err != nil {
	 	panic(err)
	 }

	if data.Username == username {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = data.Username
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		 json.NewEncoder(w).Encode( map[string]string{
			"token": t,
		})

		sqlStatement = `UPDATE users SET token=$1 WHERE username=$2 and password=$3;`

		_, err = db.Exec(sqlStatement, t,data.Username,data.Password)
		if err != nil {
			panic(err)
		}

		//fmt.Println(t)

	}else{
		w.Write([]byte("Wrong Username or Password"))
	} 
})

// Logout is a ...
var Logout = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Already logout"))
})

func checkToken (username string, token string) bool{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var userid string

	sqlStatement := `SELECT userid FROM users WHERE username=$1 and token=$2;`
	row := db.QueryRow(sqlStatement, username,token)
	err = row.Scan(&userid)
	switch err {
	case sql.ErrNoRows: return false
	case nil: return true
	default: panic(err)
	}

}