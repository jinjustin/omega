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
	//"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"

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
	myRouter := mux.NewRouter()

	serve := gin.Default()

	myRouter.HandleFunc("/",testAPI) //work
	myRouter.HandleFunc("/test",test2).Methods("POST") //work
	myRouter.Handle("/createcourse",coursecontroller.CreateCourse).Methods("POST") //work
	myRouter.Handle("/getcourselist",middleware(coursecontroller.GetCourseList)).Methods("POST") //work
	myRouter.Handle("/deletecourse",coursecontroller.DeleteCourse).Methods("POST") //work
	myRouter.Handle("/getteacherinfo",teachercontroller.GetTeacherInfo).Methods("POST") //work
	myRouter.Handle("/editteacherinfo",teachercontroller.EditTeacherInfo).Methods("POST") //work
	myRouter.Handle("/getdescription",coursecontroller.GetDescription).Methods("POST") //bug
	myRouter.Handle("/editdescription",coursecontroller.EditDescription).Methods("POST") //work
	myRouter.Handle("/getannouncement",coursecontroller.GetAnnouncement).Methods("POST")//work
	myRouter.Handle("/editannouncement",coursecontroller.EditAnnouncement).Methods("POST")//work
	myRouter.Handle("/addstudent",coursemembercontroller.AddStudentToCourse).Methods("POST") // work
	myRouter.Handle("/addteacher",coursemembercontroller.AddTeacherToCourse).Methods("POST") //almost-work (Bug)
	myRouter.Handle("/getstudentincourse",coursemembercontroller.GetStudentInCourse).Methods("POST") //work
	myRouter.Handle("/getteacherincourse",coursemembercontroller.GetTeacherInCourse).Methods("POST") //work
	myRouter.Handle("/deleteteacherincourse",coursemembercontroller.DeleteTeacherInCourse).Methods("POST") //work
	myRouter.Handle("/deletestudentincourse",coursemembercontroller.DeleteStudentInCourse).Methods("POST") //work
	myRouter.Handle("/getuserrole",coursemembercontroller.GetUserRole).Methods("POST") //work
	myRouter.Handle("/createtest",testcontroller.CreateTest).Methods("POST")
	myRouter.Handle("/gettestlist",testcontroller.GetTestList).Methods("POST")
	myRouter.Handle("/gettestinfo",testcontroller.GetTestInfo).Methods("POST")
	myRouter.Handle("edittestinfo",testcontroller.EditTestInfo).Methods("POST")
	myRouter.Handle("/login", login.Login).Methods("POST") //work

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET","POST","DELETE"},
		AllowCredentials: true,
		AllowedHeaders: []string{"*"},
	})

	serve.Use(static.Serve("/", static.LocalFile("./build", true)))
	go http.ListenAndServe(":10000", c.Handler(myRouter))
	go http.ListenAndServe(":5000", c.Handler(serve))
	select{}
}

func main() {
	handleRequests()
}