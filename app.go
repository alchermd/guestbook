package main

import (
	"database/sql"
	"net/http"
	"html/template"
	"log"
	"fmt"
	"time"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/go-sql-driver/mysql"
)

type GuestbookIndexData struct {
	SigneesCount int
	Flashes []interface{}
}

type GuestbookMessagesData struct {
	Messages []Message
}

type Message struct {
	Id int
	Name string
	Message string
	CreatedAt time.Time	
}

func main() {
	username := os.Getenv("DB_USER")
	if username == "" {
		username = "root"
	}
	
	password := os.Getenv("DB_PASSWORD")
	db, err := connectToDatabase(username, password, "127.0.0.1", "3306", "guestbook")

	if err != nil {
		log.Fatal(err)
	}

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Serving /")

		count := 0
		err := db.QueryRow("SELECT COUNT(id) AS count FROM messages").Scan(&count)

		if err != nil {
			log.Fatal(err)	
		}

		session, err := store.Get(r, "session-name")
	    if err != nil {
	        log.Fatal(err)
	    }

	    flashes := session.Flashes()
	    session.Save(r, w)

		data := GuestbookIndexData {
			SigneesCount: count,
			Flashes: flashes,
		}

		tpl := template.Must(template.ParseFiles("views/index.html"))
		tpl.Execute(w, data)
	}).Methods("GET")

	router.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Serving /messages")

		rows, err := db.Query("SELECT id, name, message, created_at FROM messages")

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		var messages []Message
		for rows.Next() {
			var m Message
			err := rows.Scan(&m.Id, &m.Name, &m.Message, &m.CreatedAt)

			if err != nil {
				log.Fatal(err)
			}

			messages = append(messages, m)
		}

		data := GuestbookMessagesData {
			Messages: messages,
		}

		tpl := template.Must(template.ParseFiles("views/messages.html"))
		tpl.Execute(w, data)
	}).Methods("GET")


	router.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Serving /messages [POST]")

		name := r.FormValue("name")

		if name == "" {
			name = "Anonymous"
		}

		_, err := db.Exec("INSERT INTO messages(name, message) VALUES(?, ?)", name, r.FormValue("message"))

		if err != nil {
			log.Fatal(err)
		}

		session, err := store.Get(r, "session-name")
	    if err != nil {
	        log.Fatal(err)
	    }

	    session.AddFlash("Message successfully saved.")
	    err = session.Save(r, w)

	    if err != nil {
	    	log.Fatal(err)
	    }

		http.Redirect(w, r, "/", 301)
	}).Methods("POST")

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