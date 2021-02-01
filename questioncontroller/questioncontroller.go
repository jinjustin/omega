package question

import (
	"fmt"

	//"encoding/json"
	"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"database/sql"
	"omega/database"
	"omega/question"
	//"omega/authentication"

	"encoding/json"
	"io/ioutil"
	"net/http"
	//"github.com/sqs/goreturns/returns"
)

func createQuestion(questionGroupID string, inputQuestion string, questionType string) []byte {

	var q question.Question

	if inputQuestion == "" {
		inputQuestion = "Insert your question here."
	}

	if checkQuestionType(questionGroupID, questionType) == false{
		q = question.Question{
			QuestionID : "",
			QuestionGroupID: "",
			Question: "Question Type doesn't match with group.",
		}

		return q.GetQuestionDetail()
	}

	questionID := generateID()

	q = question.Question{
		QuestionID : questionID,
		QuestionGroupID: questionGroupID,
		Question: inputQuestion,
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO question (questionid, questiongroupid, question)VALUES ($1, $2, $3)`

	_, err = db.Exec(sqlStatement, q.QuestionID, q.QuestionGroupID, q.Question)
	if err != nil {
		panic(err)
	}

	return q.GetQuestionDetail()
}

func generateID() string {
	n := 3
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return s
}

func checkQuestionType(questionGroupID string, questionType string) bool{

	var typeOfQuestionGroup string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT type FROM questiongroup WHERE questiongroupid=$1;`
	rows, err := db.Query(sqlStatement, questionGroupID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&typeOfQuestionGroup)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	if questionType == typeOfQuestionGroup {
		return true
	}

	return false
}

//API

//CreateQuestion is a API that use to create questiongroup in the course.
var CreateQuestion = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		QuestionID string
		Question string
		Type string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(createQuestion(input.QuestionID,input.Question,input.Type))
})