package kbart

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/miku/holdingfile"
)

func TestEmbargeDisallowEarlier(t *testing.T) {
	var cases = []struct {
		e               embargo
		disallowEarlier bool
	}{
		{embargo("1"), false},
		{embargo("R1"), true},
		{embargo("R1D"), true},
		{embargo("R10M"), true},
		{embargo("?10M"), false},
	}

	for _, c := range cases {
		got := c.e.DisallowEarlier()
		if got != c.disallowEarlier {
			t.Errorf("embargo.DisallowEarlier() got %v, want %v", got, c.disallowEarlier)
		}
	}
}

func TestEmbargoDuration(t *testing.T) {
	var cases = []struct {
		e   embargo
		d   time.Duration
		err error
	}{
		{embargo("1"), time.Duration(0), ErrIncompleteEmbargo},
		{embargo("R1"), time.Duration(0), ErrIncompleteEmbargo},
		{embargo("R1D"), -24 * time.Hour, nil},
		{embargo("R10M"), -7200 * time.Hour, nil},
		{embargo("?10M"), time.Duration(0), ErrIncompleteEmbargo},
	}

	for _, c := range cases {
		got, err := c.e.AsDuration()
		if err != c.err {
			t.Errorf("embargo.AsDuration() got %v, want %v", err, c.err)
		}
		if got != c.d {
			t.Errorf("embargo.DisallowEarlier() got %v, want %v", got, c.d)
		}
	}
}

func TestFromReader(t *testing.T) {
	var tests = []struct {
		r       io.Reader
		entries Entries
		err     error
	}{
		{r: strings.NewReader("first line is discarded"),
			entries: Entries{Map: map[holdingfile.ISSN][]holdingfile.Entry{}},
			err:     nil},
		{r: strings.NewReader("xxx\nyyy\nzzz"),
			entries: Entries{Map: map[holdingfile.ISSN][]holdingfile.Entry{}},
			err:     ErrIncompleteLine},
		{r: strings.NewReader(`publication_title	print_identifier	online_identifier	date_first_issue_online	num_first_vol_online	num_first_issue_online	date_last_issue_online	num_last_vol_online	num_last_issue_online	title_url	first_author	title_id	embargo_info	coverage_depth	coverage_notes	publisher_name	own_anchor	il_relevance	il_nationwide	il_electronic_transmission	il_comment	all_issns	zdb_id
Bill of Rights Journal (via Hein Online)	0006-2499		1968	1		1996	29		http://heinonline.org/HOL/Index?index=journals/blorij&collection=journals		227801		Volltext		via Hein Online		Keine Fernleihe				0006-2499	2805467-2`),
			entries: Entries{
				Map: map[holdingfile.ISSN][]holdingfile.Entry{
					holdingfile.ISSN("0006-2499"): []holdingfile.Entry{
						holdingfile.Entry{
							Begin: holdingfile.Signature{
								Date:   "1968",
								Volume: "1",
								Issue:  "",
							},
							End: holdingfile.Signature{
								Date:   "1996",
								Volume: "29",
								Issue:  "",
							},
							Embargo:                time.Duration(0),
							EmbargoDisallowEarlier: false,
						}}}},
			err: nil},
	}

	for _, test := range tests {
		entries, err := ReadEntries(test.r)
		if err != test.err {
			t.Errorf("FromReader error got %+v, want %+v", err, test.err)
		}
		if !reflect.DeepEqual(entries, test.entries) {
			t.Errorf("FromReader Entries got %+v, want %+v", entries, test.entries)
			for _, s := range pretty.Diff(test.entries, entries) {
				t.Errorf(s)
			}
		}
	}
}
