package answercontroller

import (
	//"fmt"
	//"github.com/jinjustin/omega/answer"

	"encoding/json"
	//"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"database/sql"

	"github.com/jinjustin/omega/answer"
	"github.com/jinjustin/omega/database"

	//"github.com/jinjustin/omega/authentication"

	"net/http"
	//"encoding/json"
	"io/ioutil"
)

func submitAnswer(testID string, studentID string, studentAnswer []answer.Info) error {

	//PlaceHolder
	totalScore := 0

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	b, err := json.Marshal(studentAnswer)
	if err != nil {
		panic(err)
	}

	sqlStatement := `INSERT INTO answer (testid, studentid, studentanswer, totalscore, completepercent)VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(sqlStatement, testID, studentID, b, totalScore, "0")
	if err != nil {
		return err
	}

	return nil
}

func getStudentAnswer(testID string, studentID string) ([]byte, error){

	var b []byte

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStatement := `SELECT studentanswer FROM answer WHERE testid=$1 and studentid=$2`
	rows, err := db.Query(sqlStatement, testID, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&b)
		if err != nil {
			return nil ,err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return b, err
}

func autoScoring (studentAnswer []answer.Info, testID string) error{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	var correctcheck string
	var groupID string
	var score string

	for num, a := range studentAnswer{

		if a.QuestionType == "choice"{
			sqlStatement := `SELECT correctcheck FROM choice WHERE choiceid=$1 and questionid=$2`
			rows, err := db.Query(sqlStatement, a.Answer, a.QuestionID)
			if err != nil {
				return nil
			}
			defer rows.Close()

			for rows.Next() {
				err = rows.Scan(&correctcheck)
				if err != nil {
					return nil
				}
			}
			err = rows.Err()
			if err != nil {
				return nil
			}

			sqlStatement = `SELECT groupid FROM question WHERE testid=$1 and questionid=$2`
			rows, err = db.Query(sqlStatement, testID, a.QuestionID)
			if err != nil {
				return nil
			}
			defer rows.Close()

			for rows.Next() {
				err = rows.Scan(&groupID)
				if err != nil {
					return nil
				}
			}
			err = rows.Err()
			if err != nil {
				return nil
			}

			sqlStatement = `SELECT score FROM questiongroup WHERE id=$1`
			rows, err = db.Query(sqlStatement, groupID)
			if err != nil {
				return nil
			}
			defer rows.Close()

			for rows.Next() {
				err = rows.Scan(&score)
				if err != nil {
					return nil
				}
			}
			err = rows.Err()
			if err != nil {
				return nil
			}

			if correctcheck == "true"{
				studentAnswer[num].Score = score
			}
		}
	}

	return nil
}

//API

//SubmitAnswer is a function that use to store student answer to database.
var SubmitAnswer = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var studentAnswer []answer.Info

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil{
		http.Error(w, "can't read body", http.StatusBadRequest)
            return
	}

	err = json.Unmarshal(reqBody,&studentAnswer)
	if err != nil{
		http.Error(w, "Can't convert JSON into map", http.StatusBadRequest)
            return
	}

	testID := r.Header.Get("TestId")

	studentID := r.Header.Get("StudentID")

	err = submitAnswer(testID, studentID, studentAnswer)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})

//GetAnswer is a function that use to get student answer from database.
var GetAnswer = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	testID := r.Header.Get("TestId")

	studentID := r.Header.Get("StudentID")

	b, err := getStudentAnswer(testID, studentID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
})