package classroomdeletercontroller

import (
	"omega/classroom"
	"omega/database"
	"database/sql"
)

//DeleteClassroom is use to delete classroom
func DeleteClassroom(classID string,userID string) []byte{
	c := classroom.Classroom{
		ClassID: "Can't find.",
		ClassCode: "",
		ClassName: "",
		Year: "",
		Permission: "",
	}

	if(checkClass(classID)==true && checkUser(classID,userID) == true){
		var className string
		var classCode string
		var year string
		var permission string

		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `SELECT classname,classcode,year,permission FROM class WHERE classid=$1;`
		row := db.QueryRow(sqlStatement, classID)
		err = row.Scan(&className,&classCode,&year,&permission)
		if err != nil{
			panic(err)
		}

		c := classroom.Classroom{
			ClassID: classID,
			ClassCode: classCode,
			ClassName: className,
			Year: year,
			Permission: permission,
		}

		sqlStatement = `DELETE FROM class WHERE classid=$1;`
		_, err = db.Exec(sqlStatement, classID)
		if err != nil {
		panic(err)
		}

		sqlStatement = `DELETE FROM userinclass WHERE classid=$1 and userid=$2;`
		_, err = db.Exec(sqlStatement, classID,userID)
		if err != nil {
		panic(err)
		}
		return c.GetClassroomDetail()
	}
	return c.GetClassroomDetail()
}

func checkClass(classID string) bool{
	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var className string

	sqlStatement := `SELECT classname FROM class WHERE classid=$1;`
	row := db.QueryRow(sqlStatement, classID)
	err = row.Scan(&className)
	switch err {
	case sql.ErrNoRows: return false
	case nil: return true
	default: panic(err)
	}
}

func checkUser(classID string,userID string) bool{
	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var check string

	sqlStatement := `SELECT classid FROM userinclass WHERE userid=$1 and classid=$2;`
	row := db.QueryRow(sqlStatement, userID,classID)
	err = row.Scan(&check)
	switch err {
	case sql.ErrNoRows: return false
	case nil: return true
	default: panic(err)
	}
}