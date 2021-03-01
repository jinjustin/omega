package questioncontroller

import (
	"fmt"
	//"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"database/sql"
	"github.com/jinjustin/omega/database"
	"github.com/jinjustin/omega/question"
	
	//"omega/authentication"

	"encoding/json"
	"io/ioutil"
	"net/http"
	//"github.com/sqs/goreturns/returns"
)

func addNewQuestion(groupID string, testID string, questionName string, questionID string, questionType string, data string) error{
	var q question.Question

	q = question.Question{
		TestID: testID,
		GroupID: groupID,
		QuestionName: questionName,
		QuestionID: questionID,
		QuestionType: questionType,
		Data: data,
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	if testID ==""{

		checkExist := checkQuestionExist(questionID)
		fmt.Println(checkExist)

		if checkExist == sql.ErrNoRows{
			sqlStatement := `INSERT INTO question (testid, groupid, questionname, questiontype, data)VALUES ('', $1, $2, $3, $4)`
			_, err = db.Exec(sqlStatement, q.GroupID, q.QuestionName, q.QuestionType, q.Data)
			if err != nil {
				return err
			}
		}else if checkExist == nil{
			sqlStatement := `UPDATE question SET questionname=$1, questiontype=$2, data=$3 WHERE questionid=$4`
	
			_, err = db.Exec(sqlStatement, q.QuestionName, q.QuestionType, q.Data, q.QuestionID)
			if err != nil {
				return err
			}
		}else{
			return err
		}
	}else{

		checkInTest := checkQuestionInTest(questionID, testID)

		checkExist := checkQuestionExist(questionID)

		fmt.Println(checkInTest)
		fmt.Println(checkExist)

		if checkInTest == sql.ErrNoRows{
			if checkExist == sql.ErrNoRows{
				sqlStatement := `INSERT INTO question (testid, groupid, questionname, questiontype, data)VALUES ('', $1, $2, $3, $4)`
				_, err = db.Exec(sqlStatement, q.GroupID, q.QuestionName, q.QuestionType, q.Data)
				if err != nil {
					return err
				}
			}else if checkExist == nil{
				sqlStatement := `INSERT INTO question (testid, groupid, questionname, questiontype, data)VALUES ($1, $2, $3, $4, $5)`
				_, err = db.Exec(sqlStatement, q.TestID, q.GroupID, q.QuestionName, q.QuestionType, q.Data)
				if err != nil {
					return err
				}
			}else{
				return err
			}
		}else if checkInTest == nil {
			sqlStatement := `UPDATE question SET questionname=$1, questiontype=$2, data=$3 WHERE questionid=$4`
			_, err = db.Exec(sqlStatement, q.QuestionName, q.QuestionType, q.Data, q.QuestionID)
			if err != nil {
				return err
			}
		}else{
			return err
		}
	}
	return err
}

func getQuestion(groupID string, testID string, questionID string) (question.Question, error){

	q := question.Question{
		TestID: testID,
		GroupID: groupID,
		QuestionID: questionID,
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return q, err
	}
	defer db.Close()

	sqlStatement := `SELECT questionname, questiontype, data FROM question WHERE groupid=$1 and testid=$2 and questionid=$3`
	rows, err := db.Query(sqlStatement, groupID, testID, questionID)
	if err != nil {
		return q, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&q.QuestionName, &q.QuestionType, &q.Data)
		if err != nil {
			return q, err
		}
	}
	err = rows.Err()
	if err != nil {
		return q, err
	}

	return q, err
}

func deleteQuestion(testID string, questionID string, groupID string) error {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	if testID==""{
		sqlStatement := `DELETE from question WHERE questionid=$1 and groupid=$2 ;`
		_, err = db.Exec(sqlStatement, questionID, groupID)
		if err != nil {
			return err
		}
	}else{
		sqlStatement := `DELETE from question WHERE questionid=$1 and groupid=$2 and testid=$3;`
		_, err = db.Exec(sqlStatement, questionID, groupID, testID)
		if err != nil {
			return err
		}
	}

	return err
}

func getAllQuestionInGroup(testID string, groupID string) ([]question.AllQuestionInGroup, error) {

	var allQuestionInGroup []question.AllQuestionInGroup

	var a question.AllQuestionInGroup

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if testID == ""{
		sqlStatement := `SELECT questionid, questionname FROM question WHERE testid='' and groupid=$1`
		rows, err := db.Query(sqlStatement, groupID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	
		for rows.Next() {
			err = rows.Scan(&a.QuestionID, &a.QuestionName)
			if err != nil {
				return nil, err
			}

			allQuestionInGroup = append(allQuestionInGroup, a)
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}
	}else{
		sqlStatement := `SELECT questionid, questionname FROM question WHERE testid=$1 and groupid=$2`
		rows, err := db.Query(sqlStatement, testID, groupID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	
		for rows.Next() {
			err = rows.Scan(&a.QuestionID, &a.QuestionName)
			if err != nil {
				return nil, err
			}

			allQuestionInGroup = append(allQuestionInGroup, a)
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}
	}

	return allQuestionInGroup, err
}

func checkQuestionInGroup(questionID string, groupID string) error {
	var questionName string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()
	sqlStatement := `SELECT questionname FROM question WHERE questionid=$1 and groupid=$2;`
	row := db.QueryRow(sqlStatement, questionID, groupID)
	err = row.Scan(&questionName)
	return err
}

func checkQuestionInTest(questionID string, testID string) error {
	var questionName string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()
	sqlStatement := `SELECT questionname FROM question WHERE questionid=$1 and testid=$2;`
	row := db.QueryRow(sqlStatement, questionID, testID)
	err = row.Scan(&questionName)
	return err
}

func checkQuestionExist(questionID string) error {
	var questionName string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()
	sqlStatement := `SELECT questionname FROM question WHERE questionid=$1;`
	row := db.QueryRow(sqlStatement, questionID)
	err = row.Scan(&questionName)
	return err
}

//API

//AddNewQuestion is a API that use to add question to question group.
var AddNewQuestion = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Data string
	}

	groupID := r.Header.Get("GroupId")
	testID := r.Header.Get("TestID")
	questionName := r.Header.Get("Question")
	questionID := r.Header.Get("QuestionID")
	questionType := r.Header.Get("Type")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)

	err := addNewQuestion(groupID,testID,questionName,questionID,questionType,input.Data)

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		fmt.Println(err)
	}else{
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 - OK"))
	}
})

//GetQuestion is a API that use to get question information in question group.
var GetQuestion = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	groupID := r.Header.Get("GroupId")
	testID := r.Header.Get("TestID")
	questionID := r.Header.Get("QuestionID")

	q, err := getQuestion(groupID,testID,questionID)

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		fmt.Println(err)
	}else{
		w.WriteHeader(http.StatusOK)
		w.Write(q.GetQuestionDetail())
	}
})

//DeleteQuestion is a API that use to delete question in question group (in test or in all question group).
var DeleteQuestion = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	groupID := r.Header.Get("GroupId")
	testID := r.Header.Get("TestID")
	questionID := r.Header.Get("QuestionID")

	err := deleteQuestion(groupID,testID,questionID)

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		fmt.Println(err)
	}else{
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 - OK"))
	}
})

//GetAllQuestionInGroup is a API that use to get id and name of all question in question group.
var GetAllQuestionInGroup = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	testID := r.Header.Get("TestID")
	groupID := r.Header.Get("GroupId")

	allQuestionInGroup, err := getAllQuestionInGroup(testID,groupID)

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		fmt.Println(err)
	}else{
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(allQuestionInGroup)
	}
})