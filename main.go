package main

import (
	"fmt"
	"net/http"
	"github.com/jinjustin/omega/coursecontroller"
	"github.com/jinjustin/omega/coursemembercontroller"
	"github.com/jinjustin/omega/login"
	"github.com/jinjustin/omega/teachercontroller"
	"github.com/jinjustin/omega/testcontroller"
	"github.com/jinjustin/omega/authentication"
	"github.com/jinjustin/omega/questiongroupcontroller"
	"github.com/jinjustin/omega/storage"
	"github.com/jinjustin/omega/questioncontroller"

	//"omega/database"
	"context"
	//"encoding/json"
	"io/ioutil"
	//"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	//"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/contrib/static"

	//"omega/student"
	//"database/sql"
	//_ "github.com/lib/pq"
)

func testAPI(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "omega")
}

func test2(w http.ResponseWriter, r *http.Request){

	reqBody, _ := ioutil.ReadAll(r.Body)
	w.Write(reqBody)
}

func middlewareTeacher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		role := authentication.GetUserRole(r)
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
				if role == "teacher"{
					ctx := context.WithValue(r.Context(), "props", claims)
					// Access context values in handlers like this
					// props, _ := r.Context().Value("props").(jwt.MapClaims)
					next.ServeHTTP(w, r.WithContext(ctx))
				}else {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("You don't have permission to access this"))
				}
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}
		}
	})
}

func middlewareStudent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		role := authentication.GetUserRole(r)
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
				if role == "student"{
					ctx := context.WithValue(r.Context(), "props", claims)
					// Access context values in handlers like this
					// props, _ := r.Context().Value("props").(jwt.MapClaims)
					next.ServeHTTP(w, r.WithContext(ctx))
				}else {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("You don't have permission to access this"))
				}
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}
		}
	})
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

	myRouter.HandleFunc("/",testAPI) //work
	myRouter.HandleFunc("/test",test2).Methods("POST") //work
	myRouter.Handle("/createcourse",middlewareTeacher(coursecontroller.CreateCourse)).Methods("POST") //work
	myRouter.Handle("/getcourselist",middleware(coursecontroller.GetCourseList)).Methods("POST") //work
	myRouter.Handle("/deletecourse",middlewareTeacher(coursecontroller.DeleteCourse)).Methods("POST") //work
	myRouter.Handle("/getteacherinfo",middlewareTeacher(teachercontroller.GetTeacherInfo)).Methods("POST") //work
	myRouter.Handle("/editteacherinfo",middlewareTeacher(teachercontroller.EditTeacherInfo)).Methods("POST") //work
	myRouter.Handle("/getdescription",middlewareTeacher(coursecontroller.GetDescription)).Methods("POST") //bug
	myRouter.Handle("/editdescription",middlewareTeacher(coursecontroller.EditDescription)).Methods("POST") //work
	myRouter.Handle("/getannouncement",middlewareTeacher(coursecontroller.GetAnnouncement)).Methods("POST")//work
	myRouter.Handle("/editannouncement",middlewareTeacher(coursecontroller.EditAnnouncement)).Methods("POST")//work
	myRouter.Handle("/addstudent",middlewareTeacher(coursemembercontroller.AddStudentToCourse)).Methods("POST") // work
	myRouter.Handle("/addstudentbyfile",coursemembercontroller.AddStudentFromExcelFile).Methods("POST") 
	myRouter.Handle("/addteacher",middlewareTeacher(coursemembercontroller.AddTeacherToCourse)).Methods("POST") //almost-work (Bug)
	myRouter.Handle("/acceptjoincourse",coursemembercontroller.StudentAcceptJoinCourse).Methods("GET")
	myRouter.Handle("/getstudentincourse",coursemembercontroller.GetStudentInCourse).Methods("POST") //work
	myRouter.Handle("/getteacherincourse",coursemembercontroller.GetTeacherInCourse).Methods("POST") //work
	myRouter.Handle("/deleteteacherincourse",middlewareTeacher(coursemembercontroller.DeleteTeacherInCourse)).Methods("POST") //work
	myRouter.Handle("/deletestudentincourse",middlewareTeacher(coursemembercontroller.DeleteStudentInCourse)).Methods("POST") //work
	myRouter.Handle("/getrole",authentication.GetRole).Methods("POST") //work
	myRouter.Handle("/deletetest",middlewareTeacher(testcontroller.DeleteTest)).Methods("POST")
	myRouter.Handle("/login", login.Login).Methods("POST") //work
	myRouter.Handle("/getusername", authentication.GetUsers).Methods("POST")
	
	myRouter.Handle("/grouptestlistupdate",middlewareTeacher(questiongroupcontroller.GroupTestListUpdate)).Methods("POST")
	myRouter.Handle("/getgroupintest",middlewareTeacher(questiongroupcontroller.GetGroupInTest)).Methods("GET")
	myRouter.Handle("/testbankupdate",middlewareTeacher(questiongroupcontroller.TestbankUpdate)).Methods("POST")
	myRouter.Handle("/getgroupintestbank",middlewareTeacher(questiongroupcontroller.GetGroupInTestbank)).Methods("GET")
	myRouter.Handle("/allgrouptestlist",middlewareTeacher(questiongroupcontroller.AllGroupTestList)).Methods("GET")
	myRouter.Handle("/postdetailtest", middlewareTeacher(testcontroller.PostDetailTest)).Methods("POST")
	myRouter.Handle("/getdetailtest", middlewareTeacher(testcontroller.GetDetailTest)).Methods("GET")
	myRouter.Handle("/deletetest", middlewareTeacher(testcontroller.DeleteTest)).Methods("POST")
	myRouter.Handle("/changedraftstatus", middlewareTeacher(testcontroller.ChangeDraftStatus)).Methods("POST")

	myRouter.Handle("/addnewquestion", middlewareTeacher(questioncontroller.AddNewQuestion)).Methods("POST")
	myRouter.Handle("/getquestion", middlewareTeacher(questioncontroller.GetQuestion)).Methods("GET")
	myRouter.Handle("/deletequestion", middlewareTeacher(questioncontroller.DeleteQuestion)).Methods("GET")
	myRouter.Handle("/getallquestioningroup", middlewareTeacher(questioncontroller.GetAllQuestionInGroup)).Methods("GET")

	myRouter.Handle("/uploadpic", storage.UploadPic).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET","POST","DELETE"},
		AllowCredentials: true,
		AllowedHeaders: []string{"*"},
	})

	http.ListenAndServe(":10000", c.Handler(myRouter))
}

func main() {
	handleRequests()
}