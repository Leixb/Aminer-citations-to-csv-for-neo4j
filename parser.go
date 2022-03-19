package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var ids map[string]string
var counter uint64

type Venue struct {
	ID    string `json:"_id"`
	SID   string `json:"sid"`
	Name  string `json:"name_d"`
	T     string `json:"t"`
	Vtype int    `json:"type"` // 2 -> workshop, 10 -> conf. 1/11 -> journal, 3/12 -> book
	Raw   string `json:"raw"`
}

type Author struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type Article struct {
	ID         string   `json:"_id"`
	Title      string   `json:"title"`
	Authors    []Author `json:"authors"`
	Venue      Venue    `json:"venue"`
	Year       int      `json:"year"`
	Keywords   []string `json:"keywords"`
	Fos        []string `json:"fos"`
	N_citation int      `json:"n_citation"`
	Page_start string   `json:"page_start"`
	Page_end   string   `json:"page_end"`
	Lang       string   `json:"lang"`
	Volume     string   `json:"volume"`
	Issue      string   `json:"issue"`
	Issn       string   `json:"issn"`
	Isbn       string   `json:"isbn"`
	Doi        string   `json:"doi"`
	Pdf        string   `json:"pdf"`
	Url        []string `json:"url"`
	Abstract   string   `json:"abstract"`
	References []string `json:"references"`
}

func getId(value string) (string, bool) {
	if val, ok := ids[value]; ok {
		return val, true
	}

	counter++
	scount := strconv.FormatUint(counter, 10)
	ids[value] = scount

	return scount, false
}

func create_csv(filename string) *csv.Writer {
	f, _ := os.Create(filename)
	w := csv.NewWriter(f)
	w.Comma = ';'

	return w
}

type VenueType int

const (
	Journal VenueType = iota
	Conference
	Workshop
)

func process_venue(v *Venue, year int, volume *string, isbn *string, city_names []string,
	f_journal *csv.Writer, f_conference *csv.Writer, f_workshop *csv.Writer,
	f_edition *csv.Writer, f_volume *csv.Writer, rel_belongs *csv.Writer,
) (string, bool) {
	// Make sure ID is valid
	if v.ID == "" {
		if v.SID == "" {
			return "", false
		}
		v.ID = "!" + v.SID
	}

	var done bool
	var venue_id string
	var pub_id string

	venueType := Journal // Default to journal

	if v.T == "C" || v.Vtype == 10 {
		venueType = Conference
	} else if v.Vtype == 2 {
		venueType = Workshop
	}

	venue_id, done = getId(v.ID)
	// If venue does not exist, create it
	if !done {
		var f_venue *csv.Writer
		switch venueType {
		case Journal:
			f_venue = f_journal
		case Conference:
			f_venue = f_conference
		case Workshop:
			f_venue = f_workshop
		default:
			log.Fatalf("Error processing venue %v", v)
		}
		f_venue.Write([]string{
			venue_id, v.Name, v.Raw,
		})
	}

	switch venueType {
	case Journal:
		if pub_id, done = getId(fmt.Sprintf("%s-%d-%s", venue_id, year, *volume)); !done {
			f_volume.Write([]string{
				pub_id, strconv.Itoa(year), *volume, *isbn,
			})
		}
	case Conference:
		if pub_id, done = getId(fmt.Sprintf("%s-%d-C", venue_id, year)); !done {
			f_edition.Write([]string{
				pub_id, strconv.Itoa(year), city_names[rand.Intn(len(city_names))],
			})
		}
	case Workshop:
		if pub_id, done = getId(fmt.Sprintf("%s-%d-W", venue_id, year)); !done {
			f_edition.Write([]string{
				pub_id, strconv.Itoa(year), city_names[rand.Intn(len(city_names))],
			})
		}
	default:
		log.Fatalf("Error processing venue %v", v)
	}

	if !done {
		rel_belongs.Write([]string{
			pub_id, venue_id,
		})
	}

	return pub_id, true

}

func main() {
	counter = 0
	ids = map[string]string{}
	f, _ := os.Open("input.json")
	dec := json.NewDecoder(f)

	f_articles := create_csv("articles.csv")
	f_authors := create_csv("authors.csv")
	f_conference := create_csv("conference.csv")
	f_edition := create_csv("edition.csv")
	f_journal := create_csv("journal.csv")
	f_keywords := create_csv("keywords.csv")
	f_volume := create_csv("volume.csv")
	f_workshop := create_csv("workshop.csv")

	rel_authored := create_csv("rel_authored.csv")
	rel_belongs := create_csv("rel_belongs.csv")
	rel_cites := create_csv("rel_cites.csv")
	rel_keywords := create_csv("rel_keywords.csv")
	rel_published := create_csv("rel_published.csv")

	handles := []*csv.Writer{
		f_articles,
		f_authors,
		f_conference,
		f_edition,
		f_journal,
		f_keywords,
		f_volume,
		f_workshop,

		rel_authored,
		rel_belongs,
		rel_cites,
		rel_keywords,
		rel_published,
	}

	f_cities, _ := os.Open("cities")
	var city_names []string
	scanner := bufio.NewScanner(f_cities)
	for scanner.Scan() {
		city_names = append(city_names, scanner.Text())
	}
	f_cities.Close()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// read open bracket
	_, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	// while the array contains values
loop:
	for dec.More() {
		select {
		case <-sigc:
			fmt.Println("Stopping")
			break loop
		default:
		}

		var a Article
		// decode an array value (Message)
		err := dec.Decode(&a)
		if err != nil {
			log.Fatal(err)
		}

		art_id, _ := getId(a.ID)
		for _, author := range a.Authors {
			if author.ID == "" {
				continue
			}
			auth_id, done := getId(author.ID)
			if !done {
				f_authors.Write([]string{
					auth_id,
					author.Name,
				})
			}
			rel_authored.Write([]string{
				art_id,
				auth_id,
			})
		}

		pub_id, ok := process_venue(&a.Venue, a.Year, &a.Volume, &a.Isbn, city_names,
			f_journal, f_conference, f_workshop,
			f_edition, f_volume, rel_belongs,
		)
		if ok {
			rel_published.Write([]string{art_id, pub_id})
		}

		for _, keyword := range a.Keywords {
			if keyword == "" {
				continue
			}
			key_id, done := getId("@" + keyword) // add @ so that it does not clash with sids
			if !done {
				f_keywords.Write([]string{key_id, keyword})
			}
			rel_keywords.Write([]string{art_id, key_id})
		}

		for _, reference := range a.References {
			ref_id, _ := getId(reference)
			rel_cites.Write([]string{art_id, ref_id})
		}

		f_articles.Write([]string{
			art_id,
			a.Title,
			a.Abstract,
			a.Doi,
			strconv.Itoa(a.Year),
		})
	}
	for _, handle := range handles {
		handle.Flush()
	}
}
