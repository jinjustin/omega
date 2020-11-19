package classroomlistcontroller

import (
	"omega/classroom"
	"omega/database"
	"database/sql"
)


//GetClassroomList is use to get all classrooms that user is being member.
func GetClassroomList(userID string) []classroom.Classroom{

	var classIDs []string

	var classes []classroom.Classroom

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT classid FROM userinclass WHERE userid=$1;`
	rows,err := db.Query(sqlStatement, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var classID string
		err = rows.Scan(&classID)
			if err != nil {
				panic(err)
			}
		classIDs = append(classIDs,classID)
		}

		err = rows.Err()
		if err != nil {
		panic(err)
		}

	for _, a := range classIDs {
		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()
	
		sqlStatement := `SELECT classname,classcode,year,permission FROM class WHERE classid=$1;`
		rows,err := db.Query(sqlStatement, a)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
	
		for rows.Next() {
			var className string
			var classCode string
			var year string
			var permission string

			err = rows.Scan(&className,&classCode,&year,&permission)
				if err != nil {
					panic(err)
				}
			
			c := classroom.Classroom{
				ClassID: a,
				ClassCode: classCode,
				ClassName: className,
				Year: year,
				Permission: permission,
			}
			
			classes = append(classes,c)
			}
	
			err = rows.Err()
			if err != nil {
			panic(err)
			}
	}

	return classes
}