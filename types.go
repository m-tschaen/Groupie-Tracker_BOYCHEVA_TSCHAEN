package main

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
	Favorites     map[int]bool
}

type LocationMarker struct {
	City         string
	Country      string
	Count        int
	Artists      []string
	ArtistIDs    []int
	ArtistDates  [][]string
	X            int
	Y            int
}

type Continent struct {
	Name          string
	Emoji         string
	TotalConcerts int
	Locations     []LocationMarker
}

type LocationsPageData struct {
	Continents []Continent
}

type IndexPageData struct {
	Artists         []Artist
	Query           string
	MinYear         string
	MaxYear         string
	MemberOptions   []int
	SelectedMembers map[int]bool
	Favorites       map[int]bool
}