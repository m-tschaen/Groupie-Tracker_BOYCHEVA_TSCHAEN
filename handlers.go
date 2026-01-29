package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type ErrorData struct {
	Code    int
	Title   string
	Message string
}

func errorHandler(w http.ResponseWriter, tmpl *template.Template, code int, message string) {
	w.WriteHeader(code)
	
	var title string
	var msg string
	
	switch code {
	case 404:
		title = "Page introuvable"
		if message == "" {
			msg = "DÃ©solÃ©, la page que vous recherchez n'existe pas ou a Ã©tÃ© dÃ©placÃ©e."
		} else {
			msg = message
		}
	case 400:
		title = "RequÃªte invalide"
		if message == "" {
			msg = "Les paramÃ¨tres de votre requÃªte sont incorrects."
		} else {
			msg = message
		}
	case 500:
		title = "Erreur serveur"
		if message == "" {
			msg = "Une erreur interne s'est produite. Veuillez rÃ©essayer plus tard."
		} else {
			msg = message
		}
	default:
		title = "Erreur"
		msg = message
	}
	
	data := ErrorData{
		Code:    code,
		Title:   title,
		Message: msg,
	}
	
	if err := tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		http.Error(w, fmt.Sprintf("%d - %s", code, title), code)
	}
}

func homeHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tracker" {
			errorHandler(w, tmpl, 404, "")
			return
		}

		artists, err := fetchArtists()
		if err != nil {
			errorHandler(w, tmpl, 500, "Impossible de récupérer les données.")
			return
		}

		queryRaw := strings.TrimSpace(r.URL.Query().Get("query"))
		query := strings.ToLower(queryRaw)

		minYearStr := r.URL.Query().Get("minYear")
		maxYearStr := r.URL.Query().Get("maxYear")

		minYear, hasMin := parseInt(minYearStr)
		maxYear, hasMax := parseInt(maxYearStr)

		if minYearStr != "" && !hasMin {
			errorHandler(w, tmpl, 400, "L'année minimale doit être un nombre valide.")
			return
		}
		if maxYearStr != "" && !hasMax {
			errorHandler(w, tmpl, 400, "L'année maximale doit être un nombre valide.")
			return
		}
		if hasMin && (minYear < 1900 || minYear > 2100) {
			errorHandler(w, tmpl, 400, "L'année minimale doit être entre 1900 et 2100.")
			return
		}
		if hasMax && (maxYear < 1900 || maxYear > 2100) {
			errorHandler(w, tmpl, 400, "L'année maximale doit être entre 1900 et 2100.")
			return
		}
		if hasMin && hasMax && minYear > maxYear {
			errorHandler(w, tmpl, 400, "L'année minimale ne peut pas être supérieure à l'année maximale.")
			return
		}

		membersVals := r.URL.Query()["members"]
		selected := selectedMembersSet(membersVals)
		
		for memberCount := range selected {
			if memberCount < 1 || memberCount > 10 {
				errorHandler(w, tmpl, 400, "Le nombre de membres doit être entre 1 et 10.")
				return
			}
		}

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
					var relRes RelationsResponse
					if err := fetchJSON(artist.Relations, &relRes); err == nil {
						for location, dates := range relRes.DatesLocation {
							locationClean := strings.ReplaceAll(strings.ToLower(location), "_", " ")
							locationFormatted := strings.ToLower(formatPlace(location))
							
							if strings.Contains(locationClean, query) || strings.Contains(locationFormatted, query) {
								match = true
								break
							}

							for _, date := range dates {
								dateClean := strings.ToLower(formatDate(date))
								if strings.Contains(dateClean, query) {
									match = true
									break
								}
							}
							
							if match {
								break
							}
						}
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
		
		favs := readFavoritesCookie(r)
		opts := []int{1, 2, 3, 4, 5, 6, 7, 8}
		data := IndexPageData{
			Artists:         filtered,
			Query:           queryRaw,
			MinYear:         minYearStr,
			MaxYear:         maxYearStr,
			MemberOptions:   opts,
			SelectedMembers: selected,
			Favorites: favs,
		}

		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			errorHandler(w, tmpl, 500, "Impossible de récupérer les données.")
			return
		}
	}
}

func artistHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/artist/"):]
		if idStr == "" {
			errorHandler(w, tmpl, 404, "")
			return
		}
		
		if _, err := strconv.Atoi(idStr); err != nil {
			errorHandler(w, tmpl, 400, "L'ID de l'artiste doit Ãªtre un nombre valide.")
			return
		}

		artists, err := fetchArtists()
		if err != nil {
			errorHandler(w, tmpl, 500, "Impossible de rÃ©cupÃ©rer les donnÃ©es.")
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
			errorHandler(w, tmpl, 404, "Cet artiste n'existe pas.")
			return
		}

		var locRes LocationsResponse
		var dateRes DatesResponse
		var relRes RelationsResponse

		if err := fetchJSON(found.Locations, &locRes); err != nil {
			errorHandler(w, tmpl, 500, "Impossible de rÃ©cupÃ©rer les donnÃ©es.")
			return
		}
		if err := fetchJSON(found.ConcertDates, &dateRes); err != nil {
			errorHandler(w, tmpl, 500, "Impossible de rÃ©cupÃ©rer les donnÃ©es.")
			return
		}
		if err := fetchJSON(found.Relations, &relRes); err != nil {
			errorHandler(w, tmpl, 500, "Impossible de rÃ©cupÃ©rer les donnÃ©es.")
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

		favs := readFavoritesCookie(r)
		data := ArtistPageData{
			Artist:        *found,
			Locations:     formattedLocations,
			Dates:         dateRes.Dates,
			RelationItems: items,
			Favorites: favs,
		}

		if err := tmpl.ExecuteTemplate(w, "artist.html", data); err != nil {
			errorHandler(w, tmpl, 500, "Erreur lors de l'affichage de la page.")
			return
		}
	}
}

func locationsHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artists, err := fetchArtists()
		if err != nil {
			errorHandler(w, tmpl, 500, "Impossible de rÃ©cupÃ©rer les artistes pour la carte.")
			return
		}
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
			"AmÃ©rique du Nord": "",
			"AmÃ©rique du Sud":  "",
			"Europe":           "",
			"Asie":             "",
			"OcÃ©anie":          "",
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

		if err := tmpl.ExecuteTemplate(w, "locations.html", LocationsPageData{Continents: continents}); err != nil {
			errorHandler(w, tmpl, 500, "Erreur lors de l'affichage de la carte.")
			return
		}
	}
}

func compareHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artists, err := fetchArtists()
		if err != nil {
			errorHandler(w, tmpl, 500, "Impossible de rÃ©cupÃ©rer les donnÃ©es.")
			return
		}

		id1Str := r.URL.Query().Get("artist1")
		id2Str := r.URL.Query().Get("artist2")

		data := ComparePageData{
			AllArtists: artists,
		}

		
		if id1Str != "" && id2Str != "" {
			if _, err := strconv.Atoi(id1Str); err != nil {
				errorHandler(w, tmpl, 400, "L'ID du premier artiste doit Ãªtre un nombre valide.")
				return
			}
			if _, err := strconv.Atoi(id2Str); err != nil {
				errorHandler(w, tmpl, 400, "L'ID du deuxiÃ¨me artiste doit Ãªtre un nombre valide.")
				return
			}

			
			for i := range artists {
				if fmt.Sprint(artists[i].Id) == id1Str {
					data.Artist1 = &artists[i]
					break
				}
			}

			
			for i := range artists {
				if fmt.Sprint(artists[i].Id) == id2Str {
					data.Artist2 = &artists[i]
					break
				}
			}

			
			if data.Artist1 == nil {
				errorHandler(w, tmpl, 404, "Le premier artiste n'existe pas.")
				return
			}
			if data.Artist2 == nil {
				errorHandler(w, tmpl, 404, "Le deuxiÃ¨me artiste n'existe pas.")
				return
			}

			if data.Artist1 != nil {
				var relRes RelationsResponse
				if err := fetchJSON(data.Artist1.Relations, &relRes); err == nil {
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
					data.Relations1 = items
				}
			}

			
			if data.Artist2 != nil {
				var relRes RelationsResponse
				if err := fetchJSON(data.Artist2.Relations, &relRes); err == nil {
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
					data.Relations2 = items
				}
			}
		}

		if err := tmpl.ExecuteTemplate(w, "compare.html", data); err != nil {
			errorHandler(w, tmpl, 500, "Erreur lors de l'affichage de la comparaison.")
			return
		}
	}
}

func favoriteHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, ok := parseInt(idStr)
		if !ok {
			errorHandler(w, tmpl, 400, "L'ID du favori est invalide.")
			return
		}

		
		artists, err := fetchArtists()
		if err != nil {
			errorHandler(w, tmpl, 500, "Impossible de vÃ©rifier l'artiste.")
			return
		}

		found := false
		for _, artist := range artists {
			if artist.Id == id {
				found = true
				break
			}
		}
		if !found {
			errorHandler(w, tmpl, 404, "Cet artiste n'existe pas.")
			return
		}

		favs := readFavoritesCookie(r)

		if favs[id] {
			delete(favs, id)
		} else {
			favs[id] = true
		}

		writeFavoritesCookie(w, favs)
		back := r.URL.Query().Get("back")
		if back == "favorites" {
			http.Redirect(w, r, "/favorites", http.StatusSeeOther)
			return
		}
		if back == "artist" {
			http.Redirect(w, r, "/artist/"+strconv.Itoa(id), http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/tracker", http.StatusSeeOther)
	}
}

func favoritesPageHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/favorites" {
			errorHandler(w, tmpl, 404, "")
			return
		}

		artists, err := fetchArtists()
		if err != nil {
			errorHandler(w, tmpl, 500, "Impossible de rÃ©cupÃ©rer les donnÃ©es.")
			return
		}

		favs := readFavoritesCookie(r)

		onlyFav := make([]Artist, 0)
		for _, a := range artists {
			if favs[a.Id] {
				onlyFav = append(onlyFav, a)
			}
		}

		data := IndexPageData{
			Artists:   onlyFav,
			Favorites: favs,
		}

		if err := tmpl.ExecuteTemplate(w, "favorites.html", data); err != nil {
			errorHandler(w, tmpl, 500, "Impossible de rÃ©cupÃ©rer les donnÃ©es.")
			return
		}
	}
}