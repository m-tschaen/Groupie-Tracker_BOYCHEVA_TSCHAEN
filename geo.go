package main

import (
	"strings"
)

func getCountryCoordinates(country string) (int, int) {
	country = strings.ToLower(strings.TrimSpace(country))
	country = strings.ReplaceAll(country, "_", " ")

	type LatLong struct {
		Lat float64
		Lon float64
	}

	geoCoords := map[string]LatLong{
		"usa":                  {37.0, -95.0},
		"canada":               {56.0, -106.0},
		"mexico":               {23.0, -102.0},
		"brazil":               {-10.0, -55.0},
		"argentina":            {-34.0, -64.0},
		"chile":                {-30.0, -71.0},
		"peru":                 {-9.0, -75.0},
		"colombia":             {4.0, -72.0},
		"costa rica":           {10.0, -84.0},
		"uk":                   {54.0, -2.0},
		"france":               {46.0, 2.0},
		"germany":              {51.0, 10.0},
		"italy":                {42.5, 12.5},
		"spain":                {40.0, -4.0},
		"portugal":             {39.5, -8.0},
		"netherlands":          {52.5, 5.75},
		"belgium":              {50.5, 4.5},
		"ireland":              {53.0, -8.0},
		"switzerland":          {47.0, 8.0},
		"austria":              {47.5, 14.5},
		"czechia":              {49.75, 15.5},
		"poland":               {52.0, 20.0},
		"slovakia":             {48.7, 19.5},
		"hungary":              {47.0, 20.0},
		"sweden":               {62.0, 15.0},
		"norway":               {60.0, 8.0},
		"denmark":              {56.0, 10.0},
		"finland":              {64.0, 26.0},
		"greece":               {39.0, 22.0},
		"romania":              {46.0, 25.0},
		"belarus":              {53.0, 28.0},
		"qatar":                {25.5, 51.25},
		"united arab emirates": {24.0, 54.0},
		"saudi arabia":         {24.0, 45.0},
		"japan":                {36.0, 138.0},
		"china":                {35.0, 105.0},
		"south korea":          {37.0, 127.5},
		"taiwan":               {23.5, 121.0},
		"thailand":             {15.0, 100.0},
		"indonesia":            {-2.0, 118.0},
		"philippines":          {13.0, 122.0},
		"india":                {20.0, 77.0},
		"australia":            {-25.0, 133.0},
		"new zealand":          {-41.0, 174.0},
		"new caledonia":        {-21.5, 165.5},
		"french polynesia":     {-17.5, -149.5},
		"netherlands antilles": {12.2, -69.0},
	}

	var lat, lon float64
	if coord, ok := geoCoords[country]; ok {
		lat = coord.Lat
		lon = coord.Lon
	} else {
		return 600, 300
	}

	x := int((lon + 180.0) * (1200.0 / 360.0))
	y := int((90.0 - lat) * (600.0 / 180.0))

	x = clamp(x, 10, 1190)
	y = clamp(y, 10, 590)

	return x, y
}

func getContinent(country string) string {
	country = strings.ToUpper(strings.TrimSpace(country))

	northAmerica := map[string]bool{"USA": true, "CANADA": true, "MEXICO": true}
	southAmerica := map[string]bool{"BRAZIL": true, "ARGENTINA": true, "CHILE": true, "PERU": true, "COLOMBIA": true, "COSTA_RICA": true}
	europe := map[string]bool{"UK": true, "FRANCE": true, "GERMANY": true, "SPAIN": true, "ITALY": true,
		"PORTUGAL": true, "NETHERLANDS": true, "BELGIUM": true, "SWITZERLAND": true, "AUSTRIA": true,
		"CZECHIA": true, "POLAND": true, "SWEDEN": true, "NORWAY": true, "DENMARK": true, "FINLAND": true,
		"IRELAND": true, "GREECE": true, "SLOVAKIA": true, "HUNGARY": true, "ROMANIA": true, "BELARUS": true}
	asia := map[string]bool{"JAPAN": true, "CHINA": true, "SOUTH_KOREA": true, "INDIA": true, "THAILAND": true,
		"INDONESIA": true, "PHILIPPINES": true, "TAIWAN": true, "QATAR": true, "UNITED_ARAB_EMIRATES": true, "SAUDI_ARABIA": true}
	oceania := map[string]bool{"AUSTRALIA": true, "NEW_ZEALAND": true, "NEW_CALEDONIA": true, "FRENCH_POLYNESIA": true, "NETHERLANDS_ANTILLES": true}

	if northAmerica[country] {
		return "Amérique du Nord"
	}
	if southAmerica[country] {
		return "Amérique du Sud"
	}
	if europe[country] {
		return "Europe"
	}
	if asia[country] {
		return "Asie"
	}
	if oceania[country] {
		return "Océanie"
	}

	return "Autre"
}