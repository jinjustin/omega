package questiongroupcontroller

import (
	"fmt"
	//"strings"

	//"github.com/golang/protobuf/descriptor"
	//"github.com/jinjustin/omega/question"

	"github.com/jinjustin/omega/questiongroup"
	"github.com/jinjustin/omega/question"
	"github.com/jinjustin/omega/questioncontroller"
	"github.com/jinjustin/omega/testcontroller"
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
	"github.com/iancoleman/orderedmap"
	//"github.com/mitchellh/mapstructure"
)




func groupTestListUpdate(name string, questiongroupID string, questiongroupName string, numQuestion string, maxQuestion string, score string, courseID string, testID string, uuid string, headerOrder int, groupOrder int)  error{
	
	var g questiongroup.QuestionGroup

	if (checkQuestionGroupInTest(questiongroupID, testID)){

		g = questiongroup.QuestionGroup{
			Name: name,
			ID: questiongroupID,
			GroupName: questiongroupName,
			NumQuestion: numQuestion,
			MaxQuestion: maxQuestion,
			Score: score,
			CourseID: courseID,
			TestID: testID,
			UUID: uuid,
			HeaderOrder: headerOrder,
			GroupOrder: groupOrder,
		}

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			return err
		}
		defer db.Close()


		if(checkQuestionGroupExist(questiongroupID)){
			sqlStatement := `INSERT INTO questiongroup (name, id, groupname, numquestion, maxquestion, score, courseid, testid, uuid, headerorder, grouporder)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
			_, err = db.Exec(sqlStatement, g.Name, g.ID, g.GroupName, g.NumQuestion, g.MaxQuestion, g.Score, g.CourseID, "", g.UUID, g.HeaderOrder, g.GroupOrder)
			if err != nil {
				return err
			}
		}

		sqlStatement := `INSERT INTO questiongroup (name, id, groupname, numquestion, maxquestion, score, courseid, testid, uuid, headerorder, grouporder)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
		_, err = db.Exec(sqlStatement, g.Name, g.ID, g.GroupName, g.NumQuestion, g.MaxQuestion, g.Score, g.CourseID, g.TestID, g.UUID, g.HeaderOrder, g.GroupOrder)
		if err != nil {
			return err
		}

	}else{

		g = questiongroup.QuestionGroup{
			Name: name,
			ID: questiongroupID,
			GroupName: questiongroupName,
			NumQuestion: numQuestion,
			MaxQuestion: maxQuestion,
			Score: score,
			CourseID: courseID,
			TestID: testID,
			UUID: uuid,
			HeaderOrder: headerOrder,
			GroupOrder: groupOrder,
		}

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			return err
		}
		defer db.Close()
	
		sqlStatement := `UPDATE questiongroup SET name=$1, groupname=$2, numquestion=$3, maxquestion=$4, score=$5, uuid=$6, headerorder=$7, grouporder=$8 WHERE id=$9`
	
		_, err = db.Exec(sqlStatement, g.Name, g.GroupName, g.NumQuestion, g.MaxQuestion, g.Score, g.UUID, g.HeaderOrder, g.GroupOrder, g.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func testbankUpdate(name string, questiongroupID string, questiongroupName string, numQuestion string, maxQuestion string, score string, courseID string, uuid string, headerOrder int, groupOrder int) error{
	
	var g questiongroup.QuestionGroup

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	if (checkQuestionGroupInTestbank(questiongroupID)){

		g = questiongroup.QuestionGroup{
			Name: name,
			ID: questiongroupID,
			GroupName: questiongroupName,
			NumQuestion: numQuestion,
			MaxQuestion: maxQuestion,
			Score: score,
			CourseID: courseID,
			TestID: "",
			UUID: uuid,
			HeaderOrder: headerOrder,
			GroupOrder: groupOrder,
		}

		sqlStatement := `INSERT INTO questiongroup (name, id, groupname, numquestion, maxquestion, score, courseid, testid, uuid, headerorder, grouporder)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
		_, err = db.Exec(sqlStatement, g.Name, g.ID, g.GroupName, g.NumQuestion, g.MaxQuestion, g.Score, g.CourseID, g.TestID, g.UUID, g.HeaderOrder, g.GroupOrder)
		if err != nil {
			return err
		}

	}else{

		g = questiongroup.QuestionGroup{
			Name: name,
			ID: questiongroupID,
			GroupName: questiongroupName,
			NumQuestion: numQuestion,
			MaxQuestion: maxQuestion,
			Score: score,
			CourseID: courseID,
			TestID: "",
			UUID: uuid,
			HeaderOrder: headerOrder,
			GroupOrder: groupOrder,
		}
	
		sqlStatement := `UPDATE questiongroup SET name=$1, groupname=$2, numquestion=$3, maxquestion=$4, score=$5, uuid=$6, headerorder=$7, grouporder=$8 WHERE id=$9`
	
		_, err = db.Exec(sqlStatement, g.Name, g.GroupName, g.NumQuestion, g.MaxQuestion, g.Score, g.UUID, g.HeaderOrder, g.GroupOrder, g.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func getGroupInTest(courseID string, testID string) ([]byte, error){

	type GroupItem struct {
		ID string `json:"id"`
		GroupName string `json:"groupName"`
		NumQuestion string `json:"numQuestion"`
		Score string `json:"score"`
		Order int `json:"order"`
		QuestionList []question.AllQuestionInGroup `json:"questionList"`
	}

	type GroupInTest struct {
		Name string `json:"name"`
		Items []GroupItem `json:"items"`
	}

	type UUIDinGroup struct {
		Order int
		UUID string
	}

	//groupTestMap := make(map[string]GroupInTest)
	o := orderedmap.New()

	var UUIDs []UUIDinGroup
	var uuid UUIDinGroup

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStatement := `SELECT uuid, headerorder FROM questiongroup WHERE courseid=$1 and testid=$2`
	rows, err := db.Query(sqlStatement, courseID, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&uuid.UUID, &uuid.Order)
		if err != nil {
			return nil ,err
		}

		check := true

		for _, u := range UUIDs{
			if u.UUID == uuid.UUID{
				check = false
			}
		}

		if check{
			UUIDs = append(UUIDs, uuid)
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	var uuidTemp UUIDinGroup

	for i := range UUIDs {
		for j := range UUIDs{
			if UUIDs[i].Order < UUIDs[j].Order{
				uuidTemp = UUIDs[i]
				UUIDs[i] = UUIDs[j]
				UUIDs[j] = uuidTemp
			}
		}
	}

	for _, uuid := range UUIDs {

		var GroupItems []GroupItem

		var allQuestionInGroup []question.AllQuestionInGroup

		var g GroupInTest
		var i GroupItem
		var groupTemp GroupItem

		var a question.AllQuestionInGroup

		sqlStatement := `SELECT name FROM questiongroup WHERE uuid=$1 and testid=$2`
		rows, err := db.Query(sqlStatement, uuid.UUID, testID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&g.Name)
			if err != nil {
				return nil, err
			}
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}

		sqlStatement = `SELECT id, groupname, numquestion, score, grouporder FROM questiongroup WHERE uuid=$1 and testid=$2`
		rows, err = db.Query(sqlStatement, uuid.UUID, testID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&i.ID, &i.GroupName, &i.NumQuestion, &i.Score, &i.Order)
			if err != nil {
				return nil, err
			}

			sqlStatement = `SELECT questionid, questionname FROM question WHERE testid=$1 and groupid=$2`
			rows, err = db.Query(sqlStatement, testID, i.ID)
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
			i.QuestionList = allQuestionInGroup

			GroupItems = append(GroupItems, i)
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}

		for i := range GroupItems {
			for j := range GroupItems{
				if GroupItems[i].Order < GroupItems[j].Order{
					groupTemp = GroupItems[i]
					GroupItems[i] = GroupItems[j]
					GroupItems[j] = groupTemp
				}
			}
		}

		g.Items = GroupItems

		o.Set(uuid.UUID,g)
	}


	b,err := o.MarshalJSON()
	if err != nil{
		return nil, err
	}

	return b, nil
}

func getGroupInTestbank(courseID string) ([]questiongroup.GroupItem, error){

	var allQuestionInGroup []question.AllQuestionInGroup

	var groupItems []questiongroup.GroupItem
	var i questiongroup.GroupItem
	var groupTemp questiongroup.GroupItem

	var a question.AllQuestionInGroup

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	fmt.Println(courseID)
	sqlStatement := `SELECT id, groupname, grouporder FROM questiongroup WHERE courseid=$1 and testid=''`
	rows, err := db.Query(sqlStatement, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&i.ID, &i.GroupName, &i.Order)
		if err != nil {
			return nil, err
		}

		fmt.Println(i)

		sqlStatement = `SELECT questionid, questionname FROM question WHERE testid='' and groupid=$1`
		rows, err = db.Query(sqlStatement, i.ID)
		fmt.Println(err)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&a.QuestionID, &a.QuestionName)
			fmt.Println(err)
			if err != nil {
				return nil, err
			}

			fmt.Println(a)
			allQuestionInGroup = append(allQuestionInGroup, a)
		}
		err = rows.Err()
		fmt.Println(err)
		if err != nil {
			return nil, err
		}
		i.QuestionList = allQuestionInGroup

		groupItems = append(groupItems, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	for i := range groupItems {
		for j := range groupItems{
			if groupItems[i].Order < groupItems[j].Order{
				groupTemp = groupItems[i]
				groupItems[i] = groupItems[j]
				groupItems[j] = groupTemp
			}
		}
	}

	return groupItems, nil
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

func deleteQuestionGroupFromTest(groupInTest []string, testID string, courseID string){

	var groupID string

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
		err = rows.Scan(&groupID)
		if err != nil {
			panic(err)
		}

		check := true

		for _, id := range groupInTest{
			if groupID == id{
				check = false
			}
		}

		if check {
			sqlStatement := `DELETE from questiongroup WHERE id=$1 and testid=$2;`
			_, err = db.Exec(sqlStatement, groupID, testID)
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

func allgrouptestlist(courseid string) ([]questiongroup.GrouptestList, error){
	var grouptestList []questiongroup.GrouptestList

	var grouptest questiongroup.GrouptestList

	var allQuestionInGroup []question.AllQuestionInGroup

	var questionInGroup question.AllQuestionInGroup

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return grouptestList, err
	}
	defer db.Close()

	sqlStatement := `SELECT id, groupname FROM questiongroup WHERE courseid=$1 and testid='';`
	rows, err := db.Query(sqlStatement, courseid)
	if err != nil {
		return grouptestList, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&grouptest.ID, &grouptest.GroupName)
		if err != nil {
			return grouptestList, err
		}

		sqlStatement = `SELECT questionid, questionname FROM question WHERE groupid=$1 and testid='';`
		rows, err := db.Query(sqlStatement, grouptest.ID)
		if err != nil {
			return grouptestList, err
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&questionInGroup.QuestionID, &questionInGroup.QuestionName)
			if err != nil {
				return grouptestList, err
			}
			
			allQuestionInGroup = append(allQuestionInGroup, questionInGroup)
		}
		err = rows.Err()
		if err != nil {
			return grouptestList, err
		}
		grouptest.QuestionList = allQuestionInGroup
		allQuestionInGroup = nil

		grouptestList = append(grouptestList, grouptest)
	}
	err = rows.Err()
	if err != nil {
		return grouptestList, err
	}

	return grouptestList, nil
}

//API

//GroupTestListUpdate is a API that use to get add or update group test list.
var GroupTestListUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	type Item struct {
		ID string `json:"id"`
		GroupName string `json:"groupName"`
		NumQuestion string `json:"numQuestion"`
		MaxQuestion string `json:"maxQuestion"`
		Score string `json:"score"`
		QuestionList []question.AllQuestionInGroup `json:"questionList"`
	}

	type Input struct {
		Name string `json:"name"`
		Items []Item `array:"item"`
	}

	var objmap map[string]Input

	o := orderedmap.New()

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

	err = json.Unmarshal(reqBody,&objmap)
	if err != nil{
		http.Error(w, "Can't convert JSON into map", http.StatusBadRequest)
            return
	}

	var input Input

	uuids := make([]string, 0, len(o.Keys()))
	for _, uuid := range o.Keys() {
        uuids = append(uuids, uuid)
    }

	var questionInTest []string

	var questionInGroup []string

	courseID := r.Header.Get("CourseID")

	testID := r.Header.Get("TestId")

	check := false

	if testID == ""{
		testID = testcontroller.GenerateTestID()
		check = true
	}

	fmt.Println("Recieve: ", objmap)

	for headerorder, uuid := range uuids {
		input = objmap[uuid]
		for grouporder, item := range input.Items{
			fmt.Println("Header Name: ", input.Name," Group ID: ", item.ID, " GroupName : ",item.GroupName," NumQuestion: ",item.NumQuestion," MaxQuestion: ",item.MaxQuestion, " Score: ",item.Score, " CourseID: ", courseID, " TestID: ", testID, " UUID: ", uuid)
			err = groupTestListUpdate(input.Name, item.ID, item.GroupName, item.NumQuestion, item.MaxQuestion, item.Score, courseID, testID, uuid,headerorder,grouporder)
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				fmt.Println(err)
            		return
			}
			questionInTest = append(questionInTest, item.ID)
			for _, questionItem := range item.QuestionList{
				fmt.Println(questionItem.QuestionName)
				err = questioncontroller.AddNewQuestion(item.ID, testID, questionItem.QuestionName, questionItem.QuestionID,"","")
				if err != nil{
					http.Error(w, err.Error(), http.StatusInternalServerError)
					fmt.Println(err)
            			return
				}
				questionInGroup = append(questionInGroup, questionItem.QuestionID) 
			}
			err = questioncontroller.DeleteQuestionFromGroupInTest(questionInGroup,testID,item.ID)
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
            		return
			}
			questionInGroup = nil
		}
	}

	deleteQuestionGroupFromTest(questionInTest, testID, courseID)
	w.WriteHeader(http.StatusOK)
	if check{
		w.Write([]byte(testID))
	}else{
		w.Write([]byte(""))
	}
})

//GetGroupInTest is a API that use to get all group test in the test.
var GetGroupInTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseID := r.Header.Get("CourseID")

	testID := r.Header.Get("TestId")

	g, err := getGroupInTest(courseID,testID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write(g)
})

//AllGroupTestList is a API that use to get all grouptestlist in that course 
var AllGroupTestList = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseID := r.Header.Get("CourseID")

	a, err := allgrouptestlist(courseID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a)
})

//TestbankUpdate is a API that use to get add or update group test list.
var TestbankUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	type Item struct {
		ID string `json:"id"`
		GroupName string `json:"groupName"`
		NumQuestion string `json:"numQuestion"`
		MaxQuestion string `json:"maxQuestion"`
		Score string `json:"score"`
		QuestionList []question.AllQuestionInGroup `json:"questionList"`
	}

	var items []Item

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil{
		http.Error(w, "can't read body", http.StatusBadRequest)
            return
	}

	err = json.Unmarshal(reqBody,&items)
	if err != nil{
		http.Error(w, "Can't convert JSON into map", http.StatusBadRequest)
            return
	}

	var questionInTest []string

	var questionInGroup []string

	courseID := r.Header.Get("CourseID")

		for grouporder, item := range items{
			err = testbankUpdate("", item.ID, item.GroupName, item.NumQuestion, item.MaxQuestion, item.Score, courseID, "", 0, grouporder)
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				fmt.Println(err)
            		return
			}
			questionInTest = append(questionInTest, item.ID)
			for _, questionItem := range item.QuestionList{
				fmt.Println(questionItem.QuestionName)
				err = questioncontroller.AddNewQuestion(item.ID, "", questionItem.QuestionName, questionItem.QuestionID,"","")
				if err != nil{
					http.Error(w, err.Error(), http.StatusInternalServerError)
					fmt.Println(err)
            			return
				}
				questionInGroup = append(questionInGroup, questionItem.QuestionID) 
			}
			err = questioncontroller.DeleteQuestionFromTestbank(questionInGroup,item.ID)
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
            		return
			}
			questionInGroup = nil
		}
	

	deleteQuestionGroupFromTestbank(questionInTest, courseID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})

//GetGroupInTestbank is a API that use to get all questiongroup in the testbsnk.
var GetGroupInTestbank = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseID := r.Header.Get("CourseID")

	allQuestionGroup, err := getGroupInTestbank(courseID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allQuestionGroup)
})