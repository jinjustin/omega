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

/*
func createQuestionGroup(name string,courseID string, questionType string) []byte {

	var g questiongroup.QuestionGroup

	questionGroupID := generateID()

	g = questiongroup.QuestionGroup{
		QuestionGroupID: questionGroupID,
		Name: name,
		CourseID: courseID,
		Type: questionType,
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO questiongroup (questiongroupid, name, courseid, type)VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, g.QuestionGroupID,g.Name,g.CourseID,g.Type)
	if err != nil {
		panic(err)
	}

	return g.GetQuestionGroupDetail()
}

func getQuestionGroupList(courseID string) []questiongroup.QuestionGroup {
	var questionGroupList []questiongroup.QuestionGroup
	var g questiongroup.QuestionGroup

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT questiongroupid, name, type FROM questiongroup WHERE courseid=$1;`
	rows, err := db.Query(sqlStatement, courseID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&g.QuestionGroupID, &g.Name, &g.Type)
		if err != nil {
			panic(err)
		}
		g.CourseID = courseID
		questionGroupList = append(questionGroupList, g)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return questionGroupList
}

func getQuestionGroupInfo(questionGroupID string) []byte{
	var g questiongroup.QuestionGroup

	g.QuestionGroupID = questionGroupID

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT name, courseid, type FROM questiongroup WHERE questiongroupid=$1;`
	rows, err := db.Query(sqlStatement, g.QuestionGroupID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&g.Name, &g.CourseID, &g.Type)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return g.GetQuestionGroupDetail()
}

func editQuestionGroupName(questionGroupID string, name string) string{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `UPDATE questiongroup SET name=$1 WHERE questiongroupid=$2`

	_, err = db.Exec(sqlStatement, name, questionGroupID)
	if err != nil {
		panic(err)
	}

	return "success"
}

func deleteQuestionGroup(questionGroupID string) []byte{

	g := questiongroup.QuestionGroup{
		QuestionGroupID: "Can't find",
		Name: "",
		CourseID: "",
		Type: "",
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT name, courseid, type FROM questiongroup WHERE questiongroupid=$1;`
	rows, err := db.Query(sqlStatement, questionGroupID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&g.Name,&g.CourseID,&g.Type)
		if err != nil {
			panic(err)
		}
		g.QuestionGroupID = questionGroupID
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM questiongroup WHERE questiongroupid=$1;`
	_, err = db.Exec(sqlStatement, questionGroupID)
	if err != nil {
		panic(err)
	}

	sqlStatement = `DELETE FROM question WHERE questiongroupid=$1;`
	_, err = db.Exec(sqlStatement, questionGroupID)
	if err != nil {
		panic(err)
	}

	return g.GetQuestionGroupDetail()
}

func checkQuestionGroupName(courseID string, name string) bool {
	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var questionGroupID string

	sqlStatement := `SELECT questiongroupid FROM questiongroup WHERE courseid=$1 and name=$2;`
	row := db.QueryRow(sqlStatement, courseID, name)
	err = row.Scan(&questionGroupID)
	switch err {
	case sql.ErrNoRows:
		return true
	case nil:
		return false
	default:
		panic(err)
	}
}*/

func groupTestListUpdate(name string, questiongroupID string, questiongroupName string, numQuestion string, score string, courseID string, testID string, key string) []byte{
	
	var g questiongroup.QuestionGroup

	g = questiongroup.QuestionGroup{
		Name: "Error",
		QuestionGroupID: "Error",
		QuestionGroupName: "Error",
		NumQuestion: "Error",
		Score: "Error",
		CourseID: "Error",
		TestID: "Error",
		Key: "Error",
	}

	if (checkQuestionGroupExist(questiongroupID)){

		g = questiongroup.QuestionGroup{
			Name: name,
			QuestionGroupID: questiongroupID,
			QuestionGroupName: questiongroupName,
			NumQuestion: numQuestion,
			Score: score,
			CourseID: courseID,
			TestID: testID,
			Key: key,
		}

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `INSERT INTO questiongroup (name, questiongroupid, questiongroupname, numquestion, score, courseid, testid, key)VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		_, err = db.Exec(sqlStatement, g.Name, g.QuestionGroupID, g.QuestionGroupName, g.NumQuestion, g.Score, g.CourseID, g.TestID, g.Key)
		if err != nil {
			panic(err)
		}
	}else{

		g = questiongroup.QuestionGroup{
			Name: name,
			QuestionGroupID: questiongroupID,
			QuestionGroupName: questiongroupName,
			NumQuestion: numQuestion,
			Score: score,
			CourseID: courseID,
			TestID: testID,
			Key: key,
		}

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `UPDATE questiongroup SET name = $1, questiongroupname = $2, numquestion = $3, score = $4 WHERE questiongroupid = $5`

		_, err = db.Exec(sqlStatement, g.Name, g.QuestionGroupName, g.NumQuestion, g.Score, g.QuestionGroupID)
		if err != nil {
			panic(err)
		}
	}

	return g.GetQuestionGroupDetail()
}

func checkQuestionGroupExist(questionGroupID string) bool {
	var key string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlStatement := `SELECT key FROM questiongroup WHERE questiongroupid=$1;`
	row := db.QueryRow(sqlStatement, questionGroupID)
	err = row.Scan(&key)
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

//CreateQuestionGroup is a API that use to create questiongroup in the course.
/*var CreateQuestionGroup = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Name string
		CourseID string
		Type string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(createQuestionGroup(input.Name,input.CourseID,input.Type))
})

//GetQuestionGroupList is a API that use to get question group list in the course.
var GetQuestionGroupList = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		CourseID string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(getQuestionGroupList(input.CourseID))
})

//GetQuestionGroupInfo is a API that use to get information of question group by using questiongroupID.
var GetQuestionGroupInfo = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		QuestionGroupID string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(getQuestionGroupInfo(input.QuestionGroupID))
})

//EditQuestionGroupName is a API that use to edit question group name.
var EditQuestionGroupName = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		QuestionGroupID string
		Name string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(editQuestionGroupName(input.QuestionGroupID,input.Name))
})

//DeleteQuestionGroup is a API that use to delete question group by using questiongroupid
var DeleteQuestionGroup = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		QuestionGroupID string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(deleteQuestionGroup(input.QuestionGroupID))
})*/

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
		Items Item
	}

	var objmap map[string]Input

	reqBody, _ := ioutil.ReadAll(r.Body)
	//var input Input
	json.Unmarshal(reqBody, &objmap)
	keys := make([]string, 0, len(objmap))
	for k := range objmap {
        keys = append(keys, k)
    }

	key := keys[0]

	var input Input

	input = objmap[key]

	courseID := r.Header.Get("CourseID")

	testID := r.Header.Get("TestId")

	w.Write(groupTestListUpdate(input.Name, input.Items.ID, input.Items.GroupName, input.Items.NumQuestion, input.Items.Score, courseID, testID, key))
})