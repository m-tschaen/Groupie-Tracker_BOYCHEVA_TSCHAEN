package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

var (
	cachedArtists []Artist
	cacheMutex    sync.RWMutex
	lastFetch     time.Time
	cacheDuration = 10 * time.Minute
)

func fetchArtists() ([]Artist, error) {
	cacheMutex.RLock()
	if time.Since(lastFetch) < cacheDuration && len(cachedArtists) > 0 {
		artists := cachedArtists
		cacheMutex.RUnlock()
		return artists, nil
	}
	cacheMutex.RUnlock()

	return refreshCache()
}

func refreshCache() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		cacheMutex.RLock()
		defer cacheMutex.RUnlock()
		if len(cachedArtists) > 0 {
			return cachedArtists, nil
		}
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		cacheMutex.RLock()
		defer cacheMutex.RUnlock()
		if len(cachedArtists) > 0 {
			return cachedArtists, nil
		}
		return nil, err
	}

	cacheMutex.Lock()
	cachedArtists = artists
	lastFetch = time.Now()
	cacheMutex.Unlock()

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