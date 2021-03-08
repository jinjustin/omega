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

	"github.com/iancoleman/orderedmap"
)


func submitAnswer(testID string, studentID string, studentAnswer map[string]answer.Info, keys []string) error {

	o := orderedmap.New()

	for _, key := range keys{
		o.Set(key,studentAnswer[key])
	}

	b,err := o.MarshalJSON()
	if err != nil{
		return nil
	}

	//PlaceHolder
	totalScore := 0

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `INSERT INTO answer (testid, studentid, studentanswer, totalscore, completepercent)VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(sqlStatement, testID, studentID, b, totalScore, 0)
	if err != nil {
		return err
	}

	return nil
}

func autoScoring (studentAnswer map[string]answer.Info, keys []string, testID string) error{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	var correctness string

	for _, key := range keys{
		a, _ := studentAnswer[key]
		if a.QuestionType == "choice"{
			sqlStatement := `SELECT correctcheck FROM choice WHERE choiceid=$1 and questionid=$2`
			rows, err := db.Query(sqlStatement, a.Answer, key)
			if err != nil {
				return nil
			}
			defer rows.Close()

			for rows.Next() {
				err = rows.Scan(&correctness)
				if err != nil {
					return nil
				}
			}
			err = rows.Err()
			if err != nil {
				return nil
			}
			if correctness == "true"{
				
			}
		}
	}

	return nil
}

//API

//SubmitAnswer is a function that use to store student answer to database.
var SubmitAnswer = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	o := orderedmap.New()

	var studentAnswer map[string]answer.Info

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil{
		http.Error(w, "can't read body", http.StatusBadRequest)
            return
	}

	err = o.UnmarshalJSON(reqBody)
	if err != nil{
		http.Error(w, "Can't convert JSON into map", http.StatusBadRequest)
            return
	}

	keys := o.Keys()

	err = json.Unmarshal(reqBody,&studentAnswer)
	if err != nil{
		http.Error(w, "Can't convert JSON into map", http.StatusBadRequest)
            return
	}

	testID := r.Header.Get("testID")

	studentID := r.Header.Get("studentID")

	err = submitAnswer(testID, studentID, studentAnswer, keys)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})