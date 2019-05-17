package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// DBConfig holds the config values for the databases
type DBConfig struct {
	Hostname string
	Name     string // db name
	Password string
	Port     string
	Username string
}

// IndexPage holds the displayed page information
type IndexPage struct {
	Message string
	Class   string
}

func main() {

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", PageHandler)

	fmt.Println(http.ListenAndServe(":8080", nil))
}

// PageHandler queries the producer and displays the quote
func PageHandler(w http.ResponseWriter, r *http.Request) {

	config := DBConfig{
		Hostname: os.Getenv("hostname"),
		Name:     os.Getenv("name"),
		Password: os.Getenv("password"),
		Port:     os.Getenv("port"),
		Username: os.Getenv("username"),
	}

	connectionString := config.Username + ":" + config.Password + "@tcp(" + config.Hostname + ":" + config.Port + ")/" + config.Name

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	template := template.Must(template.ParseFiles("templates/index.html"))

	indexPage := IndexPage{
		Message: "Successfully monitoring the database",
		Class:   "alert-success",
	}

	err = db.Ping()
	if err != nil {
		indexPage = IndexPage{
			Message: "Unable to contact the database",
			Class:   "alert-danger",
		}
	}
	if err := template.ExecuteTemplate(w, "index.html", indexPage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
