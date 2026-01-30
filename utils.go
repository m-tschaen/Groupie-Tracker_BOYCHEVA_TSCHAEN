package main

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
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
	country := titleWords(countryRaw)

	if len(strings.ReplaceAll(countryRaw, " ", "")) <= 3 {
		country = strings.ToUpper(countryRaw)
	}

	return city + ", " + country
}

func parseInt(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}

func selectedMembersSet(vals []string) map[int]bool {
	set := make(map[int]bool)
	for _, v := range vals {
		if n, ok := parseInt(v); ok {
			set[n] = true
		}
	}
	return set
}

func readFavoritesCookie(r *http.Request) map[int]bool {
	favs := map[int]bool{}
	c, err := r.Cookie("favorites")
	if err != nil || strings.TrimSpace(c.Value) == "" {
		return favs
	}
	parts := strings.Split(c.Value, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if id, ok := parseInt(p); ok {
			favs[id] = true
		}
	}
	return favs
}

func writeFavoritesCookie(w http.ResponseWriter, favs map[int]bool) {
	values := make([]string, 0, len(favs))
	for id := range favs {
		values = append(values, strconv.Itoa(id))
	}
	sort.Strings(values)
	http.SetCookie(w, &http.Cookie{
		Name:  "favorites",
		Value: strings.Join(values, ","),
		Path:  "/",
	})
}