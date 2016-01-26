package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/miku/holdingfile"
	"github.com/miku/holdingfile/kbart"
)

func main() {
	date := flag.String("date", "", "record date")
	filename := flag.String("file", "", "holding file")
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

	r := kbart.NewReader(file)
	holdings, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	t, err := time.Parse("2006-01-02", *date)
	if err != nil {
		log.Fatal(err)
	}

	s := holdingfile.Signature{Date: *date, Volume: *volume, Issue: *issue}

	lics := holdings.Licenses(*issn)

	for i, l := range lics {
		if *verbose {
			log.Printf("%+v", l)
		}
		cov := l.Covers(s)
		wall := l.TimeRestricted(t)
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
