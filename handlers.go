package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func homeHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		artists, err := fetchArtists()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		queryRaw := strings.TrimSpace(r.URL.Query().Get("query"))
		query := strings.ToLower(queryRaw)

		minYearStr := r.URL.Query().Get("minYear")
		maxYearStr := r.URL.Query().Get("maxYear")

		minYear, hasMin := parseInt(minYearStr)
		maxYear, hasMax := parseInt(maxYearStr)

		membersVals := r.URL.Query()["members"]
		selected := selectedMembersSet(membersVals)

		filtered := make([]Artist, 0, len(artists))
		for _, artist := range artists {
			if query != "" {
				match := strings.Contains(strings.ToLower(artist.Name), query)

				if !match {
					for _, member := range artist.Members {
						if strings.Contains(strings.ToLower(member), query) {
							match = true
							break
						}
					}
				}

				if !match {
					if strings.Contains(strconv.Itoa(artist.CreationDate), query) || strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
						match = true
					}
				}
				if !match {
					continue
				}
			}
			if hasMin && artist.CreationDate < minYear {
				continue
			}
			if hasMax && artist.CreationDate > maxYear {
				continue
			}
			if len(selected) > 0 {
				if !selected[len(artist.Members)] {
					continue
				}
			}
			filtered = append(filtered, artist)
		}
		
		opts := []int{1, 2, 3, 4, 5, 6, 7, 8}
		data := IndexPageData{
			Artists:         filtered,
			Query:           queryRaw,
			MinYear:         minYearStr,
			MaxYear:         maxYearStr,
			MemberOptions:   opts,
			SelectedMembers: selected,
		}

		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func artistHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func locationsHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artists, _ := fetchArtists()
		locationMap := make(map[string]*LocationMarker)

		for _, artist := range artists {
			var relRes RelationsResponse
			if err := fetchJSON(artist.Relations, &relRes); err == nil {
				for rawLoc, dates := range relRes.DatesLocation {
					cleanLoc := strings.ReplaceAll(rawLoc, "_", " ")
					formatted := formatPlace(rawLoc)
					parts := strings.Split(cleanLoc, "-")
					country := strings.TrimSpace(parts[len(parts)-1])
					city := strings.TrimSpace(strings.Join(parts[:len(parts)-1], " "))

					formattedDates := make([]string, len(dates))
					for i, d := range dates {
						formattedDates[i] = formatDate(d)
					}

					if marker, exists := locationMap[formatted]; exists {
						marker.Count++
						marker.Artists = append(marker.Artists, artist.Name)
						marker.ArtistIDs = append(marker.ArtistIDs, artist.Id)
						marker.ArtistDates = append(marker.ArtistDates, formattedDates)
					} else {
						x, y := getCountryCoordinates(country)
						locationMap[formatted] = &LocationMarker{
							City: titleWords(city), Country: strings.ToUpper(country),
							Count: 1, Artists: []string{artist.Name}, ArtistIDs: []int{artist.Id},
							ArtistDates: [][]string{formattedDates}, X: x, Y: y,
						}
					}
				}
			}
		}

		continentMap := make(map[string]*Continent)
		continentEmojis := map[string]string{
			"Amérique du Nord": "",
			"Amérique du Sud":  "",
			"Europe":           "",
			"Asie":             "",
			"Océanie":          "",
		}

		for _, marker := range locationMap {
			continentName := getContinent(marker.Country)
			if continent, exists := continentMap[continentName]; exists {
				continent.Locations = append(continent.Locations, *marker)
				continent.TotalConcerts += marker.Count
			} else {
				emoji := continentEmojis[continentName]
				continentMap[continentName] = &Continent{
					Name:          continentName,
					Emoji:         emoji,
					TotalConcerts: marker.Count,
					Locations:     []LocationMarker{*marker},
				}
			}
		}

		continents := make([]Continent, 0, len(continentMap))
		for _, c := range continentMap {
			sort.Slice(c.Locations, func(i, j int) bool {
				return c.Locations[i].Count > c.Locations[j].Count
			})
			continents = append(continents, *c)
		}

		sort.Slice(continents, func(i, j int) bool {
			return continents[i].TotalConcerts > continents[j].TotalConcerts
		})

		tmpl.ExecuteTemplate(w, "locations.html", LocationsPageData{Continents: continents})
	}
}