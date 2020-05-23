package main

import (
	"fmt"
	"net/http"
	"text/template"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type User struct{
	Username string 
}
func logResponseWrite(w http.ResponseWriter, r *http.Request){
	fmt.Println("********************")
	fmt.Println("response:",w)
	fmt.Println("--------------------")
	fmt.Println("request:",r)
	fmt.Println("********************")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	
	tpl, err := template.ParseFiles("templates/index.gohtml")
	err = tpl.Execute(w, nil)
	checkErr(err)
	logResponseWrite(w,r)
}

func dbConn() (db *sql.DB) {
    dbDriver 	:= "mysql"
	dbName 		:= "go"
	dbUser		:= "username"
	dbPass 		:= "password"


	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	fmt.Println(db)
    if err != nil {
        panic(err.Error())
	}
	fmt.Println("db connection successfully")
    return db
}



func getUsername(w http.ResponseWriter, r *http.Request){

	usernamesearch := r.FormValue("username")
	db := dbConn()

	var user User
	res := []User{}
	results,_:= db.Query("SELECT username FROM users WHERE username = ?", usernamesearch)

	for results.Next() {
		_err := results.Scan(&user.Username)
		res = append(res, user)
		checkErr(_err)
	}
	
	
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}



func main(){
	fmt.Println("http://localhost:8080 starting...")
	fs := http.FileServer(http.Dir("./static"))
  	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/",index)
	http.HandleFunc("/getUsername",getUsername)
	http.ListenAndServe(":8080", nil)
}