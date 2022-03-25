package main

import (
	"encoding/csv"
	"flag"
	"log"
	"math/rand"
	"os"
)

func main() {
    var input, output string
    flag.StringVar(&input, "input file", "data/rel_authored.csv", "File containing the article -> author relation in csv format")
    flag.StringVar(&output, "output file", "data/rel_reviews.csv", "Output file where the review relation will be saved")
    flag.Parse()

    f, err := os.Open(input)
    if err != nil {
        log.Fatal(err)
    }
    reader := csv.NewReader(f)

    reader.Comma = ';'

    lines, err := reader.ReadAll()
    if err != nil {
        log.Fatal(err)
    }

    article_authors := make(map[string]map[string]bool)

    authors := make(map[string]bool)

    for _, lines := range lines {
        article := lines[0]
        author := lines[1]

        if article_authors[article] == nil {
            article_authors[article] = make(map[string]bool)
        }

        article_authors[article][author] = true
        authors[author] = true
    }

    auth_list := make([]string, len(authors))

    i := 0
    for k := range authors {
        auth_list[i] = k
        i++
    }

    fr, err := os.Create(output)
    if err != nil {
        log.Fatal(err)
    }
    w := csv.NewWriter(fr)
    defer w.Flush()
    w.Comma = ';';

    for article := range article_authors {
        for i := 0; i < 3 + rand.Intn(2); i++ {
            var reviewer string
            for {
                reviewer = auth_list[rand.Intn(len(auth_list))]
                if !article_authors[article][reviewer] {
                    break
                }
            }
            w.Write([]string{ article, reviewer })
        }
    }
}
