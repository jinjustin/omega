package classroomcreatorcontroller

import (
	"omega/classroom"
	"fmt"
	//"encoding/json"
	"crypto/rand"
	//"github.com/jmoiron/sqlx"
	"omega/database"
	"database/sql"
)

//CreateNewClass is function that use to create classroom
func CreateNewClass(className string,classCode string,year string,permission string,userID string) []byte {

	classID := generateClassID()

	c := classroom.Classroom{
		ClassID: "",
		ClassName: "",
		ClassCode: "ClassCode Error",
		Year: "",
		Permission: "",
   }

	if(checkClassCode(classCode,year) == true){

		c = classroom.Classroom{
			ClassID: classID,
			ClassName: className,
			ClassCode: classCode,
			Year: year,
			Permission: permission,
	   }

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `INSERT INTO class (classid,classname,classcode, year, permission)VALUES ($1, $2, $3, $4, $5)`

		_, err = db.Exec(sqlStatement, c.ClassID,c.ClassName, c.ClassCode, c.Year, c.Permission)
		if err != nil {
		panic(err)
		}

		sqlStatement = `INSERT INTO userinclass (classid,userid)VALUES ($1, $2)`

		_, err = db.Exec(sqlStatement, c.ClassID, userID)
		if err != nil {
		panic(err)
		}

		return c.GetClassroomDetail()
	}

	return c.GetClassroomDetail()
}

func checkClassCode(classCode string,year string) bool{
	
	var classID string

	if len(classCode) != 8{
		return false
	}

	db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
	defer db.Close()
	sqlStatement := `SELECT classid FROM class WHERE classcode=$1 and year=$2;`
	row := db.QueryRow(sqlStatement, classCode,year)
	err = row.Scan(&classID)
	switch err {
	case sql.ErrNoRows: return true
	case nil: return false
	default: panic(err)
	}
}

//GenerateClassID is function that use to Generate Class ID
func generateClassID() string{
	n := 3
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
        panic(err)
    }
	s := fmt.Sprintf("%X", b)
	return s
}
