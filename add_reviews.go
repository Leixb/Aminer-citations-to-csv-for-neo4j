package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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

func getRelAsMap(filename string) map[string]string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(f)

	reader.Comma = ';'

	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	rel := make(map[string]string)
	for _, line := range lines {
		rel[line[0]] = line[1]
	}
	return rel
}

// some random names
var names = []string{ "John", "Paul", "George", "Ringo", "Pete", "James", "Amelia", "Olivia", "Emma", "Isabella", "Ava", "Sophia", "Charlotte", "Mia", "Lucy", "Grace", "Ruby", "Ella" }
var surnames = []string{ "Smith", "Johnson", "Williams", "Jones", "Davis", "Miller", "Wilson", "Moore", "Taylor", "Anderson", "Jackson", "White", "Harris", "Martin", "Thompson" }

func getRandomNameWithID() (string, string) {
	name := names[rand.Intn(len(names))] + " " + surnames[rand.Intn(len(surnames))]
	id, _ := getId("&&" + name)
	return name, id
}

// Adds fake reviews as a relationship. Making sure the reviewers are not
// the paper authors.
func AddReviews(input, output_folder string) {

	rand.Seed(42)

	article_authors, auth_list := getRelationships(input)

	paperPublishedIn := getRelAsMap(filepath.Join(output_folder, "rel_published.csv"))
	// publicationInVenue := getRelAsMap(filepath.Join(output_folder, "rel_belongs.csv"))
	submittedTo := getRelAsMap(filepath.Join(output_folder, "rel_submittedTo.csv"))

	committees := create_csv(output_folder, "committees.csv")
	committees.Write([]string{"ID"})
	defer committees.Flush()

	managers := create_csv(output_folder, "managers.csv")
	managers.Write([]string{"ID", "pName"})
	defer managers.Flush()

	rel_handled_by := create_csv(output_folder, "rel_handled_by.csv")
	defer rel_handled_by.Flush()

	rel_assigns := create_csv(output_folder, "rel_assigns.csv")
	defer rel_assigns.Flush()

	rel_member_of := create_csv(output_folder, "rel_member_of.csv")
	defer rel_member_of.Flush()

	rel_approves := create_csv(output_folder, "rel_approves.csv")
	defer rel_approves.Flush()

	rel_rejects := create_csv(output_folder, "rel_rejects.csv")
	defer rel_rejects.Flush()

	rel_makes_review := create_csv(output_folder, "rel_makes_review.csv")
	defer rel_makes_review.Flush()

	reviews := create_csv(output_folder, "reviews.csv")
	defer reviews.Flush()
	reviews.Write([]string{"ID", "text"})

	venueManagerMap := make(map[string]string)

	for article := range article_authors {

		committeeId := nextId()
		committees.Write([]string{committeeId})

		// If paper is in list of publications, it has been approved
		_, approved := paperPublishedIn[article]
		venueId := submittedTo[article]

		if _, ok := venueManagerMap[venueId]; !ok {
			manager, manager_id := getRandomNameWithID()
			managers.Write([]string{manager_id, manager})
			rel_handled_by.Write([]string{venueId, manager_id})

			venueManagerMap[venueId] = manager_id
		}

		manager_id := venueManagerMap[venueId]

		rel_assigns.Write([]string{manager_id, committeeId})

		// Add between 3 and 5 reviewers for each paper
		n := 3 + rand.Intn(2)

		// create review node
		rev_id := nextId()

		if approved { //approved if in published list
			rel_approves.Write([]string{rev_id, article})
		} else { //rejected
			rel_rejects.Write([]string{rev_id, article})
		}

		reviews.Write([]string{rev_id, "PLACEHOLDER TEXT"})
		rel_makes_review.Write([]string{committeeId, rev_id})

		for _, reviewer := range getReviewers(auth_list, article_authors[article], n) {
			rel_member_of.Write([]string{reviewer, committeeId})
		}
	}
}
