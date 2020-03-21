package main

import (
	"net/http"
	"html/template"
	"log"

	"github.com/gorilla/mux"
)

type GuestbookIndexData struct {
	SigneesCount int
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Serving /")

		tpl := template.Must(template.ParseFiles("views/index.html"))
		data := GuestbookIndexData {
			SigneesCount: 0,
		}

		tpl.Execute(w, data)
	})

	fs := http.FileServer(http.Dir("static/"))
	log.Print("Serving static assets at /static")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	log.Print("Starting server at port 3000")
	http.ListenAndServe(":3000", router)
}