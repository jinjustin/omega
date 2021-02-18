package questiongroupcontroller

import (
	//"fmt"
	//"strings"

	"github.com/jinjustin/omega/questiongroup"

	//"encoding/json"
	//"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"database/sql"

	"github.com/jinjustin/omega/database"

	//"omega/authentication"

	"encoding/json"
	"io/ioutil"
	"net/http"
	//"github.com/sqs/goreturns/returns"
)

func groupTestListUpdate(name string, questiongroupID string, questiongroupName string, numQuestion string, score string, courseID string, testID string, uuid string) {
	
	var g questiongroup.QuestionGroup

	if (checkQuestionGroupExist(questiongroupID)){

		g = questiongroup.QuestionGroup{
			Name: name,
			ID: questiongroupID,
			GroupName: questiongroupName,
			NumQuestion: numQuestion,
			Score: score,
			CourseID: courseID,
			TestID: testID,
			UUID: uuid,
		}

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `INSERT INTO questiongroup (name, id, groupname, numquestion, score, courseid, testid, uuid)VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		_, err = db.Exec(sqlStatement, g.Name, g.ID, g.GroupName, g.NumQuestion, g.Score, g.CourseID, g.TestID, g.UUID)
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
		}

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `UPDATE questiongroup SET name = $1, groupname = $2, numquestion = $3, score = $4, testid= $5 uuid=$6 WHERE questiongroupid = $7`

		_, err = db.Exec(sqlStatement, g.Name, g.GroupName, g.NumQuestion, g.Score, g.TestID,g.UUID, g.ID)
		if err != nil {
			panic(err)
		}
	}

}

func checkQuestionGroupExist(questionGroupID string) bool {
	var uuid string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlStatement := `SELECT uuid FROM questiongroup WHERE id=$1;`
	row := db.QueryRow(sqlStatement, questionGroupID)
	err = row.Scan(&uuid)
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

//GroupTestListUpdate is a API that use to get add or update group test list.
var GroupTestListUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	type Item struct {
		ID string
		GroupName string
		NumQuestion string
		Score string
	}

	type Input struct {
		Name string
		Items []Item
	}

	var objmap map[string]Input

	reqBody, _ := ioutil.ReadAll(r.Body)
	//var input Input
	json.Unmarshal(reqBody, &objmap)
	uuids := make([]string, 0, len(objmap))
	for i := range objmap {
        uuids = append(uuids, i)
    }

	var input Input

	courseID := r.Header.Get("CourseID")

	testID := r.Header.Get("TestId")

	for _, uuid := range uuids {
		input = objmap[uuid]
		for _, item := range input.Items{
			groupTestListUpdate(input.Name, item.ID, item.GroupName, item.NumQuestion, item.Score, courseID, testID, uuid)
		}
	}

	w.Write([]byte("success"))
})