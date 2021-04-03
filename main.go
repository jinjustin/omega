package main

import (
	"fmt"
	"net/http"
	"github.com/jinjustin/omega/coursecontroller"
	"github.com/jinjustin/omega/coursemembercontroller"
	"github.com/jinjustin/omega/login"
	"github.com/jinjustin/omega/teachercontroller"
	"github.com/jinjustin/omega/studentcontroller"
	"github.com/jinjustin/omega/testcontroller"
	"github.com/jinjustin/omega/authentication"
	"github.com/jinjustin/omega/questiongroupcontroller"
	"github.com/jinjustin/omega/storage"
	"github.com/jinjustin/omega/questioncontroller"
	"github.com/jinjustin/omega/answercontroller"

	//"omega/database"
	"context"
	//"encoding/json"
	"io/ioutil"
	//"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/go-co-op/gocron"
	"time"
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

func middlewareAdmin(next http.Handler) http.Handler {
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
				if role == "admin"{
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

func executeCronJob() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minute().Do(testcontroller.UpdateTestSituation())
}

func handleRequests() {
	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/",testAPI) 
	myRouter.HandleFunc("/test",test2).Methods("POST") 
	myRouter.Handle("/createcourse",middlewareTeacher(coursecontroller.CreateCourse)).Methods("POST") 
	myRouter.Handle("/getcourselist",middleware(coursecontroller.GetCourseList)).Methods("POST") 
	myRouter.Handle("/deletecourse",middlewareTeacher(coursecontroller.DeleteCourse)).Methods("POST")

	myRouter.Handle("/getteacherinfo",middlewareTeacher(teachercontroller.GetTeacherInfo)).Methods("POST") 
	myRouter.Handle("/editteacherinfo",middlewareTeacher(teachercontroller.EditTeacherInfo)).Methods("POST")
	
	myRouter.Handle("/getstudentinfo",middlewareStudent(studentcontroller.GetStudentInfo)).Methods("GET")
	myRouter.Handle("/getstudentinfo",middlewareStudent(studentcontroller.EditTeacherInfo)).Methods("POST")

	myRouter.Handle("/getdescription",middlewareTeacher(coursecontroller.GetDescription)).Methods("POST") 
	myRouter.Handle("/editdescription",middlewareTeacher(coursecontroller.EditDescription)).Methods("POST") 
	myRouter.Handle("/getannouncement",middlewareTeacher(coursecontroller.GetAnnouncement)).Methods("POST")
	myRouter.Handle("/editannouncement",middlewareTeacher(coursecontroller.EditAnnouncement)).Methods("POST")

	//coursemembercontroller
	myRouter.Handle("/addstudent",middlewareTeacher(coursemembercontroller.AddStudentToCourse)).Methods("POST") 
	myRouter.Handle("/addstudentbyfile",coursemembercontroller.AddStudentFromExcelFile).Methods("POST") 
	myRouter.Handle("/addteacher",middlewareTeacher(coursemembercontroller.AddTeacherToCourse)).Methods("POST") 
	myRouter.Handle("/acceptjoincourse",coursemembercontroller.StudentAcceptJoinCourse).Methods("GET")
	myRouter.Handle("/getstudentincourse",coursemembercontroller.GetStudentInCourse).Methods("POST") 
	myRouter.Handle("/getteacherincourse",coursemembercontroller.GetTeacherInCourse).Methods("POST") 
	myRouter.Handle("/deleteteacherincourse",middlewareTeacher(coursemembercontroller.DeleteTeacherInCourse)).Methods("POST") 
	myRouter.Handle("/deletestudentincourse",middlewareTeacher(coursemembercontroller.DeleteStudentInCourse)).Methods("POST")
	myRouter.Handle("/teacheraddcourse",middlewareTeacher(coursemembercontroller.TeacherAddCourse)).Methods("POST")
	
	myRouter.Handle("/getrole",authentication.GetRole).Methods("POST") 
	myRouter.Handle("/deletetest",middlewareTeacher(testcontroller.DeleteTest)).Methods("POST")
	myRouter.Handle("/login", login.Login).Methods("POST") 
	myRouter.Handle("/getusername", authentication.GetUsers).Methods("POST")
	
	myRouter.Handle("/grouptestlistupdate",middlewareTeacher(questiongroupcontroller.GroupTestListUpdate)).Methods("POST")
	myRouter.Handle("/grouptestlistupdate",middleware(questiongroupcontroller.GetGroupInTest)).Methods("GET")

	myRouter.Handle("/testbankupdate",middlewareTeacher(questiongroupcontroller.TestbankUpdate)).Methods("POST")
	myRouter.Handle("/testbankupdate",middlewareTeacher(questiongroupcontroller.GetGroupInTestbank)).Methods("GET")

	myRouter.Handle("/allgrouptestlist",middlewareTeacher(questiongroupcontroller.AllGroupTestList)).Methods("GET")
	myRouter.Handle("/allgrouptestlist",middlewareTeacher(questiongroupcontroller.AllGroupTestListPost)).Methods("POST")
	myRouter.Handle("/getallheaderintest",middlewareTeacher(questiongroupcontroller.GetAllHeaderInTest)).Methods("GET")

	myRouter.Handle("/updatedetailtest", middlewareTeacher(testcontroller.PostDetailTest)).Methods("POST")
	myRouter.Handle("/updatedetailtest", middlewareTeacher(testcontroller.GetDetailTest)).Methods("GET")
	myRouter.Handle("/deletetest", middlewareTeacher(testcontroller.DeleteTest)).Methods("DELETE")
	myRouter.Handle("/changedraftstatus", middlewareTeacher(testcontroller.ChangeDraftStatus)).Methods("POST")
	myRouter.Handle("/changedraftstatus", middlewareTeacher(testcontroller.GetDraftStatus)).Methods("GET")
	myRouter.Handle("/getalltestincourse", middlewareTeacher(testcontroller.GetAllTestInCourse)).Methods("GET")
	myRouter.Handle("/studentgettestlist", middlewareStudent(testcontroller.StudentGetTestListByDay)).Methods("GET")
	myRouter.Handle("/getallfinishedtest", middlewareTeacher(testcontroller.GetAllFinishTestInCourse)).Methods("GET")

	myRouter.Handle("/addnewquestion", middlewareTeacher(questioncontroller.UpdateQuestion)).Methods("POST")
	myRouter.Handle("/getquestion", middlewareTeacher(questioncontroller.GetQuestion)).Methods("GET")
	//myRouter.Handle("/deletequestion", middlewareTeacher(questioncontroller.DeleteQuestion)).Methods("DELETE")
	myRouter.Handle("/getallquestioningroup", middlewareTeacher(questioncontroller.GetAllQuestionInGroup)).Methods("GET")
	myRouter.Handle("/inputexam", middlewareStudent(questioncontroller.InputExam)).Methods("POST")

	myRouter.Handle("/updateallquestionintest", middlewareTeacher(questioncontroller.UpdateAllQuestionInTest)).Methods("POST")
	myRouter.Handle("/updateallquestionintest", middlewareTeacher(questioncontroller.GetAllQuestionInTest)).Methods("GET")
	myRouter.Handle("/getallquestionfortest", middlewareStudent(questioncontroller.GetAllQuestionForTest)).Methods("GET")

	myRouter.Handle("/getstudentcourselist", middlewareTeacher(coursecontroller.GetStudentCourse)).Methods("GET")

	myRouter.Handle("/uploadpic", storage.UploadPic).Methods("POST")

	myRouter.Handle("/changestudentpassword",middlewareStudent(coursemembercontroller.ChangeStudentPassword)).Methods("POST")

	myRouter.Handle("/submitanswer", answercontroller.SubmitAnswer).Methods("POST")
	myRouter.Handle("/getstudentAnswer", answercontroller.GetAnswer).Methods("GET")
	myRouter.Handle("/scoreAnswer", answercontroller.ScoringAnswer).Methods("POST")
	myRouter.Handle("/getallstudentanswerinformation", answercontroller.GetAllStudentAnswerInformation).Methods("GET")
	myRouter.Handle("/getstatisticvalue", answercontroller.GetStatisticValue).Methods("GET")
	myRouter.Handle("/studentgetscore",middlewareStudent(answercontroller.StudentGetScore)).Methods("GET")
	myRouter.Handle("/getallstudentscore", answercontroller.GetAllStudentScore).Methods("GET")

	myRouter.Handle("/addteachertosystem", coursemembercontroller.AddTeacherToSystem).Methods("POST")

	myRouter.Handle("/updatetestsituation", testcontroller.TestUpdateSituation).Methods("GET")
	
	//myRouter.Handle("/testautoscoring", answercontroller.TestAutoScoring).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET","POST","DELETE"},
		AllowCredentials: true,
		AllowedHeaders: []string{"*"},
	})

	http.ListenAndServe(":10000", c.Handler(myRouter))
}

func main() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Second().Do(testcontroller.UpdateTestSituation())
	go s.StartAsync()
	handleRequests()
}