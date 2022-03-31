package main

import (
	"flag"
	"fmt"
    "path/filepath"
)

func main() {
	var filename, output, cities_file string
    var add_reviews, reviews_as_node bool

	flag.StringVar(&filename, "input", "input.json", "JSON file with the data to convert")
	flag.StringVar(&output, "output", "data", "Folder where all the csv files will be saved")
	flag.StringVar(&cities_file, "cities", "cities.txt", "File containting all the cities to add")
    flag.BoolVar(&add_reviews, "add-reviews", true, "Add fake review to the data")
    flag.BoolVar(&reviews_as_node, "reviews-as-node", false, "Add reviews as nodes")
	flag.Parse()

    fmt.Println("Parsing JSON and generating CSV files...")

	GenerateFiles(filename, output, cities_file)

    if add_reviews {
        if reviews_as_node {
            fmt.Println("Adding fake review as nodes...")
            AddReviewsNode(
                filepath.Join(output, "rel_authored.csv"),
                output,
            )
        } else {
            fmt.Println("Adding fake review as edges...")
            AddReviewsEdge(
                filepath.Join(output, "rel_authored.csv"),
                output,
            )
        }
    }

    fmt.Println("DONE")
    fmt.Println("Run: ./neo4j_import.sh to load the data in neo4j")
}
