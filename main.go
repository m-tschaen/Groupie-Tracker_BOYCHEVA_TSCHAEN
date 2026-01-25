package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"
	"unicode"
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

type RelationItem struct {
	Place string
	Dates []string
}

type ArtistPageData struct {
	Artist        Artist
	Locations     []string
	Dates         []string
	RelationItems []RelationItem
}

func formatDate(s string) string {
	return strings.TrimPrefix(s, "*")
}

func titleWords(s string) string {
	s = strings.ToLower(s)
	words := strings.Fields(s)
	for i, w := range words {
		r := []rune(w)
		if len(r) == 0 {
			continue
		}
		r[0] = unicode.ToUpper(r[0])
		words[i] = string(r)
	}
	return strings.Join(words, " ")
}

func formatPlace(raw string) string {
	raw = strings.ReplaceAll(raw, "_", " ")
	parts := strings.Split(raw, "-")
	if len(parts) < 2 {
		return titleWords(raw)
	}

	countryRaw := strings.TrimSpace(parts[len(parts)-1])
	cityRaw := strings.TrimSpace(strings.Join(parts[:len(parts)-1], " "))

	city := titleWords(cityRaw)

	countryClean := strings.ReplaceAll(countryRaw, "_", " ")
	countryWords := strings.Fields(countryClean)

	if len(countryWords) == 1 && len(countryWords[0]) <= 3 {
		return city + ", " + strings.ToUpper(countryWords[0])
	}

	return city + ", " + titleWords(countryClean)
}

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
		}).ParseFiles("templates/index.html", "templates/artist.html"),
	)

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

		query := strings.ToLower(r.URL.Query().Get("query"))
		if query != "" {
			var filtered []Artist
			for _, artist := range artists {
				match := strings.Contains(strings.ToLower(artist.Name), query)
				for _, member := range artist.Members {
					if strings.Contains(strings.ToLower(member), query) {
						match = true
					}
				}
                if strings.Contains(fmt.Sprint(artist.CreationDate), query) || strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
					match = true
				}
				if match {
					filtered = append(filtered, artist)
				}
			}
			artists = filtered
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

		keys := make([]string, 0, len(relRes.DatesLocation))
		for k := range relRes.DatesLocation {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		items := make([]RelationItem, 0, len(keys))
		for _, k := range keys {
			dates := relRes.DatesLocation[k]
			for i := range dates {
				dates[i] = formatDate(dates[i])
			}

			items = append(items, RelationItem{
				Place: formatPlace(k),
				Dates: dates,
			})
		}

		formattedLocations := make([]string, 0, len(locRes.Locations))
		for _, l := range locRes.Locations {
			formattedLocations = append(formattedLocations, formatPlace(l))
		}

		data := ArtistPageData{
			Artist:        *found,
			Locations:     formattedLocations,
			Dates:         dateRes.Dates,
			RelationItems: items,
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
