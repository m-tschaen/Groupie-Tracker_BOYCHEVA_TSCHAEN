package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
	Members      []string `json:"members"`
}

type LocationsResponse struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type DatesResponse struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type RelationsResponse struct {
	Id            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

type ArtistPageData struct {
	Artist    Artist
	Locations []string
	Dates     []string
	Relations map[string][]string
}

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/artist.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		artists, err := fetchArtists()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, "index.html", artists)
	})

	http.HandleFunc("/artist/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/artist/"):]
		if idStr == "" {
			http.NotFound(w, r)
			return
		}

		artists, err := fetchArtists()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var found *Artist
		for i := range artists {
			if fmt.Sprint(artists[i].Id) == idStr {
				found = &artists[i]
				break
			}
		}
		if found == nil {
			http.NotFound(w, r)
			return
		}

		var locRes LocationsResponse
		var dateRes DatesResponse
		var relRes RelationsResponse

		if err := fetchJSON(found.Locations, &locRes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := fetchJSON(found.ConcertDates, &dateRes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := fetchJSON(found.Relations, &relRes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := ArtistPageData{
			Artist:    *found,
			Locations: locRes.Locations,
			Dates:     dateRes.Dates,
			Relations: relRes.DatesLocation,
		}

		tmpl.ExecuteTemplate(w, "artist.html", data)
	})

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func fetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}
	return artists, nil
}

func fetchJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
