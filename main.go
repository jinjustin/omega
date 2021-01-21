package main

import (
	"fmt"
	"net/http"
	"omega/coursecontroller"
	"omega/coursemembercontroller"
	"omega/login"
	"omega/teachercontroller"
	"omega/testcontroller"

	//"omega/database"
	"context"
	//"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	//"omega/student"
	//"database/sql"
	//_ "github.com/lib/pq"
)

func testAPI(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
}

func test2(w http.ResponseWriter, r *http.Request){
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

	myRouter.HandleFunc("/",testAPI)
	myRouter.HandleFunc("/test",test2).Methods("POST")
	myRouter.Handle("/createcourse",coursecontroller.CreateCourse).Methods("POST")
	myRouter.Handle("/getcourselist",middleware(coursecontroller.GetCourseList)).Methods("POST")
	myRouter.Handle("/deletecourse",coursecontroller.DeleteCourse).Methods("POST")
	myRouter.Handle("/getteacherinfo",teachercontroller.GetTeacherInfo).Methods("POST")
	myRouter.Handle("/editteacherinfo",teachercontroller.EditTeacherInfo).Methods("POST")
	myRouter.Handle("/getdescription",coursecontroller.GetDescription).Methods("POST")
	myRouter.Handle("/editdescription",coursecontroller.EditDescription).Methods("POST")
	myRouter.Handle("/getannouncement",coursecontroller.GetAnnouncement).Methods("POST")
	myRouter.Handle("/editannouncement",coursecontroller.EditAnnouncement).Methods("POST")
	myRouter.Handle("/addstudent",coursemembercontroller.AddStudentToCourse).Methods("POST")
	myRouter.Handle("/addteacher",coursemembercontroller.AddTeacherToCourse).Methods("POST")
	myRouter.Handle("/getstudentincourse",coursemembercontroller.GetStudentInCourse).Methods("POST")
	myRouter.Handle("/getteacherincourse",coursemembercontroller.GetTeacherInCourse).Methods("POST")
	myRouter.Handle("/deleteteacherincourse",coursemembercontroller.DeleteTeacherInCourse).Methods("POST")
	myRouter.Handle("/deletestudentincourse",coursemembercontroller.DeleteStudentInCourse).Methods("POST")
	myRouter.Handle("/getuserrole",coursemembercontroller.GetUserRole).Methods("POST")
	myRouter.Handle("/createtest",testcontroller.CreateTest).Methods("POST")
	myRouter.Handle("/login", login.Login).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000",c.Handler(myRouter)))
}

func main() {
	handleRequests()
}