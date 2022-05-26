package main

// Copyright (C) 2022  LeixB

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

import (
	"flag"
	"fmt"
	"path/filepath"
)

func main() {
	var filename, output_folder, cities_file string
	var add_reviews, add_venue_related bool

	flag.StringVar(&filename, "input", "input.json", "JSON file with the data to convert")
	flag.StringVar(&output_folder, "output-folder", "data", "Folder where all the csv files will be saved")
	flag.StringVar(&cities_file, "cities", "cities.txt", "File containting all the cities to add")
	flag.BoolVar(&add_reviews, "add-reviews", true, "Add fake review to the data")
	flag.BoolVar(&add_venue_related, "add-venue-area", true, "Add fake venue relations to the data")
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

	if add_venue_related {
		fmt.Println("Adding fake venue relations as edges...")
		AddVenueRelated(
			output_folder,
		)
	}

	fmt.Println("CSV creation complete, the data can be found in the folder:", output_folder)
	fmt.Println("Run: ./neo4j_import_A2.sh to load the data in neo4j for task A.2")
	fmt.Println("Run: ./neo4j_import_A3.sh to load the data in neo4j for task A.3 (has reviews as nodes, and include companies and universities)")
}
