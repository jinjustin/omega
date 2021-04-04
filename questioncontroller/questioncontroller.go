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
	"github.com/jinjustin/omega/testcontroller"
	"github.com/jinjustin/omega/authentication"
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
				fmt.Println("Create new question and Add to test.")
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
				fmt.Println("Add question to test.")
				sqlStatement := `INSERT INTO question (testid, groupid, questionname, questionid, questiontype)VALUES ($1, $2, $3, $4, $5)`
				_, err = db.Exec(sqlStatement, q.TestID, q.GroupID, q.QuestionName, q.QuestionID, q.QuestionType)
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
		}else if checkInTest == nil {
			fmt.Println("Update question in test.")

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

	//var questionChoice choice.Choice

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
					var questionChoice choice.Choice
					var io choice.ImageObject
					var url string
					err = choiceRows.Scan(&questionChoice.ChoiceID,&questionChoice.Data,&url,&questionChoice.Check)
					if err != nil {
						return nil, err
					}
					if url == ""{
						questionChoice.ImageLink = make([]choice.ImageObject, 0)
					}else{
						io.URL = url
						io.UID = "-1"
						questionChoice.ImageLink = make([]choice.ImageObject, 1)
						questionChoice.ImageLink[0] = io
						//questionChoice.ImageLink = append(questionChoice.ImageLink,io)
					}
					
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
					var questionChoice choice.Choice
					var io choice.ImageObject
					var url string
					err = choiceRows.Scan(&questionChoice.ChoiceID,&questionChoice.Data,&url,&questionChoice.Check)
					if err != nil {
						return nil, err
					}
					if url == ""{
						questionChoice.ImageLink = make([]choice.ImageObject, 0)
					}else{
						io.URL = url
						io.UID = "-1"
						questionChoice.ImageLink = make([]choice.ImageObject, 1)
						questionChoice.ImageLink[0] = io
						//questionChoice.ImageLink = append(questionChoice.ImageLink,io)
					}
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

	if questionWithChoices == nil{
		questionWithChoices = make([]question.WithChoice,0)
	}

	b,err := json.Marshal(questionWithChoices)
	if err != nil{
		return nil, err
	}

	return b, nil
}

func getAllQuestionForTest(courseID string) ([]byte, error) {

	var choiceWIthoutCorrectCheck choice.WithoutCorrectCheck

	var qac question.AndChoiceWithoutCorrectCheck

	var questionAndChoicesWithoutCorrectChecks []question.AndChoiceWithoutCorrectCheck

	var choiceWIthoutCorrectChecks []choice.WithoutCorrectCheck

	var groupIDs []string

	var groupID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

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
			err = questionRows.Scan(&qac.QuestionID, &qac.QuestionName, &qac.QuestionType)
			if err != nil {
				return nil, err
			}
			qac.GroupID = id

			sqlStatement = `SELECT data FROM questiondata WHERE groupid=$1 and questionid=$2`
			questionDataRows, err := db.Query(sqlStatement, id, qac.QuestionID)
			if err != nil {
				return nil, err
			}
			defer questionDataRows.Close()
		
			for questionDataRows.Next() {
				err = questionDataRows.Scan(&qac.Data)
				if err != nil {
					return nil, err
				}
			}

			sqlStatement = `SELECT choiceid, data, imagelink FROM choice WHERE questionid=$1`
			choiceRows, err := db.Query(sqlStatement, qac.QuestionID)
			if err != nil {
				return nil, err
			}
			defer choiceRows.Close()
		
			for choiceRows.Next() {
				var io choice.ImageObject
				var url string
				err = choiceRows.Scan(&choiceWIthoutCorrectCheck.ChoiceID,&choiceWIthoutCorrectCheck.Data,&url)
				if err != nil {
					return nil, err
				}
				if url == ""{
					choiceWIthoutCorrectCheck.ImageLink = make([]choice.ImageObject, 0)
				}else{
					io.URL = url
					io.UID = "-1"
					choiceWIthoutCorrectCheck.ImageLink = append(choiceWIthoutCorrectCheck.ImageLink,io)
				}
				
				choiceWIthoutCorrectChecks = append(choiceWIthoutCorrectChecks, choiceWIthoutCorrectCheck)
			}

			if choiceWIthoutCorrectChecks == nil{
				qac.ChoiceDetail = make([]choice.WithoutCorrectCheck,0)
			}else{
				qac.ChoiceDetail = choiceWIthoutCorrectChecks
			}
			
			choiceWIthoutCorrectChecks = nil
			questionAndChoicesWithoutCorrectChecks = append(questionAndChoicesWithoutCorrectChecks, qac)
		}
		err = questionRows.Err()
		if err != nil {
			return nil, err
		}
	} 

	if questionAndChoicesWithoutCorrectChecks == nil{
		questionAndChoicesWithoutCorrectChecks = make([]question.AndChoiceWithoutCorrectCheck,0)
	}

	b,err := json.Marshal(questionAndChoicesWithoutCorrectChecks)
	if err != nil{
		return nil, err
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

func UpdateQuestionToTest(testID string, questionID string, groupID string, questionName string) error{

	var questiontype string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `SELECT questiontype FROM question WHERE questionid=$1 and groupid=$2`
	questionGroupRows, err := db.Query(sqlStatement, questionID, groupID)
	if err != nil {
		return err
	}
	defer questionGroupRows.Close()

	for questionGroupRows.Next() {
		err = questionGroupRows.Scan(&questiontype)
		if err != nil {
			return err
		}
	}
	err = questionGroupRows.Err()
	if err != nil {
		return err
	}

	checkInTest := checkQuestionInTest(questionID, testID)

	if checkInTest == sql.ErrNoRows{
		sqlStatement = `INSERT INTO question (testid, groupid, questionid, questionname, questiontype)VALUES ($1, $2, $3, $4, $5)`
		_, err = db.Exec(sqlStatement, testID, groupID, questionID, questionName, questiontype)
		if err != nil {
			return err
		}
	}else if err != nil{
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

			sqlStatement = `DELETE from questiondata WHERE questionid=$1;`
			_, err = db.Exec(sqlStatement, questionID)
			if err != nil {
				return err
			}

			sqlStatement = `DELETE from choice WHERE questionid=$1;`
			_, err = db.Exec(sqlStatement, questionID)
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

func getExam(testID string, studentID string) ([]byte, error) {
	var exam []byte

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlStatement := `SELECT exam FROM exam WHERE testid=$1 and studentid=$2;`
	row := db.QueryRow(sqlStatement, testID, studentID)
	err = row.Scan(&exam)
	if err != nil{
		return nil, err
	}
	return exam, nil
}

func checkExamExist(testID string, studentID string)(error) {
	var exam []byte

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()
	sqlStatement := `SELECT exam FROM exam WHERE testid=$1 and studentid=$2;`
	row := db.QueryRow(sqlStatement, testID, studentID)
	err = row.Scan(&exam)
	return err
}

func inputExam(exam []byte, testID string, studentID string) (error) {

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	checkExist := checkExamExist(testID, studentID)

	if checkExist == sql.ErrNoRows{
		sqlStatement := `INSERT INTO exam (exam, testid, studentid)VALUES ($1, $2, $3)`
		_, err = db.Exec(sqlStatement, exam, testID, studentID)
		if err != nil {
			return err
		}
	}else if checkExist == nil{
		sqlStatement := `UPDATE exam SET exam=$1 WHERE testid=$2 and studentid=$3`

		_, err = db.Exec(sqlStatement, exam, testID, studentID)
		if err != nil {
			return err
		}
	}else if checkExist != nil{
		return err
	}
	return nil
}


//API

//UpdateQuestion is a API that use to add question to question group.
var UpdateQuestion = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Data string
	}

	groupID := r.Header.Get("GroupId")
	testID := r.Header.Get("TestId")
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
		return 
	}else{
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 - OK"))
	}
})

//GetQuestion is a API that use to get question information in question group.
var GetQuestion = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	groupID := r.Header.Get("GroupId")
	testID := r.Header.Get("TestId")
	questionID := r.Header.Get("QuestionID")

	q, err := getQuestion(groupID,testID,questionID)

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error Contact JJ immediately!"))
		return
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

	testID := r.Header.Get("TestId")
	groupID := r.Header.Get("GroupId")

	allQuestionInGroup, err := getAllQuestionInGroup(testID,groupID)

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error Contact JJ immediately!"))
	}else{
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(allQuestionInGroup)
	}
})

//UpdateAllQuestionInTest is a function that use to update all question in test.
var UpdateAllQuestionInTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var questionWithChoices []question.WithChoice

	//var questionInTest []string

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
	fmt.Println(testID)
	check := false

	if testID == ""{
		fmt.Println("Enter Case")
		testID = testcontroller.GenerateTestID()
		check = true
	}

	for _, q := range questionWithChoices {
		//fmt.Println(q.GroupID, testID, q.QuestionName, q.QuestionID, q.QuestionType, q.Data)
		err = AddNewQuestion(q.GroupID, "", q.QuestionName, q.QuestionID, q.QuestionType, q.Data)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
				return
		}

		for _, choice := range q.ChoiceDetail{

			imageLink := ""
			for _, i := range choice.ImageLink{
				imageLink = i.URL
			}
			err = choicecontroller.AddNewChoice(choice.ChoiceID, q.QuestionID, choice.Data, imageLink, choice.Check)
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
            		return
			}

			choiceInQuestion = append(choiceInQuestion, choice.ChoiceID)
		}
		choicecontroller.DeleteChoiceFromQuestion(choiceInQuestion, q.QuestionID)
		choiceInQuestion = nil
		//questionInTest = append(questionInTest, q.QuestionID)
	}

	//DeleteQuestionFromTest(questionInTest, testID)
	w.WriteHeader(http.StatusOK)
	if check{
		w.Write([]byte(testID))
	}else{
		w.Write([]byte(""))
	}
	
})


//GetAllQuestionInTest is a API that use to get information of all question in test.
var GetAllQuestionInTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseID := r.Header.Get("CourseID")
	testID := r.Header.Get("TestId")

	allQuestionInTest, err := getAllQuestionInTest(courseID,testID)

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	w.WriteHeader(http.StatusOK)
	w.Write(allQuestionInTest)
})

//GetAllQuestionForTest is a API that use to get information of all question in course without answer.
var GetAllQuestionForTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	courseID := r.Header.Get("CourseID")
	testID := r.Header.Get("TestId")
	studentID := authentication.GetUsername(r)
	

	exam, err := getExam(testID, studentID)
	if err == sql.ErrNoRows{
		allQuestionAndChoiceWithoutCorrectcheck, err := getAllQuestionForTest(courseID)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		w.WriteHeader(http.StatusOK)
		w.Write(allQuestionAndChoiceWithoutCorrectcheck)
	}else{
		w.WriteHeader(http.StatusOK)
		w.Write(exam)
	}
})

//GetAllQuestionForTest is a API that use to get information of all question in course without answer.
var InputExam = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	testID := r.Header.Get("TestId")
	studentID := authentication.GetUsername(r)

	reqBody, _ := ioutil.ReadAll(r.Body)

	err := inputExam(reqBody, testID, studentID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Ok"))
})
