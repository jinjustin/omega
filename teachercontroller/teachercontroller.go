package teachercontroller

import(
	"omega/teacher"
	"omega/database"
	"database/sql"
)

//GetTeacherInfo is a function that use to get teacher information by input teacher username
func GetTeacherInfo(username string) []byte{

	var userID string
	var firstname string
	var surname string
	var email string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM users WHERE username=$1;`
	rows, err := db.Query(sqlStatement, username)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	sqlStatement = `SELECT firstname,surname,email FROM teacher WHERE userid=$1;`
	rows, err = db.Query(sqlStatement, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&firstname,&surname,&email)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	t := teacher.Teacher{
		UserID: "",
		Firstname: firstname,
		Surname: surname,
		Email: email,
	}

	return t.GetTeacherDetail()
}

//EditTeacherInfo is a function that use to edit teacher information
func EditTeacherInfo(firstname string,surname string,email string,username string)[]byte{
	
	var userID string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM users WHERE username=$1;`
	rows, err := db.Query(sqlStatement, username)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	sqlStatement = `UPDATE teacher SET firstname = $1,surname = $2,email = $3 WHERE userid = $4`

	_, err = db.Exec(sqlStatement, firstname, surname, email, userID)
	if err != nil {
		panic(err)
	}

	sqlStatement = `SELECT firstname,surname,email FROM teacher WHERE userid=$1;`
	rows, err = db.Query(sqlStatement, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&firstname,&surname,&email)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	t := teacher.Teacher{
		UserID: "",
		Firstname: firstname,
		Surname: surname,
		Email: email,
	}

	return t.GetTeacherDetail()
}