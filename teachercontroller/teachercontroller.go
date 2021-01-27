package teachercontroller

import(
	"omega/teacher"
	"omega/database"
	"omega/authentication"
	"database/sql"

	"net/http"
	"encoding/json"
	"io/ioutil"
)

func getTeacherInfo(username string) []byte{

	var userID string
	var firstname string
	var surname string
	var email string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM users WHERE username=$1;`
	rows, err := db.Query(sqlStatement, username)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	sqlStatement = `SELECT firstname,surname,email FROM teacher WHERE userid=$1;`
	rows, err = db.Query(sqlStatement, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&firstname,&surname,&email)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	t := teacher.Teacher{
		UserID: "",
		Firstname: firstname,
		Surname: surname,
		Email: email,
	}

	return t.GetTeacherDetail()
}

func editTeacherInfo(firstname string,surname string,email string,username string)[]byte{
	
	var userID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM users WHERE username=$1;`
	rows, err := db.Query(sqlStatement, username)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	sqlStatement = `UPDATE teacher SET firstname = $1,surname = $2,email = $3 WHERE userid = $4`

	_, err = db.Exec(sqlStatement, firstname, surname, email, userID)
	if err != nil {
		panic(err)
	}

	sqlStatement = `SELECT firstname,surname,email FROM teacher WHERE userid=$1;`
	rows, err = db.Query(sqlStatement, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&firstname,&surname,&email)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	t := teacher.Teacher{
		UserID: "",
		Firstname: firstname,
		Surname: surname,
		Email: email,
	}

	return t.GetTeacherDetail()
}

//API

//GetTeacherInfo is a API that use for get teacher information.
var GetTeacherInfo = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	username := authentication.GetUsername(r)
	w.Write(getTeacherInfo(username))
})

//EditTeacherInfo is a API that use for edit teacher information.
var EditTeacherInfo = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct{
		Firstname string
		Surname string
		Email string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	username := authentication.GetUsername(r)
	w.Write(editTeacherInfo(input.Firstname,input.Surname,input.Email,username))
})


