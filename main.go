package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

func main() {
	var filename, output_folder, cities_file string
	var add_reviews bool

	flag.StringVar(&filename, "input", "input.json", "JSON file with the data to convert")
	flag.StringVar(&output_folder, "output-folder", "data", "Folder where all the csv files will be saved")
	flag.StringVar(&cities_file, "cities", "cities.txt", "File containting all the cities to add")
	flag.BoolVar(&add_reviews, "add-reviews", true, "Add fake review to the data")
	flag.Parse()

	fmt.Println("Parsing JSON and generating CSV files into folder:", output_folder)

	GenerateFiles(filename, output_folder, cities_file)

	if add_reviews {
		fmt.Println("Adding fake review as nodes...")
		AddReviews(
			filepath.Join(output_folder, "rel_authored.csv"),
			output_folder,
		)
	}

	fmt.Println("CSV creation complete, the data can be found in the folder:", output_folder)
	fmt.Println("Run: ./neo4j_import_A2.sh to load the data in neo4j for task A.2")
	fmt.Println("Run: ./neo4j_import_A3.sh to load the data in neo4j for task A.3 (has reviews as nodes, and include companies and universities)")
}
