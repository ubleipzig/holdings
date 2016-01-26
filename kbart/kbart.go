package kbart

import (
	"encoding/csv"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/miku/holdingfile"
)

var (
	ErrIncompleteLine     = errors.New("incomplete KBART line")
	ErrIncompleteEmbargo  = errors.New("incomplete embargo")
	ErrInvalidEmbargo     = errors.New("invalid embargo")
	ErrMissingIdentifiers = errors.New("missing identifiers")
)

var delayPattern = regexp.MustCompile(`([P|R])([0-9]+)([Y|M|D])`)

type embargo string

// entry represents the various columns.
type columns struct {
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
	Embargo                  embargo
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

// Convert string like P12M, P1M, R10Y into a time.Duration.
func (e embargo) AsDuration() (time.Duration, error) {
	var d time.Duration

	emb := strings.TrimSpace(string(e))
	if len(emb) == 0 {
		return d, nil
	}

	var parts = delayPattern.FindStringSubmatch(emb)
	if len(parts) < 4 {
		return d, ErrIncompleteEmbargo
	}

	i, err := strconv.Atoi(parts[2])
	if err != nil {
		return d, ErrInvalidEmbargo
	}

	switch parts[3] {
	case "D":
		return time.Duration(-i) * 24 * time.Hour, nil
	case "M":
		return time.Duration(-i) * 24 * time.Hour * 30, nil
	case "Y":
		return time.Duration(-i) * 24 * time.Hour * 30 * 365, nil
	default:
		return d, ErrInvalidEmbargo
	}
}

func (e embargo) DisallowEarlier() bool {
	return strings.HasPrefix(strings.TrimSpace(string(e)), "R")
}

type Reader struct {
	r                        *csv.Reader
	SkipFirstRow             bool
	IgnoreMissingIdentifiers bool
	IgnoreIncompleteLines    bool
	IgnoreInvalidEmbargo     bool
	currentRow               int
}

func NewReader(r io.Reader) *Reader {
	cr := csv.NewReader(r)
	cr.Comma = '\t'
	cr.FieldsPerRecord = 23
	cr.LazyQuotes = true
	return &Reader{r: cr}
}

// ReadAll loads entries from a reader. Must be a tab-separated CSV with
// exactly one header row.
func (r *Reader) ReadAll() (holdingfile.Entries, error) {
	entries := make(holdingfile.Entries)

	for {
		cols, entry, err := r.Read()

		if err == io.EOF {
			break
		}

		switch err {
		case ErrMissingIdentifiers:
			if !r.IgnoreMissingIdentifiers {
				return entries, err
			}
		case ErrIncompleteLine:
			if !r.IgnoreIncompleteLines {
				return entries, err
			}
		case ErrInvalidEmbargo:
			if !r.IgnoreInvalidEmbargo {
				return entries, err
			}
		}

		pi := strings.TrimSpace(cols.PrintIdentifier)
		oi := strings.TrimSpace(cols.OnlineIdentifier)

		if pi == "" && oi == "" {
			if !r.IgnoreMissingIdentifiers {
				return entries, ErrMissingIdentifiers
			}
		}
		if pi != "" {
			entries[pi] = append(entries[pi], holdingfile.License(entry))
		}
		if oi != "" {
			entries[oi] = append(entries[oi], holdingfile.License(entry))
		}
	}

	return entries, nil
}

// Read reads a single line.
func (r *Reader) Read() (columns, holdingfile.Entry, error) {
	var entry holdingfile.Entry
	var cols columns

	if r.SkipFirstRow && r.currentRow == 0 {
		if _, err := r.r.Read(); err != nil {
			return cols, entry, err
		}
	}

	record, err := r.r.Read()
	r.currentRow++

	if err == io.EOF {
		return cols, entry, io.EOF
	}
	if err != nil {
		return cols, entry, err
	}
	if len(record) < 23 {
		return cols, entry, ErrIncompleteLine
	}

	cols = columns{
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
		Embargo:                  embargo(record[12]),
		CoverageDepth:            record[13],
		CoverageNotes:            record[14],
		PublisherName:            record[15],
		Anchor:                   record[16],
		InterlibraryRelevance:    record[17],
		InterlibraryNationwide:   record[18],
		InterlibraryTransmission: record[19],
		InterlibraryComment:      record[20],
		ZDBID:                    record[22],
	}

	emb, err := cols.Embargo.AsDuration()
	if err != nil {
		return cols, entry, err
	}

	entry = holdingfile.Entry{
		Begin: holdingfile.Signature{
			Date:   cols.FirstIssueDate,
			Volume: cols.FirstVolume,
			Issue:  cols.FirstIssue,
		},
		End: holdingfile.Signature{
			Date:   cols.LastIssueDate,
			Volume: cols.LastVolume,
			Issue:  cols.LastIssue,
		},
		Embargo:                emb,
		EmbargoDisallowEarlier: cols.Embargo.DisallowEarlier(),
	}

	return cols, entry, nil
}
