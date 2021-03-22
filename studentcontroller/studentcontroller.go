package studentcontroller

import(
	"github.com/jinjustin/omega/student"
	"github.com/jinjustin/omega/database"
	"github.com/jinjustin/omega/authentication"
	"database/sql"

	"net/http"
	"encoding/json"
	"io/ioutil"
)

func getStudentInfo(username string) []byte{

	var userID string
	var firstname string
	var surname string
	var studentID string
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

	sqlStatement = `SELECT studentid, firstname, surname, email FROM teacher WHERE userid=$1;`
	rows, err = db.Query(sqlStatement, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&studentID, &firstname, &surname, &email)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	s := student.Student{
		UserID: "",
		StudentID: studentID,
		Firstname: firstname,
		Surname: surname,
		Email: email,
	}

	return s.GetStudentDetail()
}

func editStudentInfo(firstname string, surname string, username string) error{
	
	var userID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM users WHERE username=$1;`
	rows, err := db.Query(sqlStatement, username)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	sqlStatement = `UPDATE student SET firstname = $1, surname = $2 WHERE userid = $4`

	_, err = db.Exec(sqlStatement, firstname, surname, userID)
	if err != nil {
		return err
	}

	return nil
}

//API

//GetStudentInfo is a API that use for get student information by send authentication token.
var GetStudentInfo = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	username := authentication.GetUsername(r)
	w.Write(getStudentInfo(username))
})

//EditTeacherInfo is a API that use for edit teacher information.
var EditTeacherInfo = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct{
		Firstname string
		Surname string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	username := authentication.GetUsername(r)
	err := editStudentInfo(input.Firstname,input.Surname,username)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})