package choicecontroller

import (
	//"fmt"
	//"encoding/json"

	"database/sql"
	"github.com/jinjustin/omega/database"

	"github.com/jinjustin/omega/choice"
)

//AddNewChoice is a function that use to update choice in question.
func AddNewChoice(choiceID string, questionID string, data string, imageLink string, check string) error{
	var c choice.Choice

	io := choice.ImageObject{
		UID: -1,
		URL: imageLink,
	}

	var ios []choice.ImageObject

	ios = append(ios, io)

	c = choice.Choice{
		ChoiceID: choiceID,
		QuestionID: questionID,
		Data: data,
		ImageLink: ios,
		Check: check,
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	checkExist := checkChoiceExist(choiceID, questionID)

	if checkExist == sql.ErrNoRows{
		sqlStatement := `INSERT INTO choice (choiceid, questionid, data, imagelink, correctcheck)VALUES ($1, $2, $3, $4, $5)`
		_, err = db.Exec(sqlStatement, c.ChoiceID, c.QuestionID, c.Data, c.ImageLink[0].URL, c.Check)
		if err != nil {
			return err
		}

	}else if checkExist == nil{
		sqlStatement := `UPDATE choice SET data=$1, imagelink=$2, correctcheck=$3 WHERE choiceid=$4`
		_, err = db.Exec(sqlStatement, c.Data, c.ImageLink[0].URL, c.Check, c.ChoiceID)
		if err != nil {
			return err
		}
	}else{
		return checkExist
	}
	
	return nil
}

func checkChoiceExist(choiceID string, questionID string) error {
	var correctCheck string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()
	sqlStatement := `SELECT correctcheck FROM choice WHERE choiceid=$1 and questionid=$2;`
	row := db.QueryRow(sqlStatement, choiceID, questionID)
	err = row.Scan(&correctCheck)
	return err
}

//DeleteChoiceFromQuestion is a function that use to delete choice in question.
func DeleteChoiceFromQuestion(choiceInQuestion []string, questionID string) error{

	var choiceID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStatement := `SELECT choiceid FROM choice WHERE questionid=$1;`
	rows, err := db.Query(sqlStatement, questionID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&choiceID)
		if err != nil {
			return err
		}

		check := true

		for _, id := range choiceInQuestion{
			if choiceID == id{
				check = false
			}
		}

		if check {
			sqlStatement := `DELETE from choice WHERE choiceid=$1;`
			_, err = db.Exec(sqlStatement, choiceID)
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