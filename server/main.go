package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type TemplateData struct{
    Message string
}

func frontpage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(
		"templates/base.gtpl",
		"templates/navbar.gtpl",
		"templates/frontpage.gtpl",
	)
	t.Execute(w, nil)
}

type editHandler struct{
}

func (h editHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
    fmt.Println(w, "hello")
}

func main() {
    // Setup in memory session store for user tokens
    sessionStore = make(map[string]Client)


    // Setup static router
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

    // Setup non-dynamic routes
	http.HandleFunc("/", frontpage)
	http.HandleFunc("/login", ProcessLogin)
    http.Handle("/edit", Authenticate(editHandler{}))

    fmt.Println("listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer", err)
	}
}
