package answercontroller

import (
	"fmt"
	//"github.com/jinjustin/omega/answer"

	"encoding/json"
	//"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"database/sql"

	"github.com/jinjustin/omega/answer"
	"github.com/jinjustin/omega/database"
	"github.com/jinjustin/omega/question"
	"github.com/jinjustin/omega/authentication"

	//"github.com/jinjustin/omega/authentication"

	"net/http"
	//"encoding/json"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

/*func submitAnswer(testID string, studentID string, qwc []question.AndChoiceWithoutCorrectCheck) error {

	var data string
	var groupID string
	var score string
	var a answer.Info
	var studentAnswer []answer.Info

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	for num, a := range studentAnswer{
		sqlStatement := `SELECT data FROM questiondata WHERE questionid=$1;`
		rows, err := db.Query(sqlStatement, a.QuestionID)
		if err != nil {
			return err
		}
		defer rows.Close()
	
		for rows.Next() {
			err = rows.Scan(&data)
			if err != nil {
				return err
			}
		}
		err = rows.Err()
		if err != nil {
			return err
		}

		studentAnswer[num].Data = data

		sqlStatement = `SELECT groupid FROM question WHERE questionid=$1 and testid=$2;`
		questionRows, err := db.Query(sqlStatement, a.QuestionID, testID)
		if err != nil {
			return err
		}
		defer questionRows.Close()
	
		for questionRows.Next() {
			err = questionRows.Scan(&groupID)
			if err != nil {
				return err
			}
		}
		err = questionRows.Err()
		if err != nil {
			return err
		}

		sqlStatement = `SELECT score FROM questiongroup WHERE id=$1 and testid=$2;`
		questionGroupRows, err := db.Query(sqlStatement, groupID, testID)
		if err != nil {
			return err
		}
		defer questionGroupRows.Close()
	
		for questionGroupRows.Next() {
			err = questionGroupRows.Scan(&score)
			if err != nil {
				return err
			}
		}
		err = questionGroupRows.Err()
		if err != nil {
			return err
		}

		studentAnswer[num].MaxScore = score
	}

	b, err := json.Marshal(studentAnswer)
	if err != nil {
		panic(err)
	}

	checkExist := checkAnswerExist(testID, studentID)

	if checkExist == sql.ErrNoRows{
		sqlStatement := `INSERT INTO answer (testid, studentid, studentanswer, totalscore, checkedanswer, completepercent)VALUES ($1, $2, $3, $4, $5, $6);`
		_, err = db.Exec(sqlStatement, testID, studentID, b, "0", "0", "0.00")
		if err != nil {
			return err
		}
	}else if checkExist == nil{
		sqlStatement := `UPDATE answer SET studentanswer=$1 WHERE testid=$2 and studentid=$3;`
		_, err = db.Exec(sqlStatement, b, testID, studentID)
		if err != nil {
			return err
		}
	}else{
		return checkExist
	}

	err = autoScoring(studentAnswer,testID,studentID)
	if err != nil{
		return err
	}

	return nil
}*/

func submitAnswer(testID string, studentID string, questionAndChoiceWithoutAnswer []question.AndChoiceWithoutCorrectCheck) error {

	var groupID string
	var score string
	var studentAnswer []answer.Info

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	for _, qwc := range questionAndChoiceWithoutAnswer{
		var a answer.Info
		a.QuestionID = qwc.QuestionID
		a.QuestionType = qwc.QuestionType
		a.Data = qwc.Data

		sqlStatement := `SELECT groupid FROM question WHERE questionid=$1 and testid=$2;`
		questionRows, err := db.Query(sqlStatement, a.QuestionID, testID)
		if err != nil {
			return err
		}
		defer questionRows.Close()
	
		for questionRows.Next() {
			err = questionRows.Scan(&groupID)
			if err != nil {
				return err
			}
		}
		err = questionRows.Err()
		if err != nil {
			return err
		}

		sqlStatement = `SELECT score FROM questiongroup WHERE id=$1 and testid=$2;`
		questionGroupRows, err := db.Query(sqlStatement, groupID, testID)
		if err != nil {
			return err
		}
		defer questionGroupRows.Close()
	
		for questionGroupRows.Next() {
			err = questionGroupRows.Scan(&score)
			if err != nil {
				return err
			}
		}
		err = questionGroupRows.Err()
		if err != nil {
			return err
		}

		a.MaxScore = score
		if a.QuestionType == "Choice"{
			for _, ans := range qwc.ChoiceDetail{
				if ans.Answer == "true"{
					a.Answer = append(a.Answer, ans.ChoiceID)
				}
			}
		}else if a.QuestionType == "Pair"{
			for _, ans := range qwc.ChoiceDetail{
				pairAnswer := ans.ChoiceID + ":" + ans.Answer
				a.Answer = append(a.Answer, pairAnswer)
			}
		}else if a.QuestionType == "Short Answer"{
			for _, ans := range qwc.ChoiceDetail{
				shortAnswer := ans.Answer
				if shortAnswer == ""{
					a.Answer = append(a.Answer, "")
				}else{
					a.Answer = append(a.Answer, ans.Answer)
				}
			}
		}else if a.QuestionType == "Write-up"{
			for _, ans := range qwc.ChoiceDetail{
				writeAnswer := ans.Answer
				if writeAnswer == ""{
					a.Answer = append(a.Answer, "")
				}else{
					a.Answer = append(a.Answer, ans.Answer)
				}
			}
		}else if a.QuestionType == "Upload Answer"{
			for _, ans := range qwc.ChoiceDetail{
				writeAnswer := ans.Answer
				if writeAnswer == ""{
					a.Answer = append(a.Answer, "")
				}else{
					a.Answer = append(a.Answer, ans.Answer)
				}
			}
		}
		a.Score = "-"
		studentAnswer = append(studentAnswer, a)
	}

	b, err := json.Marshal(studentAnswer)
	if err != nil {
		panic(err)
	}

	checkExist := checkAnswerExist(testID, studentID)

	if checkExist == sql.ErrNoRows{
		sqlStatement := `INSERT INTO answer (testid, studentid, studentanswer, totalscore, checkedanswer, completepercent)VALUES ($1, $2, $3, $4, $5, $6);`
		_, err = db.Exec(sqlStatement, testID, studentID, b, "0", "0", "0.00")
		if err != nil {
			return err
		}
	}else if checkExist == nil{
		sqlStatement := `UPDATE answer SET studentanswer=$1 WHERE testid=$2 and studentid=$3;`
		_, err = db.Exec(sqlStatement, b, testID, studentID)
		if err != nil {
			return err
		}
	}else{
		return checkExist
	}

	err = autoScoring(studentAnswer,testID,studentID)
	if err != nil{
		return err
	}

	return nil
}



func getStudentAnswer(testID string, studentID string, uuid string) ([]answer.Info, error){

	var b []byte

	var allStudentAnswer []answer.Info

	var selectedStudentAnswer []answer.Info

	var groupID string
	var groupIDs []string

	var questionID string
	var questionIDs []string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStatement := `SELECT studentanswer FROM answer WHERE testid=$1 and studentid=$2;`
	rows, err := db.Query(sqlStatement, testID, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&b)
		if err != nil {
			return nil ,err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b,&allStudentAnswer)
	if err != nil{
        return nil, err
	}

	sqlStatement = `SELECT id FROM questiongroup WHERE uuid=$1 and testid=$2;`
	rows, err = db.Query(sqlStatement, uuid, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&groupID)
		if err != nil {
			return nil ,err
		}
		groupIDs = append(groupIDs, groupID)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	for _, g := range groupIDs{
		sqlStatement = `SELECT questionid FROM question WHERE groupid=$1 and testid=$2;`
		rows, err = db.Query(sqlStatement, g, testID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	
		for rows.Next() {
			err = rows.Scan(&questionID)
			if err != nil {
				return nil ,err
			}
			questionIDs = append(questionIDs, questionID)
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}
	}

	for _, id := range questionIDs{
		for _, a := range allStudentAnswer{
			if id == a.QuestionID && (a.QuestionType != "Choice" && a.QuestionType != "Pair" && a.QuestionType != "Short Answer"){
				selectedStudentAnswer = append(selectedStudentAnswer, a)
			} 
		}
	}
	return selectedStudentAnswer, err
}

func getAllStudentAnswerInformation(testID string) ([]answer.StudentAnswerInformation, error){

	var studentAnswerInfo answer.StudentAnswerInformation

	var studentanswerInfos []answer.StudentAnswerInformation

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStatement := `SELECT studentid, completepercent FROM answer WHERE testid=$1;`
	answerRows, err := db.Query(sqlStatement, testID)
	if err != nil {
		return nil, err
	}
	defer answerRows.Close()

	for answerRows.Next() {
		err = answerRows.Scan(&studentAnswerInfo.StudentID, &studentAnswerInfo.CompletePercent)
		if err != nil {
			return nil ,err
		}

		sqlStatement = `SELECT firstname, surname FROM student WHERE studentid=$1;`
		studentRows, err := db.Query(sqlStatement, studentAnswerInfo.StudentID)
		if err != nil {
			return nil, err
		}
		defer studentRows.Close()
		for studentRows.Next() {
			err = studentRows.Scan(&studentAnswerInfo.Firstname, &studentAnswerInfo.Surname)
			if err != nil {
				return nil ,err
			}
		}
		err = studentRows.Err()
		if err != nil {
			return nil, err
		}
		studentanswerInfos = append(studentanswerInfos, studentAnswerInfo)
	}
	err = answerRows.Err()
	if err != nil {
		return nil, err
	}

	for num1, i := range studentanswerInfos{
		for num2, j := range studentanswerInfos{
			iStudentIDint, _ := strconv.Atoi(i.StudentID)
			jStudentIDint, _ := strconv.Atoi(j.StudentID)

			if iStudentIDint < jStudentIDint{
				studentanswerInfos[num1], studentanswerInfos[num2] = studentanswerInfos[num2], studentanswerInfos[num1]
			}
		}
	}

	return studentanswerInfos, err
}

func scoringAnswer(testID string, studentID string, questionID string, score string) error{

	_, errTest := strconv.Atoi(score)

	if errTest != nil{
		return errTest
	}

	var b []byte
	var allStudentAnswer []answer.Info
	var checkedAnswer string
	var totalScore string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `SELECT studentanswer, checkedanswer, totalscore FROM answer WHERE testid=$1 and studentid=$2;`
	answerRows, err := db.Query(sqlStatement, testID, studentID)
	if err != nil {
		return nil
	}
	defer answerRows.Close()

	for answerRows.Next() {
		err = answerRows.Scan(&b,&checkedAnswer,&totalScore)
		if err != nil {
			return nil
		}
	}
	err = answerRows.Err()
	if err != nil {
		return nil
	}

	err = json.Unmarshal(b,&allStudentAnswer)
	if err != nil{
		return err
	}

	checkedAnswerf, err := strconv.ParseFloat(checkedAnswer, 64)
	if err != nil{
		return err
	}

	for num, a := range allStudentAnswer{

		if a.QuestionID == questionID{
			if a.Score == "-"{
				checkedAnswerf += 1.0
			}
			allStudentAnswer[num].Score = score
		}
	}

	completePercent := (checkedAnswerf/float64(len(allStudentAnswer)))*100

	b, err = json.Marshal(allStudentAnswer)
	if err != nil {
		panic(err)
	}

	totalScoreInt, err := strconv.Atoi(totalScore)
	if err != nil{
		return err
	}

	newScoreInt, err := strconv.Atoi(score)
	if err != nil{
		return err
	}

	totalScoreInt += newScoreInt

	totalscoreString := strconv.Itoa(totalScoreInt)

	completePercentString := fmt.Sprintf("%.2f", completePercent)
	checkedAnswerString := fmt.Sprintf("%.0f", checkedAnswerf)
	
	sqlStatement = `UPDATE answer SET studentanswer=$1, totalscore=$2, checkedanswer=$3, completepercent=$4 WHERE testid=$5 and studentid=$6;`
	
	_, err = db.Exec(sqlStatement, b, totalscoreString,checkedAnswerString, completePercentString, testID, studentID)
	if err != nil {
		return err
	}

	return nil
}

func autoScoring(studentAnswer []answer.Info, testID string, studentID string) error{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	var correctcheck string
	var data string
	var groupID string
	var score string
	
	totalScore := 0
	checkedAnswer := 0.0

	for num, a := range studentAnswer{

		if a.QuestionType == "Choice"{

			sqlStatement := `SELECT groupid FROM question WHERE testid=$1 and questionid=$2;`
			questionRows, err := db.Query(sqlStatement, testID, a.QuestionID)
			if err != nil {
				return err
			}
			defer questionRows.Close()

			for questionRows.Next() {
				err = questionRows.Scan(&groupID)
				if err != nil {
					return err
				}
			}
			err = questionRows.Err()
			if err != nil {
				return err
			}

			sqlStatement = `SELECT score FROM questiongroup WHERE id=$1;`
			questionGroupRows, err := db.Query(sqlStatement, groupID)
			if err != nil {
				return err
			}
			defer questionGroupRows.Close()

			for questionGroupRows.Next() {
				err = questionGroupRows.Scan(&score)
				if err != nil {
					return err
				}
			}
			err = questionGroupRows.Err()
			if err != nil {
				return err
			}

			check := true

			for _, ans := range a.Answer{
				sqlStatement = `SELECT correctcheck FROM choice WHERE choiceid=$1 and questionid=$2;`
				choiceRows, err := db.Query(sqlStatement, ans, a.QuestionID)
				if err != nil {
					return err
				}
				defer choiceRows.Close()
	
				for choiceRows.Next() {
					err = choiceRows.Scan(&correctcheck)
					if err != nil {
						return err
					}
				}
				err = choiceRows.Err()
				if err != nil {
					return err
				}

				if correctcheck == "false"{
					check = false
				}
			}

			if check{
				studentAnswer[num].Score = score
				s, _ := strconv.Atoi(score)
				totalScore += s
			}
			checkedAnswer += 1.0
		}else if a.QuestionType == "Pair"{
			sqlStatement := `SELECT groupid FROM question WHERE testid=$1 and questionid=$2;`
			questionRows, err := db.Query(sqlStatement, testID, a.QuestionID)
			if err != nil {
				return err
			}
			defer questionRows.Close()

			for questionRows.Next() {
				err = questionRows.Scan(&groupID)
				if err != nil {
					return err
				}
			}
			err = questionRows.Err()
			if err != nil {
				return err
			}

			sqlStatement = `SELECT score FROM questiongroup WHERE id=$1;`
			questionGroupRows, err := db.Query(sqlStatement, groupID)
			if err != nil {
				return err
			}
			defer questionGroupRows.Close()

			for questionGroupRows.Next() {
				err = questionGroupRows.Scan(&score)
				if err != nil {
					return err
				}
			}
			err = questionGroupRows.Err()
			if err != nil {
				return err
			}

			check := true

			for _, ans := range a.Answer{

				pairs := strings.Split(ans, ":")

				sqlStatement = `SELECT correctcheck FROM choice WHERE choiceid=$1 and questionid=$2;`
				choiceRows, err := db.Query(sqlStatement, pairs[0], a.QuestionID)
				if err != nil {
					return err
				}
				defer choiceRows.Close()
	
				for choiceRows.Next() {
					err = choiceRows.Scan(&correctcheck)
					if err != nil {
						return err
					}
				}
				err = choiceRows.Err()
				if err != nil {
					return err
				}

				if correctcheck != pairs[1]{
					check = false
				}
			}

			if check{
				studentAnswer[num].Score = score
				s, _ := strconv.Atoi(score)
				totalScore += s
			}
			checkedAnswer += 1.0
		}else if a.QuestionType == "Short Answer"{
			sqlStatement := `SELECT groupid FROM question WHERE testid=$1 and questionid=$2;`
			questionRows, err := db.Query(sqlStatement, testID, a.QuestionID)
			if err != nil {
				return err
			}
			defer questionRows.Close()

			for questionRows.Next() {
				err = questionRows.Scan(&groupID)
				if err != nil {
					return err
				}
			}
			err = questionRows.Err()
			if err != nil {
				return err
			}

			sqlStatement = `SELECT score FROM questiongroup WHERE id=$1;`
			questionGroupRows, err := db.Query(sqlStatement, groupID)
			if err != nil {
				return err
			}
			defer questionGroupRows.Close()

			for questionGroupRows.Next() {
				err = questionGroupRows.Scan(&score)
				if err != nil {
					return err
				}
			}
			err = questionGroupRows.Err()
			if err != nil {
				return err
			}

			check := true

			for _, ans := range a.Answer{

				sqlStatement = `SELECT data FROM choice WHERE questionid=$1;`
				choiceRows, err := db.Query(sqlStatement, a.QuestionID)
				if err != nil {
					return err
				}
				defer choiceRows.Close()
	
				for choiceRows.Next() {
					err = choiceRows.Scan(&data)
					if err != nil {
						return err
					}
				}
				err = choiceRows.Err()
				if err != nil {
					return err
				}

				if data != ans{
					check = false
				}
			}

			if check{
				studentAnswer[num].Score = score
				s, _ := strconv.Atoi(score)
				totalScore += s
			}
			checkedAnswer += 1.0
		}
	}


	completePercent := (checkedAnswer/float64(len(studentAnswer)))*100

	b, err := json.Marshal(studentAnswer)
	if err != nil {
		panic(err)
	}

	totalscoreString := strconv.Itoa(totalScore)

	completePercentString := fmt.Sprintf("%.2f", completePercent)
	checkedAnswerString := fmt.Sprintf("%.0f", checkedAnswer)
	
	sqlStatement := `UPDATE answer SET studentanswer=$1, totalscore=$2, checkedanswer=$3, completepercent=$4 WHERE testid=$5 and studentid=$6;`
	
	_, err = db.Exec(sqlStatement, b, totalscoreString,checkedAnswerString, completePercentString, testID, studentID)
	if err != nil {
		return err
	}

	return nil
}

func CalculateStatistic(testID string) (answer.StatisticValue ,error) {

	var statisticValue answer.StatisticValue

	var score string
	var scores []string

	var max int
	var min int
	var mean float64
	var sd float64
	totalScore := 0
	totalValue := 0.00

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return statisticValue, err
	}
	defer db.Close()

	sqlStatement := `SELECT totalscore FROM answer WHERE testid=$1;`
	answerRows, err := db.Query(sqlStatement, testID)
	if err != nil {
		return statisticValue, err
	}
	defer answerRows.Close()

	for answerRows.Next() {
		err = answerRows.Scan(&score)
		if err != nil {
			return statisticValue, nil
		}
		fmt.Println("score: ", score)
		scores = append(scores, score)
	}
	err = answerRows.Err()
	if err != nil {
		return statisticValue, nil
	}

	fmt.Println("scores: ",scores)

	if scores != nil{
		max, err = strconv.Atoi(scores[0])
		if err != nil{
			return statisticValue, nil
		}
		min, err = strconv.Atoi(scores[0])
		if err != nil{
			return statisticValue, nil
		}

		for _, s := range scores{
			iScore, err := strconv.Atoi(s)
			if err != nil{
				return statisticValue, nil
			}
			if iScore > max{
				max = iScore
			}
			if iScore < min{
				min = iScore
			}

			totalScore += iScore
		}

		fmt.Println("totalscore: ", totalScore)

		mean = float64(totalScore)/float64(len(scores))

		fmt.Println("mean: ",mean)

		for _, s := range scores{
			iScore, err := strconv.Atoi(s)
			if err != nil{
				return statisticValue, nil
			}
			fScore := float64(iScore)

			value := fScore - mean

			value2 := value*value
			
			totalValue += value2
		}

		sd = math.Sqrt((totalValue/float64(len(scores)-1)))
		statisticValue.Max = strconv.Itoa(max)
		statisticValue.Min = strconv.Itoa(min)
		statisticValue.Mean = fmt.Sprintf("%.2f", mean)
		statisticValue.SD = fmt.Sprintf("%.2f", sd)
	}

	return statisticValue, nil
}

func studentGetScore(studentID string) ([]answer.StudentScore, error){

	var studentScore []answer.StudentScore
	var ss answer.StudentScore

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStatement := `SELECT testid, totalscore FROM answer WHERE studentid=$1 and completepercent='100.00';`
	answerRows, err := db.Query(sqlStatement, studentID)
	if err != nil {
		return nil, err
	}
	defer answerRows.Close()

	for answerRows.Next() {
		err = answerRows.Scan(&ss.TestID,&ss.TotalScore)
		if err != nil {
			return nil, err
		}
		studentScore = append(studentScore, ss)
	}
	err = answerRows.Err()
	if err != nil {
		return nil, err
	}

	var topic string
	for num, s := range studentScore{
		sqlStatement := `SELECT topic FROM test WHERE testid=$1;`
		testRows, err := db.Query(sqlStatement, s.TestID)
		if err != nil {
			return nil, err
		}
		defer testRows.Close()
	
		for testRows.Next() {
			err = testRows.Scan(&topic)
			if err != nil {
				return nil, err
			}
		}
		err = testRows.Err()
		if err != nil {
			return nil, err
		}

		studentScore[num].Topic = topic

		statisticValue, err := CalculateStatistic(s.TestID)
		if err != nil{
			return nil, err
		}
		studentScore[num].Max = statisticValue.Max
		studentScore[num].Min = statisticValue.Min
		studentScore[num].Mean = statisticValue.Mean
		studentScore[num].SD = statisticValue.SD
	}

	return studentScore, nil
}

func checkAnswerExist(testID string, studentID string) error {

	var dummy string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()
	sqlStatement := `SELECT checkedanswer FROM answer WHERE testid=$1 and studentid=$2;`
	row := db.QueryRow(sqlStatement, testID,studentID)
	err = row.Scan(&dummy)
	return err
}

//API

//SubmitAnswer is a function that use to store student answer to database.
var SubmitAnswer = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var questionAndChoiceWithoutCorrectCheck []question.AndChoiceWithoutCorrectCheck

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil{
		http.Error(w, "can't read body", http.StatusBadRequest)
            return
	}

	err = json.Unmarshal(reqBody,&questionAndChoiceWithoutCorrectCheck)
	if err != nil{
		http.Error(w, "Can't convert JSON into map", http.StatusBadRequest)
            return
	}

	fmt.Println(questionAndChoiceWithoutCorrectCheck)

	testID := r.Header.Get("TestId")

	studentID := authentication.GetUsername(r)

	err = submitAnswer(testID, studentID, questionAndChoiceWithoutCorrectCheck)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})

//GetAnswer is a function that use to get student answer from database.
var GetAnswer = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	testID := r.Header.Get("TestId")

	studentID := r.Header.Get("StudentID")

	uuid := r.Header.Get("UUID")

	selectedAnswer, err := getStudentAnswer(testID, studentID, uuid)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
            return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(selectedAnswer)
})

//GetAnswer is a function that use to get student answer from database.
var GetAllStudentAnswerInformation = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	testID := r.Header.Get("TestId")

	studentAnswerInformations, err := getAllStudentAnswerInformation(testID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(studentAnswerInformations)
})


var ScoringAnswer = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	testID := r.Header.Get("TestId")
	studentID := r.Header.Get("studentID")

	type Input struct{
		QuestionID string `json:"questionID"`
		Score string `json:"score"`
	}

	var input Input

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil{
		http.Error(w, "can't read body", http.StatusBadRequest)
            return
	}

	err = json.Unmarshal(reqBody, &input)
	if err != nil{
		http.Error(w, "Can't convert", http.StatusBadRequest)
            return
	}

	err = scoringAnswer(testID, studentID, input.QuestionID, input.Score)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
})

var GetStatisticValue = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	testID := r.Header.Get("TestId")

	statisticValue, err := CalculateStatistic(testID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(statisticValue)
})

var StudentGetScore = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	studentID := authentication.GetUsername(r)

	studentScores, err := studentGetScore(studentID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(studentScores)
})