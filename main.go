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
		}).ParseFiles("templates/welcome.html", "templates/index.html", "templates/artist.html", "templates/locations.html", "templates/favorites.html"),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "welcome.html", nil)
    })
	http.HandleFunc("/tracker", homeHandler(tmpl))
	http.HandleFunc("/artist/", artistHandler(tmpl))
	http.HandleFunc("/locations", locationsHandler(tmpl))
	http.HandleFunc("/favorite", favoriteHandler)
	http.HandleFunc("/favorites", favoritesPageHandler(tmpl))

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}