package kbart

import (
	"bufio"
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

// delayPattern fixes allowed embargo strings.
var delayPattern = regexp.MustCompile(`([P|R])([0-9]+)([Y|M|D])`)

// embargo is a string representing a delay, e.g. P1Y, R10M.
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

// DisallowEarlier returns true if dates *before* the boundary should be
// disallowed.
func (e embargo) DisallowEarlier() bool {
	return strings.HasPrefix(strings.TrimSpace(string(e)), "R")
}

// Reader reads tab-separated KBART. The encoding/csv package did not like
// that particular format so we use a simple bufio.Reader for now.
type Reader struct {
	r          *bufio.Reader
	currentRow int

	SkipFirstRow           bool
	SkipMissingIdentifiers bool
	SkipIncompleteLines    bool
	SkipInvalidEmbargo     bool
}

// NewReader creates a new KBART reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{SkipFirstRow: true, SkipMissingIdentifiers: true, r: bufio.NewReader(r)}
}

// ReadAll loads entries from a reader.
func (r *Reader) ReadAll() (holdingfile.Entries, error) {
	entries := make(holdingfile.Entries)

	for {
		cols, entry, err := r.Read()

		if err == io.EOF {
			break
		}

		switch err {
		case ErrMissingIdentifiers:
			if !r.SkipMissingIdentifiers {
				return entries, err
			}
		case ErrIncompleteLine:
			if !r.SkipIncompleteLines {
				return entries, err
			}
		case ErrInvalidEmbargo:
			if !r.SkipInvalidEmbargo {
				return entries, err
			}
		}

		pi := strings.TrimSpace(cols.PrintIdentifier)
		oi := strings.TrimSpace(cols.OnlineIdentifier)

		if pi == "" && oi == "" {
			if !r.SkipMissingIdentifiers {
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
		if _, err := r.r.ReadString('\n'); err != nil {
			return cols, entry, err
		}
	}
	r.currentRow++

	var line string
	var err error

	for {
		line, err = r.r.ReadString('\n')
		if strings.TrimSpace(line) != "" {
			break
		}
		if err == io.EOF {
			return cols, entry, io.EOF
		}
	}

	record := strings.Split(line, "\t")

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
		PrintIdentifier:  record[1],
		OnlineIdentifier: record[2],
		FirstIssueDate:   record[3],
		FirstVolume:      record[4],
		FirstIssue:       record[5],
		LastIssueDate:    record[6],
		LastVolume:       record[7],
		LastIssue:        record[8],
		Embargo:          embargo(record[12]),
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
