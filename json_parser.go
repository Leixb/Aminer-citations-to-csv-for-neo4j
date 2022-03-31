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
	"path/filepath"
	"strconv"
	"strings"
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
	ID    string `json:"_id"`
	Name  string `json:"name"`
    Org   string `json:"org"`
    OrgID string `json:"orgid"`
    GID   string `json:"gid"`
}

type Paper struct {
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
	ISSN       string   `json:"issn"`
	ISBN       string   `json:"isbn"`
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

func create_csv(folder string, filename string) *csv.Writer {
	f, err := os.Create(filepath.Join(folder, filename))
	if err != nil {
		log.Fatal("Failed creating file", filename, "in folder", folder)
	}
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

var venue_type map[string]VenueType

func process_venue(v *Venue, year int, volume, isbn *string, city_names []string,
	f_journal, f_conference, f_workshop,
	f_edition, f_volume, rel_belongs *csv.Writer,
) (string, bool) {
	// Make sure ID is valid
	var done bool
    var venue_id, pub_id string

	venueType := Journal // Default to journal

    if v.Name == "" {
        v.Name = v.Raw
    }

	venue_id, done = getId(v.Name)
	if done { // Get venue type from map
		venueType = venue_type[venue_id]
	} else { // If venue does not exist, create it
		var f_venue *csv.Writer

		// Determine venue type
		if v.T == "C" || v.Vtype == 10 {
			venueType = Conference
		} else if v.Vtype == 2 {
			venueType = Workshop
		}

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

		// Register venue type
		venue_type[venue_id] = venueType

		f_venue.Write([]string{
			venue_id, v.Name,
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

func GenerateFiles(filename, output_folder, cities_file string) {

	counter = 0
	ids = make(map[string]string)
	venue_type = make(map[string]VenueType)

	if err := os.MkdirAll(output_folder, 0755); err != nil {
		log.Fatal("Output directory file creation failed on", output_folder)
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to open input file:", filename)
	}
	dec := json.NewDecoder(f)

	f_articles := create_csv(output_folder, "articles.csv")
	f_authors := create_csv(output_folder, "authors.csv")
	f_conference := create_csv(output_folder, "conference.csv")
	f_edition := create_csv(output_folder, "edition.csv")
	f_journal := create_csv(output_folder, "journal.csv")
	f_keywords := create_csv(output_folder, "keywords.csv")
	f_volume := create_csv(output_folder, "volume.csv")
	f_workshop := create_csv(output_folder, "workshop.csv")
    f_university := create_csv(output_folder, "university.csv")
    f_company := create_csv(output_folder, "company.csv")

	rel_authored := create_csv(output_folder, "rel_authored.csv")
	rel_belongs := create_csv(output_folder, "rel_belongs.csv")
	rel_cites := create_csv(output_folder, "rel_cites.csv")
	rel_keywords := create_csv(output_folder, "rel_keywords.csv")
	rel_published := create_csv(output_folder, "rel_published.csv")
	rel_affiliated := create_csv(output_folder, "rel_affiliated.csv")

	for _, handle := range []*csv.Writer{
		f_articles, f_authors, f_conference, f_edition, f_journal, f_keywords, f_volume, f_workshop, f_university, f_company,
		rel_authored, rel_belongs, rel_cites, rel_keywords, rel_published, rel_affiliated} {
		defer handle.Flush()
	}

	f_cities, err := os.Open(cities_file)
	if err != nil {
		log.Fatal("Failed to open city file list", cities_file)
	}
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
	_, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	// while the array contains values
	for dec.More() {
		select {
		case <-sigc:
			fmt.Println("Stopping")
			return
		default:
		}

		var paper Paper
		// decode an array value (Message)
		err := dec.Decode(&paper)
		if err != nil {
			log.Fatal(err)
		}

		art_id, _ := getId(paper.ID)
		for _, author := range paper.Authors {
			if author.ID == "" {
				continue
			}
			auth_id, done := getId(author.ID)
			if !done {
				f_authors.Write([]string{
					auth_id,
					author.Name,
				})

                if author.Org != "" {
                    name, _, _ := strings.Cut(author.Org, ",")

                    orgid, done := getId(strings.ToLower(name))
                    if !done {
                        if strings.HasPrefix(author.Org, "Uni") {
                            f_university.Write([]string{
                                orgid, name,
                            })
                        } else {
                            f_company.Write([]string{
                                orgid, name,
                            })
                        }
                    }

                    rel_affiliated.Write([]string{
                        auth_id, orgid,
                    })
                }
			}
			rel_authored.Write([]string{
				art_id,
				auth_id,
			})
		}

		pub_id, ok := process_venue(&paper.Venue, paper.Year, &paper.Volume, &paper.ISBN, city_names,
			f_journal, f_conference, f_workshop,
			f_edition, f_volume, rel_belongs,
		)
		if ok {
			rel_published.Write([]string{art_id, pub_id})
		}

		for _, keyword := range paper.Keywords {
			if keyword == "" {
				continue
			}
			key_id, done := getId("@" + keyword) // add @ so that it does not clash with sids
			if !done {
				f_keywords.Write([]string{key_id, keyword})
			}
			rel_keywords.Write([]string{art_id, key_id})
		}

		for _, reference := range paper.References {
			ref_id, _ := getId(reference)
            // Skip self references
            if ref_id == art_id {
                continue
            }
			rel_cites.Write([]string{art_id, ref_id})
		}

		f_articles.Write([]string{
			art_id,
			paper.Title,
			paper.Abstract,
			paper.Doi,
			strconv.Itoa(paper.Year),
		})
	}
}
