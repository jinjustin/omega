package main

import (
	"fmt"
	"net/http"
	"omega/coursecontroller"
	"omega/coursemembercontroller"
	"omega/teachercontroller"
	"omega/login"

	//"omega/database"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"context"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/dgrijalva/jwt-go"
	//"omega/student"
	//"database/sql"
	//_ "github.com/lib/pq"
)

func createCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseName string
		CourseID string
		Year string
		Permission string
		Announcement string
		Description string
		Username string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	fmt.Println(input)
	w.Write(coursecontroller.CreateCourse(input.CourseName,input.CourseID,input.Year,input.Permission,input.Announcement,input.Description,input.Username))
}

var getCourseList = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	type Input struct{
		Username string
		//Token string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(coursecontroller.GetCourseList(input.Username))
})

func deleteCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
		Username string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursecontroller.DeleteCourse(input.CourseCode,input.Username))
}

func getTeacherInfo(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		Username string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(teachercontroller.GetTeacherInfo(input.Username))
}

func editTeacherInfo(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		Firstname string
		Surname string
		Email string
		Username string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(teachercontroller.EditTeacherInfo(input.Firstname,input.Surname,input.Email,input.Username))
}

func getDescription(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(coursecontroller.GetDescription(input.CourseCode))
}

func editDescription(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
		Description string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(coursecontroller.EditDescription(input.CourseCode,input.Description))
}

func getAnnouncement(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(coursecontroller.GetAnnouncement(input.CourseCode))
}

func editAnnouncement(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
		Announcement string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	json.NewEncoder(w).Encode(coursecontroller.EditAnnouncement(input.CourseCode,input.Announcement))
}

func addStudentToCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		StudentID string
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursemembercontroller.AddStudentToCourse(input.StudentID,input.CourseCode))
}

func addTeacherToCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		Username string
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursemembercontroller.AddTeacherToCourse(input.Username,input.CourseCode))
}

func getStudentInCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursemembercontroller.GetStudentInCourse(input.CourseCode))
}

func getTeacherInCourse(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		CourseCode string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var input Input
	json.Unmarshal(reqBody, &input)
	w.Write(coursemembercontroller.GetTeacherInCourse(input.CourseCode))
}

var getUserRole = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	type Input struct{
		Username string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	w.Write([]byte(coursemembercontroller.GetUserRole(input.Username)))
})

var deleteTeacherInCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	type Input struct{
		CourseCode string
		Username string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	w.Write([]byte(coursemembercontroller.DeleteTeacherInCourse(input.CourseCode,input.Username)))
})

var deleteStudentInCourse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	type Input struct{
		CourseCode string
		Username string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
    var input Input
	json.Unmarshal(reqBody, &input)
	w.Write([]byte(coursemembercontroller.DeleteStudentInCourse(input.CourseCode,input.Username)))
})

func test(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
    fmt.Fprintf(w, "Welcome to the HomePage!")
}

func test2(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	type Input struct{
		test string
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	w.Write(reqBody)
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("secret"), nil
			})
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), "props", claims)
				// Access context values in handlers like this
				// props, _ := r.Context().Value("props").(jwt.MapClaims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}
		}
	})
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET","POST","DELETE"},
	})

	myRouter.HandleFunc("/",test)
	myRouter.HandleFunc("/test",test2).Methods("POST")
	myRouter.HandleFunc("/createcourse",createCourse).Methods("POST")
	myRouter.Handle("/getcourselist",middleware(getCourseList)).Methods("POST")
	myRouter.HandleFunc("/deletecourse",deleteCourse).Methods("POST")
	myRouter.HandleFunc("/getteacherinfo",getTeacherInfo).Methods("POST")
	myRouter.HandleFunc("/editteacherinfo",editTeacherInfo).Methods("POST")
	myRouter.HandleFunc("/getdescription",getDescription).Methods("POST")
	myRouter.HandleFunc("/editdescription",editDescription).Methods("POST")
	myRouter.HandleFunc("/getannouncement",getAnnouncement).Methods("POST")
	myRouter.HandleFunc("/editannouncement",editAnnouncement).Methods("POST")
	myRouter.HandleFunc("/addstudent",addStudentToCourse).Methods("POST")
	myRouter.HandleFunc("/addteacher",addTeacherToCourse).Methods("POST")
	myRouter.HandleFunc("/getstudentincourse",getStudentInCourse).Methods("POST")
	myRouter.HandleFunc("/getteacherincourse",getTeacherInCourse).Methods("POST")
	myRouter.Handle("/deleteteacherincourse",deleteTeacherInCourse).Methods("POST")
	myRouter.Handle("/deletestudentincourse",deleteStudentInCourse).Methods("POST")
	myRouter.Handle("/getuserrole",getUserRole).Methods("POST")
	myRouter.HandleFunc("/login", login.Login).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000",c.Handler(myRouter)))

	
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	handleRequests()
}