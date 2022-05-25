package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strconv"
)

type Set map[string]bool

func getRelationships(authored_file string) (map[string]Set, []string) {
	f, err := os.Open(authored_file)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(f)

	reader.Comma = ';'

	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// For each article save a set of its authors
	article_authors := make(map[string]Set)

	// Additionally, keep a set of all authors
	authors := make(Set)

	for _, lines := range lines {
		article := lines[0]
		author := lines[1]

		if article_authors[article] == nil {
			article_authors[article] = make(Set)
		}

		article_authors[article][author] = true
		authors[author] = true
	}

	auth_list := make([]string, len(authors))

	// Convert the author set to a list so that we can access randomly
	i := 0
	for k := range authors {
		auth_list[i] = k
		i++
	}

	return article_authors, auth_list
}

func getReviewers(author_list []string, article_authors Set, n int) []string {
	result := make([]string, n)
	for i := 0; i < n; i++ {
		var reviewer string
		for { // Make sure the reviewer is not one of the authors
			reviewer = author_list[rand.Intn(len(author_list))]
			if !article_authors[reviewer] {
				break
			}
		}
		result[i] = reviewer
	}

	return result
}

// Adds fake reviews as a relationship. Making sure the reviewers are not
// the paper authors.
func AddReviews(input, output_folder string) {

	rand.Seed(42)

	article_authors, auth_list := getRelationships(input)

	rel_gives_review := create_csv(output_folder, "rel_gives_review.csv")
	defer rel_gives_review.Flush()

	rel_reviews := create_csv(output_folder, "rel_reviews.csv")
	defer rel_reviews.Flush()

	rel_reviewed_about := create_csv(output_folder, "rel_review_about_paper.csv")
	defer rel_reviewed_about.Flush()

	reviews := create_csv(output_folder, "reviews.csv")
	defer reviews.Flush()
	reviews.Write([]string{"ID", "text", "accepted"})

	for article := range article_authors {
		// Add between 3 and 5 reviewers for each paper
		n := 3 + rand.Intn(2)

		// create review node
		counter++
		rev_id := strconv.FormatUint(counter, 10)

		// Set approved with 80% probability
		approved := "true"
		if rand.Float32() > 0.8 {
			approved = "false"
		}

		reviews.Write([]string{rev_id, "PLACEHOLDER TEXT", approved})
		rel_reviewed_about.Write([]string{rev_id, article})

		for _, reviewer := range getReviewers(auth_list, article_authors[article], n) {
			rel_gives_review.Write([]string{reviewer, rev_id})
			rel_reviews.Write([]string{reviewer, article})
		}
	}
}
