package questioncontroller

import (
	"fmt"

	//"encoding/json"
	"crypto/rand"
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

func addNewQuestion(groupID string, testID string, question string, questionID string, questionType string, data string){
	var q question.Question

	if (checkQuestionInGroup(questionID, groupID)){

		q = question.Question{
			TestID: testID,
			GroupID: groupID,

		}

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()


		if(checkQuestionGroupExist(questiongroupID)){
			sqlStatement := `INSERT INTO questiongroup (name, id, groupname, numquestion, score, courseid, testid, uuid, headerorder, grouporder)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
			_, err = db.Exec(sqlStatement, g.Name, g.ID, g.GroupName, g.NumQuestion, g.Score, g.CourseID, "", g.UUID, g.HeaderOrder, g.GroupOrder)
			if err != nil {
				panic(err)
			}
		}

		sqlStatement := `INSERT INTO questiongroup (name, id, groupname, numquestion, score, courseid, testid, uuid, headerorder, grouporder)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
		_, err = db.Exec(sqlStatement, g.Name, g.ID, g.GroupName, g.NumQuestion, g.Score, g.CourseID, g.TestID, g.UUID, g.HeaderOrder, g.GroupOrder)
		if err != nil {
			panic(err)
		}

	}else{

		g = questiongroup.QuestionGroup{
			Name: name,
			ID: questiongroupID,
			GroupName: questiongroupName,
			NumQuestion: numQuestion,
			Score: score,
			CourseID: courseID,
			TestID: testID,
			UUID: uuid,
			HeaderOrder: headerOrder,
			GroupOrder: groupOrder,
		}

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()
	
		sqlStatement := `UPDATE questiongroup SET name=$1, groupname=$2, numquestion=$3, score=$4, uuid=$5, headerorder=$6, grouporder=$7 WHERE id=$8`
	
		_, err = db.Exec(sqlStatement, g.Name, g.GroupName, g.NumQuestion, g.Score, g.UUID, g.HeaderOrder, g.GroupOrder, g.ID)
		if err != nil {
			panic(err)
		}
	}
}

func checkQuestionInGroup(questionID string, groupID string) bool {
	var question string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlStatement := `SELECT question FROM question WHERE questionid=$1 and groupid=$2;`
	row := db.QueryRow(sqlStatement, questionID, groupID)
	err = row.Scan(&question)
	switch err {
	case sql.ErrNoRows:
		return true
	case nil:
		return false
	default:
		panic(err)
	}
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
	w.Write([]byte("Place Holder"))
})