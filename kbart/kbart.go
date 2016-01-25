package kbart

import (
	"encoding/csv"
	"io"

	"github.com/miku/holdingfile"
)

// Knowledge Bases And Related Tools working group.
type Entry struct {
	PublicationTitle         string
	PrintIdentifier          string
	OnlineIdentifier         string
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

func (e Entry) Covers(s holdingfile.Signature) (bool, error) {
	return false, nil
}

// Entries holds a list of license entries.
type Entries []Entry

// FromReader loads entries from a tab-separated file.
func (e Entries) FromReader(r io.Reader) error {
	reader := csv.NewReader(r)
	reader.Comma = '\t'

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		e = append(e, Entry{
			PublicationTitle:         record[0],
			PrintIdentifier:          record[1],
			OnlineIdentifier:         record[2],
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
			Publisher:                record[21],
			ZDBID:                    record[23]})
	}
	return nil
}

func (e Entries) Licenses(issn holdingfile.ISSN) []holdingfile.License {
	for _, entry := range e {
		if entry.PrintIdentifier == issn || entry.OnlineIdentifier == issn {
		}
	}
}
