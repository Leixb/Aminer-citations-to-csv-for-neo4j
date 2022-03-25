package main

import (
	"flag"
	"fmt"
	"internal/parse_json"
	"internal/reviews"
)

func main() {
	var filename, output, cities_file string
	flag.StringVar(&filename, "input", "input.json", "JSON file with the data to convert")
	flag.StringVar(&output, "output", "data", "Folder where all the csv files will be saved")
	flag.StringVar(&cities_file, "cities", "cities.txt", "File containting all the cities to add")
	flag.Parse()

    fmt.Println("Parsing JSON and generating CSV files...")

	parse_json.GenerateFiles(filename, output, cities_file)

    fmt.Println("CSV file generation complete, adding fake reviews...")

	reviews.AddReviews(
		"data/rel_authored.csv",
		"data/rel_reviews.csv",
	)

    fmt.Println("DONE")
    fmt.Println("Run: ./neo4j_import.sh to load the data in neo4j")
}
