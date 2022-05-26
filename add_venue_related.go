package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
)

func getVenues(input string) []string {
	f, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(f)

	reader.Comma = ';'

	reader.Read() // skip header

	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// return first column
	var venues []string
	for _, line := range lines {
		venues = append(venues, line[0])
	}
	return venues
}

func getAreas(input string) []string {
	f, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(f)

	reader.Comma = ';'

	reader.Read() // skip header

	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// return first column
	var areas []string
	for _, line := range lines {
		areas = append(areas, line[0])
	}
	return areas
}

// Adds fake reviews as a relationship. Making sure the reviewers are not
// the paper authors.
func AddVenueRelated(output_folder string) {

	rand.Seed(42)

	venueIDs := append(getVenues("./data/journal.csv"), getVenues("./data/conference.csv")...)
	venueIDs = append(venueIDs, getVenues("./data/workshop.csv")...)
	areaIDs := getAreas("./data/area.csv")

	rel_venue_related := create_csv(output_folder, "rel_venue_related.csv")
	defer rel_venue_related.Flush()

	for _, venue := range venueIDs {
		// Add between 3 and 5 related areas for each venue
		n := 3 + rand.Intn(2)

		for i := 0; i < n; i++ {
			// Pick a random area
			area := areaIDs[rand.Intn(len(areaIDs))]
			rel_venue_related.Write([]string{venue, area})
		}
	}
}
