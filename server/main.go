package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func frontpage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(
		"templates/base.gtpl",
		"templates/navbar.gtpl",
		"templates/frontpage.gtpl",
	)
	t.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles(
			"templates/base.gtpl",
			"templates/navbar.gtpl",
			"templates/login.gtpl",
		)

		t.Execute(w, nil)
	} else {
		r.ParseForm()
		fmt.Println("Parsing login for: ", r.Form["username"])

	}
}

func main() {
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", frontpage)
	http.HandleFunc("/login", login)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer", err)
	}
}
