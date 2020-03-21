package main

import (
	"database/sql"
	"net/http"
	"html/template"
	"log"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

type GuestbookIndexData struct {
	SigneesCount int
}

func main() {
	router := mux.NewRouter()
	username := os.Getenv("DB_USER")
	if username == "" {
		username = "root"
	}
	
	password := os.Getenv("DB_PASSWORD")
	db, err := connectToDatabase(username, password, "127.0.0.1", "3306", "guestbook")

	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Serving /")

		count := 0
		err := db.QueryRow("SELECT COUNT(id) AS count FROM messages").Scan(&count)

		if err != nil {
			log.Fatal(err)	
		}

		data := GuestbookIndexData {
			SigneesCount: count,
		}

		tpl := template.Must(template.ParseFiles("views/index.html"))
		tpl.Execute(w, data)
	})

	fs := http.FileServer(http.Dir("static/"))
	log.Print("Serving static assets at /static")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	log.Print("Starting server at port 3000")
	http.ListenAndServe(":3000", router)
}

func connectToDatabase(username string, password string, host string, port string, dbname string) (*sql.DB, error) {
	dns := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", username, password, host, port, dbname)
	return sql.Open("mysql", dns)
}