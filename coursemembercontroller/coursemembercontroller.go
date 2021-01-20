package coursemembercontroller

import(
	"omega/student"
	"omega/teacher"
	"omega/database"
	"database/sql"
	"encoding/json"
)

//AddStudentToCourse is ฟังก์ชันใช้สำหรับให้ผู้สอนเชิญผู้เรียนเข้าไปใน Course หรือผู้เรียนต้องการเข้าร่วม Course
func AddStudentToCourse(studentID string,courseCode string) []byte{

	var userID string
	var firstName string
	var surName string

	s := student.Student{
		UserID: "Can't Invite this student",
		StudentID: "",
		Firstname: "",
		Surname: "",
		Email: "",
	}

   db, err := sql.Open("postgres", database.PsqlInfo())
   if err != nil {
	   panic(err)
   }
   defer db.Close()

   sqlStatement := `SELECT userid,firstname,surname FROM student WHERE studentid=$1;`
   rows,err := db.Query(sqlStatement, studentID)
   if err != nil {
	   panic(err)
   }
   defer rows.Close()

   for rows.Next() {
	err = rows.Scan(&userID,&firstName,&surName)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
	panic(err)
	}

	if(checkMemberInCourse(userID,courseCode)){
		s = student.Student{
			UserID: "",
			StudentID: studentID,
			Firstname: firstName,
			Surname: surName,
			Email: "",
		}
		
		sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`

		_, err = db.Exec(sqlStatement, courseCode, userID, "student", "pending")
		if err != nil {
			panic(err)
		}
	}

	return s.GetStudentDetail()
}

//AddTeacherToCourse is ฟังก์ชันที่ใช้เมื่อ ผู้สอน ต้องการเข้าร่วม Course
func AddTeacherToCourse(username string,courseCode string) []byte{

	var userID string
	var firstName string
	var surName string

	t := teacher.Teacher{
		UserID: "Can't Join this course",
		Firstname: "",
		Surname: "",
		Email: "",
	}

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

	sqlStatement = `SELECT firstname,surname FROM teacher WHERE userid=$1;`
	rows,err = db.Query(sqlStatement, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
 
	for rows.Next() {
	 err = rows.Scan(&firstName,&surName)
		 if err != nil {
			 panic(err)
		 }
	 }
	 err = rows.Err()
	 if err != nil {
	 panic(err)
	 }

	if(checkMemberInCourse(userID,courseCode)){
		t = teacher.Teacher{
			UserID: "",
			Firstname: firstName,
			Surname: surName,
			Email: "",
		}

		sqlStatement = `INSERT INTO coursemember (coursecode, userid, role, status)VALUES ($1, $2, $3, $4)`

		_, err = db.Exec(sqlStatement, courseCode, userID, "teacher", "pending")
		if err != nil {
		panic(err)
		}
	}

	return t.GetTeacherDetail()
}

//GetStudentInCourse is a function that use to get data of all students in the course
func GetStudentInCourse(courseCode string) []byte{
	var userIDs []string

	type studentInCourse struct{
		StudentID string
		Firstname string
		Surname string
		Status string
	}

	var studentInCourses []studentInCourse

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM coursemember WHERE coursecode=$1 and role=$2 and status=$3;`
	rows, err := db.Query(sqlStatement, courseCode,"student","join")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
		userIDs = append(userIDs, userID)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for _, a := range userIDs {
		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `SELECT studentid,firstname,surname FROM student WHERE userid=$1;`
		rows, err := db.Query(sqlStatement, a)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var studentID string
			var firstname string
			var surname string

			err = rows.Scan(&studentID, &firstname, &surname)
			if err != nil {
				panic(err)
			}

			s := studentInCourse{
				StudentID: studentID,
				Firstname: firstname,
				Surname: surname,
				Status: "join",
			}

			studentInCourses = append(studentInCourses, s)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
	}

	userIDs = nil

	sqlStatement = `SELECT userid FROM coursemember WHERE coursecode=$1 and role=$2 and status=$3;`
	rows, err = db.Query(sqlStatement, courseCode,"student","pending")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
		userIDs = append(userIDs, userID)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for _, a := range userIDs {
		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `SELECT studentid,firstname,surname FROM student WHERE userid=$1;`
		rows, err := db.Query(sqlStatement, a)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var studentID string
			var firstname string
			var surname string

			err = rows.Scan(&studentID, &firstname, &surname)
			if err != nil {
				panic(err)
			}

			s := studentInCourse{
				StudentID: studentID,
				Firstname: firstname,
				Surname: surname,
				Status: "pending",
			}

			studentInCourses = append(studentInCourses, s)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
	}

	b,err := json.Marshal(studentInCourses)
	if err != nil{
		panic(err)
	}

	return b
}

//GetTeacherInCourse is a function that use to get data of all students in the course
func GetTeacherInCourse(courseCode string) []byte{
	var userIDs []string

	type teacherInCourse struct{
		Firstname string
		Surname string
		Status string
	}

	var teacherInCourses []teacherInCourse

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT userid FROM coursemember WHERE coursecode=$1 and role=$2 and status=$3;`
	rows, err := db.Query(sqlStatement, courseCode,"teacher","join")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
		userIDs = append(userIDs, userID)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for _, a := range userIDs {
		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `SELECT firstname,surname FROM teacher WHERE userid=$1;`
		rows, err := db.Query(sqlStatement, a)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var firstname string
			var surname string

			err = rows.Scan(&firstname, &surname)
			if err != nil {
				panic(err)
			}

			t := teacherInCourse{
				Firstname: firstname,
				Surname: surname,
				Status: "join",
			}

			teacherInCourses = append(teacherInCourses, t)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
	}

	userIDs = nil

	sqlStatement = `SELECT userid FROM coursemember WHERE coursecode=$1 and role=$2 and status=$3;`
	rows, err = db.Query(sqlStatement, courseCode,"teacher","pending")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			panic(err)
		}
		userIDs = append(userIDs, userID)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for _, a := range userIDs {
		db, err := sql.Open("postgres", database.PsqlInfo())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sqlStatement := `SELECT firstname,surname FROM teacher WHERE userid=$1;`
		rows, err := db.Query(sqlStatement, a)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var firstname string
			var surname string

			err = rows.Scan(&firstname, &surname)
			if err != nil {
				panic(err)
			}

			t := teacherInCourse{
				Firstname: firstname,
				Surname: surname,
				Status: "pending",
			}

			teacherInCourses = append(teacherInCourses, t)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
	}

	b,err := json.Marshal(teacherInCourses)
	if err != nil{
		panic(err)
	}

	return b
}

//ApproveJoinCourse is ฟังก์ชันสำหรับให้ผู้สอนรองรับการเข้าร่วม Course ของผู้เรียน
func ApproveJoinCourse(userID string,courseCode string) string{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `UPDATE coursemember SET status=$1 WHERE coursecode=$2 and userid=$3;`

	_, err = db.Exec(sqlStatement,"join", courseCode, userID)
	if err != nil {
	panic(err)
	}

	return "success"
}

//DeclineJoinCourse is ฟังก์ชันสำหรับให้ผู้สอนรองรับการเข้าร่วม Course ของผู้เรียน
func DeclineJoinCourse(userID string,courseCode string) string{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `DELETE from coursemember WHERE coursecode=$1 and userid=$2;`

	_, err = db.Exec(sqlStatement,courseCode, userID)
	if err != nil {
	panic(err)
	}

	return "success"
}

//GetUserRole is a function that use to get user role (student or teacher)
func GetUserRole(username string) string{
	var role string

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT role FROM users WHERE username=$1;`
	rows, err := db.Query(sqlStatement, username)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&role)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return role
}

//DeleteTeacherInCourse is function that use to delete teacher in the course.
func DeleteTeacherInCourse(courseCode string, username string) []byte {
	t := teacher.Teacher{
		UserID: "Can't find.",
		Firstname:   "",
		Surname: "",
		Email:       "",
	}

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

	var firstname string
	var surname string
	var email string

	sqlStatement = `SELECT firstname, surname, email FROM teacher WHERE userid=$1;`
	row := db.QueryRow(sqlStatement, userID)
	err = row.Scan(&firstname, &surname, &email)
	if err != nil {
		panic(err)
	}

	t = teacher.Teacher{
		Firstname: firstname,
		Surname:   surname,
		Email: email,
	}

	sqlStatement = `DELETE FROM coursemember WHERE userID=$1 and coursecode=$2;`
	_, err = db.Exec(sqlStatement, userID,courseCode)
	if err != nil {
		panic(err)
	}

	return t.GetTeacherDetail()
}

//DeleteStudentInCourse is function that use to delete student in the course.
func DeleteStudentInCourse(courseCode string, username string) []byte {
	s := student.Student{
		UserID: "Can't find.",
		StudentID: "",
		Firstname:   "",
		Surname: "",
		Email:       "",
	}

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

	var studentID string
	var firstname string
	var surname string
	var email string

	sqlStatement = `SELECT studentid, firstname, surname, email FROM student WHERE userid=$1;`
	row := db.QueryRow(sqlStatement, userID)
	err = row.Scan(&studentID,&firstname, &surname, &email)
	if err != nil {
		panic(err)
	}

	s = student.Student{
		StudentID: studentID,
		Firstname: firstname,
		Surname:   surname,
		Email: email,
	}

	sqlStatement = `DELETE FROM coursemember WHERE userID=$1 and coursecode=$2;`
	_, err = db.Exec(sqlStatement, userID,courseCode)
	if err != nil {
		panic(err)
	}

	return s.GetStudentDetail()

}

func checkMemberInCourse(userID string,courseCode string) bool{

	db, err := sql.Open("postgres", database.PsqlInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var status string

	sqlStatement := `SELECT status FROM coursemember WHERE userid=$1 and coursecode=$2;`
	row := db.QueryRow(sqlStatement, userID,courseCode)
	err = row.Scan(&status)
	switch err {
	case sql.ErrNoRows: return true
	case nil: return false
	default: panic(err)
	}
}

