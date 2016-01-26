package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/miku/holdings"
	"github.com/miku/holdings/google"
	"github.com/miku/holdings/kbart"
	"github.com/miku/holdings/ovid"
)

var layouts = []string{
	"2006",
	"2006-01-02",
}

func main() {
	date := flag.String("date", "", "record date")
	filename := flag.String("file", "", "holding file")
	format := flag.String("format", "kbart", "holding file format: kbart, ovid or google")
	issn := flag.String("issn", "", "record issn")
	issue := flag.String("issue", "", "record issue")
	volume := flag.String("volume", "", "record volume")
	verbose := flag.Bool("verbose", false, "be verbose")

	flag.Parse()

	if *issn == "" {
		log.Fatal("-issn is required")
	}

	if *filename == "" {
		log.Fatal("a holding -file is required")
	}

	if *date == "" {
		log.Fatal("-date is required")
	}

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}

	var hfile holdings.File

	switch *format {
	case "kbart":
		hfile = kbart.NewReader(file)
	case "ovid":
		hfile = ovid.NewReader(file)
	case "google":
		hfile = google.NewReader(file)
	default:
		log.Fatal("unknown format")

	}

	entries, err := hfile.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var t time.Time

	for _, layout := range layouts {
		t, err = time.Parse(layout, *date)
		if err == nil {
			break
		}
	}

	if t.IsZero() {
		log.Fatalf("could not parse date with any of %s", strings.Join(layouts, ", "))
	}

	s := holdings.Signature{Date: *date, Volume: *volume, Issue: *issue}

	licenses := entries.Licenses(*issn)

	for i, license := range licenses {
		if *verbose {
			log.Printf("%+v", license)
		}

		cov, wall := license.Covers(s), license.TimeRestricted(t)

		if cov == nil && wall == nil {
			fmt.Printf("%d\tOK\tNo restrictions.\n", i)
		}
		if cov != nil {
			fmt.Printf("%d\tNO\tNot covered: %s\n", i, cov)
		}
		if wall != nil {
			fmt.Printf("%d\tNO\tMoving wall applies.\n", i)
		}
	}
}
