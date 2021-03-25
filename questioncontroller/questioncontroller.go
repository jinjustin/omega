package questioncontroller

import (
	"fmt"
	//"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"database/sql"

	"github.com/jinjustin/omega/choice"
	"github.com/jinjustin/omega/choicecontroller"
	"github.com/jinjustin/omega/database"
	"github.com/jinjustin/omega/question"
	//"github.com/sqs/goreturns/returns"

	//"omega/authentication"

	"encoding/json"
	"io/ioutil"
	"net/http"

	//"github.com/sqs/goreturns/returns"

	//"github.com/iancoleman/orderedmap"
)

//AddNewQuestion is a function that use to add new question to questiongroup.
func AddNewQuestion(groupID string, testID string, questionName string, questionID string, questionType string, data string) error{
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
			sqlStatement := `INSERT INTO question (testid, groupid, questionname, questionid, questiontype)VALUES ('', $1, $2, $3, $4)`
			_, err = db.Exec(sqlStatement, q.GroupID, q.QuestionName, q.QuestionID, q.QuestionType)
			if err != nil {
				return err
			}

			sqlStatement = `INSERT INTO questiondata (groupid, questionid, data)VALUES ($1, $2, $3)`
			_, err = db.Exec(sqlStatement, q.GroupID, q.QuestionID, q.Data)
			if err != nil {
				return err
			}

		}else if checkExist == nil{
			sqlStatement := `UPDATE question SET questionname=$1, questiontype=$2 WHERE questionid=$3`
			_, err = db.Exec(sqlStatement, q.QuestionName, q.QuestionType, q.QuestionID)
			if err != nil {
				return err
			}

			sqlStatement = `UPDATE questiondata SET data=$1 WHERE questionid=$2 and groupid=$3`
			_, err = db.Exec(sqlStatement,q.Data, q.QuestionID, q.GroupID)
			if err != nil {
				return err
			}

		}else{
			return checkExist
		}
	}else{

		checkInTest := checkQuestionInTest(questionID, testID)

		checkExist := checkQuestionExist(questionID)

		if checkInTest == sql.ErrNoRows{
			if checkExist == sql.ErrNoRows{
				sqlStatement := `INSERT INTO question (testid, groupid, questionname, questionid, questiontype)VALUES ('', $1, $2, $3, $4)`
				_, err = db.Exec(sqlStatement, q.GroupID, q.QuestionName, q.QuestionID ,q.QuestionType)
				if err != nil {
					return err
				}

				sqlStatement = `INSERT INTO question (testid, groupid, questionname, questionid, questiontype)VALUES ($1, $2, $3, $4, $5)`
				_, err = db.Exec(sqlStatement, q.TestID, q.GroupID, q.QuestionName, q.QuestionID, q.QuestionType)
				if err != nil {
					return err
				}

				sqlStatement = `INSERT INTO questiondata (groupid, questionid, data)VALUES ($1, $2, $3)`
				_, err = db.Exec(sqlStatement, q.GroupID, q.QuestionID, q.Data)
				if err != nil {
					return err
				}

			}else if checkExist == nil{
				sqlStatement := `INSERT INTO question (testid, groupid, questionname, questionid, questiontype)VALUES ($1, $2, $3, $4, $5)`
				_, err = db.Exec(sqlStatement, q.TestID, q.GroupID, q.QuestionName, q.QuestionID, q.QuestionType)
				if err != nil {
					return err
				}
			}else{
				return checkExist
			}
		}else if checkInTest == nil {
			sqlStatement := `UPDATE question SET questionname=$1, questiontype=$2 WHERE questionid=$3`
			_, err = db.Exec(sqlStatement, q.QuestionName, q.QuestionType, q.QuestionID)
			if err != nil {
				return err
			}

			sqlStatement = `UPDATE questiondata SET data=$1 WHERE questionid=$2 and groupid=$3`
			_, err = db.Exec(sqlStatement, q.Data, q.QuestionID, q.GroupID)
			if err != nil {
				return err
			}
		}else{
			return checkInTest
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

//DeleteQuestion is a function that use to delete question in function group.
func DeleteQuestion(testID string, questionID string) error {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	if testID==""{
		sqlStatement := `DELETE from question WHERE questionid=$1;`
		_, err = db.Exec(sqlStatement, questionID)
		if err != nil {
			return err
		}
	}else{
		sqlStatement := `DELETE from question WHERE questionid=$1 and testid=$2;`
		_, err = db.Exec(sqlStatement, questionID, testID)
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
		questionRows, err := db.Query(sqlStatement, groupID)
		if err != nil {
			return nil, err
		}
		defer questionRows.Close()
	
		for questionRows.Next() {
			err = questionRows.Scan(&a.QuestionID, &a.QuestionName)
			if err != nil {
				return nil, err
			}

			allQuestionInGroup = append(allQuestionInGroup, a)
		}
		err = questionRows.Err()
		if err != nil {
			return nil, err
		}
	}else{
		sqlStatement := `SELECT questionid, questionname FROM question WHERE testid=$1 and groupid=$2`
		questionRows, err := db.Query(sqlStatement, testID, groupID)
		if err != nil {
			return nil, err
		}
		defer questionRows.Close()
	
		for questionRows.Next() {
			err = questionRows.Scan(&a.QuestionID, &a.QuestionName)
			if err != nil {
				return nil, err
			}

			allQuestionInGroup = append(allQuestionInGroup, a)
		}
		err = questionRows.Err()
		if err != nil {
			return nil, err
		}
	}

	return allQuestionInGroup, err
}


func getAllQuestionInTest(courseID string, testID string) ([]byte, error) {

	var questionChoice choice.Choice

	var qwc question.WithChoice

	var questionWithChoices []question.WithChoice

	var questionChoices []choice.Choice

	var groupIDs []string

	var groupID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if testID == ""{

		sqlStatement := `SELECT id FROM questiongroup WHERE testid='' and courseid=$1`
		questiongroupRows, err := db.Query(sqlStatement, courseID)
		if err != nil {
			return nil, err
		}
		defer questiongroupRows.Close()
	
		for questiongroupRows.Next() {
			err = questiongroupRows.Scan(&groupID)
			if err != nil {
				return nil, err
			}

			groupIDs = append(groupIDs, groupID)
		}
		err = questiongroupRows.Err()
		if err != nil {
			return nil, err
		}

		for _, id := range groupIDs{
			sqlStatement := `SELECT questionid, questionname, questiontype FROM question WHERE testid='' and groupid=$1`
			questionRows, err := db.Query(sqlStatement, id)
			if err != nil {
				return nil, err
			}
			defer questionRows.Close()
		
			for questionRows.Next() {
				err = questionRows.Scan(&qwc.QuestionID, &qwc.QuestionName, &qwc.QuestionType)
				if err != nil {
					return nil, err
				}
				qwc.GroupID = id
				qwc.TestID = ""

				sqlStatement = `SELECT data FROM questiondata WHERE groupid=$1 and questionid=$2`
				questionDataRows, err := db.Query(sqlStatement, id, qwc.QuestionID)
				if err != nil {
					return nil, err
				}
				defer questionDataRows.Close()
			
				for questionDataRows.Next() {
					err = questionDataRows.Scan(&qwc.Data)
					if err != nil {
						return nil, err
					}
				}

				sqlStatement = `SELECT choiceid, data, imagelink, correctcheck FROM choice WHERE questionid=$1`
				choiceRows, err := db.Query(sqlStatement, qwc.QuestionID)
				if err != nil {
					return nil, err
				}
				defer choiceRows.Close()
			
				for choiceRows.Next() {
					err = choiceRows.Scan(&questionChoice.ChoiceID,&questionChoice.Data,&questionChoice.ImageLink.URL,&questionChoice.Check)
					if err != nil {
						return nil, err
					}
					questionChoice.ImageLink.UID = -1
					questionChoices = append(questionChoices, questionChoice)
				}

				if questionChoices == nil{
					qwc.ChoiceDetail = make([]choice.Choice,0)
				}else{
					qwc.ChoiceDetail = questionChoices
				}
				
				questionChoices = nil
	
				questionWithChoices = append(questionWithChoices, qwc)
			}
			err = questionRows.Err()
			if err != nil {
				return nil, err
			}
		} 

	}else{

		sqlStatement := `SELECT id FROM questiongroup WHERE testid=$1 and courseid=$2`
		questionGroupRows, err := db.Query(sqlStatement, testID, courseID)
		if err != nil {
			return nil, err
		}
		defer questionGroupRows.Close()
	
		for questionGroupRows.Next() {
			err = questionGroupRows.Scan(&groupID)
			if err != nil {
				return nil, err
			}

			groupIDs = append(groupIDs, groupID)
		}
		err = questionGroupRows.Err()
		if err != nil {
			return nil, err
		}

		for _, id := range groupIDs{
			sqlStatement := `SELECT questionid, questionname, questiontype FROM question WHERE testid=$1 and groupid=$2`
			questionRows, err := db.Query(sqlStatement,testID, id)
			if err != nil {
				return nil, err
			}
			defer questionRows.Close()
		
			for questionRows.Next() {
				err = questionRows.Scan(&qwc.QuestionID, &qwc.QuestionName, &qwc.QuestionType)
				if err != nil {
					return nil, err
				}
				qwc.GroupID = id
				qwc.TestID = testID

				sqlStatement = `SELECT data FROM questiondata WHERE groupid=$1 and questionid=$2`
				questionDataRows, err := db.Query(sqlStatement, id, qwc.QuestionID)
				if err != nil {
					return nil, err
				}
				defer questionDataRows.Close()
			
				for questionDataRows.Next() {
					err = questionDataRows.Scan(&qwc.Data)
					if err != nil {
						return nil, err
					}
				}

				sqlStatement = `SELECT choiceid, data, imagelink, correctcheck FROM choice WHERE questionid=$1`
				choiceRows, err := db.Query(sqlStatement, qwc.QuestionID)
				if err != nil {
					return nil, err
				}
				defer choiceRows.Close()
			
				for choiceRows.Next() {
					err = choiceRows.Scan(&questionChoice.ChoiceID,&questionChoice.Data,&questionChoice.ImageLink.URL,&questionChoice.Check)
					if err != nil {
						return nil, err
					}
					questionChoice.ImageLink.UID = -1
					questionChoices = append(questionChoices, questionChoice)
				}

				if questionChoices == nil{
					qwc.ChoiceDetail = make([]choice.Choice,0)
				}else{
					qwc.ChoiceDetail = questionChoices
				}
				
				questionChoices = nil
	
				questionWithChoices = append(questionWithChoices, qwc)
			}
			err = questionRows.Err()
			if err != nil {
				return nil, err
			}
		} 
	}

	b,err := json.Marshal(questionWithChoices)
	if err != nil{
		fmt.Println(err)
	}

	return b, nil
}

//DeleteQuestionFromGroupInTest is a function that use to auto delete question in questiongroup from that testID.
func DeleteQuestionFromGroupInTest(questionInGroup []string, testID string, groupID string) error{
	
	var questionID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `SELECT questionid FROM question WHERE testid=$1 and groupid=$2;`
	rows, err := db.Query(sqlStatement, testID, groupID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&questionID)
		if err != nil {
			return err
		}

		check := true

		for _, id := range questionInGroup{
			if questionID == id{
				check = false
			}
		}

		if check {
			sqlStatement := `DELETE from question WHERE questionid=$1 and groupid=$2 and testid=$3;`
			_, err = db.Exec(sqlStatement, questionID, groupID, testID)
			if err != nil {
				return err
			}
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

//DeleteQuestionFromTest is a function that use to auto delete question from test directly.
func DeleteQuestionFromTest(questionInGroup []string, testID string) error{
	
	var questionID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `SELECT questionid FROM question WHERE testid=$1;`
	rows, err := db.Query(sqlStatement, testID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&questionID)
		if err != nil {
			return err
		}

		check := true

		for _, id := range questionInGroup{
			if questionID == id{
				check = false
			}
		}

		if check {
			sqlStatement := `DELETE from question WHERE questionid=$1 and testid=$2;`
			_, err = db.Exec(sqlStatement, questionID, testID)
			if err != nil {
				return err
			}
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

//DeleteQuestionFromTestbank is a function that use to auto delete question in questiongroup.
func DeleteQuestionFromTestbank(questionInGroup []string, groupID string) error{
	
	var questionID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `SELECT questionid FROM question WHERE testid='' and groupid=$1;`
	rows, err := db.Query(sqlStatement, groupID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&questionID)
		if err != nil {
			return err
		}

		check := true

		for _, id := range questionInGroup{
			if questionID == id{
				check = false
			}
		}

		if check {
			sqlStatement := `DELETE from question WHERE questionid=$1 and groupid=$2;`
			_, err = db.Exec(sqlStatement, questionID, groupID)
			if err != nil {
				return err
			}
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
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

//UpdateQuestion is a API that use to add question to question group.
var UpdateQuestion = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	err := AddNewQuestion(groupID,testID,questionName,questionID,questionType,input.Data)

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error Contact JJ immediately!"))
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
		w.Write([]byte("500 - Internal Server Error Contact JJ immediately!"))
		fmt.Println(err)
	}else{
		w.WriteHeader(http.StatusOK)
		w.Write(q.GetQuestionDetail())
	}
})

//DeleteQuestion is a API that use to delete question in question group (in test or in all question group).
/*var DeleteQuestion = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	testID := r.Header.Get("TestID")
	questionID := r.Header.Get("QuestionID")

	err := deleteQuestion(testID,questionID)

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error Contact JJ immediately!"))
		fmt.Println(err)
	}else{
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 - OK"))
	}
})*/

//GetAllQuestionInGroup is a API that use to get id and name of all question in question group.
var GetAllQuestionInGroup = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	testID := r.Header.Get("TestID")
	groupID := r.Header.Get("GroupId")

	allQuestionInGroup, err := getAllQuestionInGroup(testID,groupID)

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error Contact JJ immediately!"))
		fmt.Println(err)
	}else{
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(allQuestionInGroup)
	}
})

//UpdateAllQuestionInTest is a function that use to update all question in test.
var UpdateAllQuestionInTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var questionWithChoices []question.WithChoice

	var questionInTest []string

	var choiceInQuestion []string

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil{
		http.Error(w, "can't read body", http.StatusBadRequest)
            return
	}

	err = json.Unmarshal(reqBody,&questionWithChoices)
	if err != nil{
		http.Error(w, "Can't convert JSON into map", http.StatusBadRequest)
            return
	}

	testID := r.Header.Get("TestId")

	for _, q := range questionWithChoices {
		err = AddNewQuestion(q.GroupID, testID, q.QuestionName, q.QuestionID, q.QuestionType, q.Data)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
				return
		}

		for _, choice := range q.ChoiceDetail{
			err = choicecontroller.AddNewChoice(choice.ChoiceID, q.QuestionID, choice.Data, choice.ImageLink.URL, choice.Check)
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				fmt.Println(err)
            		return
			}

			choiceInQuestion = append(choiceInQuestion, choice.ChoiceID)
		}
		choicecontroller.DeleteChoiceFromQuestion(choiceInQuestion, q.QuestionID)
		choiceInQuestion = nil
		questionInTest = append(questionInTest, q.QuestionID)
	}

	DeleteQuestionFromTest(questionInTest, testID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})


//GetAllQuestionInTest is a API that use to get information of all question in test.
var GetAllQuestionInTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseID := r.Header.Get("CourseID")
	testID := r.Header.Get("TestID")

	allQuestionInTest, err := getAllQuestionInTest(courseID,testID)

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
			return
		}
	w.WriteHeader(http.StatusOK)
	w.Write(allQuestionInTest)
})
