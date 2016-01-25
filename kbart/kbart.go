package kbart

import (
	"encoding/csv"
	"errors"
	"io"
	"time"

	"github.com/miku/holdingfile"
)

var ErrorIncompleteLine = errors.New("incomplete KBART line")

// Knowledge Bases And Related Tools working group.
type Entry struct {
	PublicationTitle         string
	PrintIdentifier          holdingfile.ISSN
	OnlineIdentifier         holdingfile.ISSN
	FirstIssueDate           string
	FirstVolume              string
	FirstIssue               string
	LastIssueDate            string
	LastVolume               string
	LastIssue                string
	TitleURL                 string
	FirstAuthor              string
	TitleID                  string
	Embargo                  string
	CoverageDepth            string
	CoverageNotes            string
	PublisherName            string
	InterlibraryRelevance    string
	InterlibraryNationwide   string
	InterlibraryTransmission string
	InterlibraryComment      string
	Publisher                string
	Anchor                   string
	ZDBID                    string
}

func (e Entry) Covers(s holdingfile.Signature) error {
	return nil
}

func (e Entry) TimeRestricted(t time.Time) error {
	return nil
}

// Entries holds a list of license entries.
type Entries []Entry

// FromReader loads entries from a reader. Must be a tab-separated CSV with
// exactly one header row.
func ReadEntries(r io.Reader) (Entries, error) {
	var entries Entries

	reader := csv.NewReader(r)
	reader.Comma = '\t'

	if _, err := reader.Read(); err != nil {
		return entries, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return entries, err
		}
		if len(record) < 23 {
			return nil, ErrorIncompleteLine
		}
		entries = append(entries, Entry{
			PublicationTitle:         record[0],
			PrintIdentifier:          holdingfile.ISSN(record[1]),
			OnlineIdentifier:         holdingfile.ISSN(record[2]),
			FirstIssueDate:           record[3],
			FirstVolume:              record[4],
			FirstIssue:               record[5],
			LastIssueDate:            record[6],
			LastVolume:               record[7],
			LastIssue:                record[8],
			TitleURL:                 record[9],
			FirstAuthor:              record[10],
			TitleID:                  record[11],
			Embargo:                  record[12],
			CoverageDepth:            record[13],
			CoverageNotes:            record[14],
			PublisherName:            record[15],
			Anchor:                   record[16],
			InterlibraryRelevance:    record[17],
			InterlibraryNationwide:   record[18],
			InterlibraryTransmission: record[19],
			InterlibraryComment:      record[20],
			ZDBID:                    record[22]})
	}
	return entries, nil
}

func (e Entries) Licenses(issn holdingfile.ISSN) []holdingfile.License {
	for _, entry := range e {
		if entry.PrintIdentifier == issn || entry.OnlineIdentifier == issn {
		}
	}
	return nil
}
