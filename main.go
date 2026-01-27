package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	tmpl := template.Must(
		template.New("root").Funcs(template.FuncMap{
			"formatPlace": formatPlace,
			"formatDate":  formatDate,
		}).ParseFiles("templates/index.html", "templates/artist.html", "templates/locations.html"),
	)

	http.HandleFunc("/", homeHandler(tmpl))
	http.HandleFunc("/artist/", artistHandler(tmpl))
	http.HandleFunc("/locations", locationsHandler(tmpl))

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}