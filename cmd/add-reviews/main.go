package main

import (
	"flag"
    "internal/reviews"
)

func main() {
    var input, output string
    flag.StringVar(&input, "input file", "data/rel_authored.csv", "File containing the article -> author relation in csv format")
    flag.StringVar(&output, "output file", "data/rel_reviews.csv", "Output file where the review relation will be saved")
    flag.Parse()

    reviews.AddReviews(input, output)
}

