package questiongroupcontroller

import (
	"fmt"
	//"strings"

	//"github.com/golang/protobuf/descriptor"
	//"github.com/jinjustin/omega/question"
	"github.com/jinjustin/omega/questiongroup"
	//"github.com/jinjustin/omega/test"

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

//GrouptestList is struct that use to return grouptestlist.
type GrouptestList struct{
	ID string
	GroupName string
}

func groupTestListUpdate(name string, questiongroupID string, questiongroupName string, numQuestion string, score string, courseID string, testID string, uuid string, headerOrder int, groupOrder int) {
	
	var g questiongroup.QuestionGroup

	if (checkQuestionGroupInTest(questiongroupID, testID)){

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

func testbankUpdate(name string, questiongroupID string, questiongroupName string, numQuestion string, score string, courseID string, uuid string) {
	
	var g questiongroup.QuestionGroup

	if (checkQuestionGroupInTestbank(questiongroupID)){

		g = questiongroup.QuestionGroup{
			Name: name,
			ID: questiongroupID,
			GroupName: questiongroupName,
			NumQuestion: numQuestion,
			Score: score,
			CourseID: courseID,
			TestID: "",
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
			TestID: "",
			UUID: uuid,
		}

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()
	
		sqlStatement := `UPDATE questiongroup SET name=$1, groupname=$2, numquestion=$3, score=$4, uuid=$5 WHERE id=$6`
	
		_, err = db.Exec(sqlStatement, g.Name, g.GroupName, g.NumQuestion, g.Score, g.UUID, g.ID)
		if err != nil {
			panic(err)
		}
	}
}

func getGroupInTest(courseID string, testID string) []byte{

	type GroupItem struct {
		ID string
		GroupName string
		NumQuestion string
		Score string
		Order int
	}

	type GroupInTest struct {
		Name string
		Items []GroupItem
	}

	type UUIDinGroup struct {
		Order int
		UUID string
	}

	groupTestMap := make(map[string]GroupInTest)

	var UUIDs []UUIDinGroup
	var uuid UUIDinGroup

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT uuid, headerorder FROM questiongroup WHERE courseid=$1 and testid=$2`
	rows, err := db.Query(sqlStatement, courseID, testID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&uuid.UUID, &uuid.Order)
		if err != nil {
			panic(err)
		}
		UUIDs = append(UUIDs, uuid)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	var uuidTemp UUIDinGroup

	for i := range UUIDs {
		for j := range UUIDs{
			if UUIDs[i].Order > UUIDs[j].Order{
				uuidTemp = UUIDs[i]
				UUIDs[i] = UUIDs[j]
				UUIDs[j] = uuidTemp
			}
		}
	}

	fmt.Println(UUIDs)

	for _, uuid := range UUIDs {

		var GroupItems []GroupItem

		var g GroupInTest
		var i GroupItem
		var groupTemp GroupItem

		sqlStatement := `SELECT name FROM questiongroup WHERE uuid=$1 and testid=$2`
		rows, err := db.Query(sqlStatement, uuid.UUID, testID)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&g.Name)
			if err != nil {
				panic(err)
			}
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}

		sqlStatement = `SELECT id, groupname, numquestion, score, grouporder FROM questiongroup WHERE uuid=$1 and testid=$2`
		rows, err = db.Query(sqlStatement, uuid.UUID, testID)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&i.ID, &i.GroupName, &i.NumQuestion, &i.Score, &i.Order)
			if err != nil {
				panic(err)
			}
			GroupItems = append(GroupItems, i)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}

		for i := range GroupItems {
			for j := range GroupItems{
				if GroupItems[i].Order > GroupItems[j].Order{
					groupTemp = GroupItems[i]
					GroupItems[i] = GroupItems[j]
					GroupItems[j] = groupTemp
				}
			}
		}

		g.Items = GroupItems
		groupTestMap[uuid.UUID] = g
	}

	b,err := json.Marshal(groupTestMap)
	if err != nil{
		panic(err)
	}

	return b
}

func getGroupInTestbank(courseID string) []byte{

	type GroupItem struct {
		ID string
		GroupName string
		NumQuestion string
		Score string
	}

	type GroupInTest struct {
		Name string
		Items []GroupItem
	}

	groupTestMap := make(map[string]GroupInTest)

	var UUIDs []string
	var uuid string

	/*var GroupItems []GroupItem

	var g GroupInTest
	var i GroupItem*/

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT uuid FROM questiongroup WHERE courseid=$1 and testid=''`
	rows, err := db.Query(sqlStatement, courseID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&uuid)
		if err != nil {
			panic(err)
		}
		UUIDs = append(UUIDs, uuid)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for _, uuid := range UUIDs {

		var GroupItems []GroupItem

		var g GroupInTest
		var i GroupItem

		sqlStatement := `SELECT name FROM questiongroup WHERE uuid=$1 and testid=''`
		rows, err := db.Query(sqlStatement, uuid)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&g.Name)
			if err != nil {
				panic(err)
			}
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}

		sqlStatement = `SELECT id, groupname, numquestion, score FROM questiongroup WHERE uuid=$1 and testid=''`
		rows, err = db.Query(sqlStatement, uuid)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&i.ID, &i.GroupName, &i.NumQuestion, &i.Score)
			if err != nil {
				panic(err)
			}
			GroupItems = append(GroupItems, i)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}

		g.Items = GroupItems
		groupTestMap[uuid] = g
	}

	b,err := json.Marshal(groupTestMap)
	if err != nil{
		panic(err)
	}

	return b
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

func checkQuestionGroupInTest(questionGroupID string, testID string) bool {
	var uuid string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlStatement := `SELECT id FROM questiongroup WHERE id=$1 and testid=$2;`
	row := db.QueryRow(sqlStatement, questionGroupID,testID)
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

func checkQuestionGroupInTestbank(questionGroupID string) bool {
	var uuid string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlStatement := `SELECT id FROM questiongroup WHERE id=$1 and testid='';`
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

func deleteQuestionGroupFromTest(questionInTest []string, testID string, courseID string){

	var questionID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT id FROM questiongroup WHERE testid=$1 and courseid=$2;`
	rows, err := db.Query(sqlStatement, testID, courseID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&questionID)
		if err != nil {
			panic(err)
		}

		check := true

		for _, id := range questionInTest{
			if questionID == id{
				check = false
			}
		}

		if check {
			sqlStatement := `DELETE from questiongroup WHERE id=$1 and testid=$2;`
			_, err = db.Exec(sqlStatement, questionID, testID)
			if err != nil {
				panic(err)
			}
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	
}

func deleteQuestionGroupFromTestbank(questionInTest []string, courseID string){

	var questionID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT id FROM questiongroup WHERE courseid=$1;`
	rows, err := db.Query(sqlStatement, courseID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&questionID)
		if err != nil {
			panic(err)
		}

		check := true

		for _, id := range questionInTest{
			if questionID == id{
				check = false
			}
		}

		if check {
			sqlStatement := `DELETE from questiongroup WHERE id=$1;`
			_, err = db.Exec(sqlStatement, questionID)
			if err != nil {
				panic(err)
			}
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
}

func allgrouptestlist(courseid string) []GrouptestList{
	var grouptestList []GrouptestList

	var grouptest GrouptestList

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT id, groupname FROM questiongroup WHERE courseid=$1 and testid='';`
	rows, err := db.Query(sqlStatement, courseid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&grouptest.ID, &grouptest.GroupName)
		if err != nil {
			panic(err)
		}
		grouptestList = append(grouptestList, grouptest)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return grouptestList
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

	fmt.Println(uuids)

	var input Input

	var questionInTest []string

	courseID := r.Header.Get("CourseID")

	testID := r.Header.Get("TestId")

	for headerorder, uuid := range uuids {
		input = objmap[uuid]
		for grouporder, item := range input.Items{
			groupTestListUpdate(input.Name, item.ID, item.GroupName, item.NumQuestion, item.Score, courseID, testID, uuid,headerorder,grouporder)
			questionInTest = append(questionInTest, item.ID)
		}
	}

	deleteQuestionGroupFromTest(questionInTest, testID, courseID)
	w.Write([]byte("success"))
})

//GetGroupInTest is a API that use to get all group test in the test.
var GetGroupInTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseID := r.Header.Get("CourseID")

	testID := r.Header.Get("TestId")

	w.Write(getGroupInTest(courseID,testID))
})

//AllGroupTestList is a API that use to get all grouptestlist in that course 
var AllGroupTestList = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseID := r.Header.Get("CourseID")

	json.NewEncoder(w).Encode(allgrouptestlist(courseID))
})

//TestbankUpdate is a API that use to get add or update group test list.
var TestbankUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

	var questionInTest []string

	courseID := r.Header.Get("CourseID")

	for _, uuid := range uuids {
		input = objmap[uuid]
		for _, item := range input.Items{
			testbankUpdate(input.Name, item.ID, item.GroupName, item.NumQuestion, item.Score, courseID, uuid)
			questionInTest = append(questionInTest, item.ID)
		}
	}

	deleteQuestionGroupFromTestbank(questionInTest, courseID)
	w.Write([]byte("success"))
})

//GetGroupInTestbank is a API that use to get all questiongroup in the testbsnk.
var GetGroupInTestbank = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseID := r.Header.Get("CourseID")

	w.Write(getGroupInTestbank(courseID))
})