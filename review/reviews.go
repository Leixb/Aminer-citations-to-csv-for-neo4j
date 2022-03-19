package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
)

func main() {
    f, _ := os.Open("rel_authored.csv")
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

    fr, _ := os.Create("rel_reviews.csv")
    w := csv.NewWriter(fr)
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
    w.Flush()
}
