package kbart

import "github.com/miku/holdingfile"

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

type Entries []Entry

func (e Entries) Licenses(issn holdingfile.ISSN) []License {
	for _, entry := range e {
		if entry.PrintIdentifier == issn || entry.OnlineIdentifier == issn {

		}
	}
}
