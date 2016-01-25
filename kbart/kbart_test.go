package kbart

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/kr/pretty"
	"github.com/miku/holdingfile"
)

func TestFromReader(t *testing.T) {
	var tests = []struct {
		r       io.Reader
		entries Entries
		err     error
	}{
		{r: strings.NewReader("xxx"), err: nil},
		{r: strings.NewReader("xxx\nyyy\nzzz"), err: ErrorIncompleteLine},
		{r: strings.NewReader(`publication_title	print_identifier	online_identifier	date_first_issue_online	num_first_vol_online	num_first_issue_online	date_last_issue_online	num_last_vol_online	num_last_issue_online	title_url	first_author	title_id	embargo_info	coverage_depth	coverage_notes	publisher_name	own_anchor	il_relevance	il_nationwide	il_electronic_transmission	il_comment	all_issns	zdb_id
Bill of Rights Journal (via Hein Online)	0006-2499		1968	1		1996	29		http://heinonline.org/HOL/Index?index=journals/blorij&collection=journals		227801		Volltext		via Hein Online		Keine Fernleihe				0006-2499	2805467-2`),
			entries: Entries{
				Entry{
					PublicationTitle:         "Bill of Rights Journal (via Hein Online)",
					PrintIdentifier:          holdingfile.ISSN("0006-2499"),
					OnlineIdentifier:         holdingfile.ISSN(""),
					FirstIssueDate:           "1968",
					FirstVolume:              "1",
					FirstIssue:               "",
					LastIssueDate:            "1996",
					LastVolume:               "29",
					LastIssue:                "",
					TitleURL:                 "http://heinonline.org/HOL/Index?index=journals/blorij&collection=journals",
					FirstAuthor:              "",
					TitleID:                  "227801",
					Embargo:                  "",
					CoverageDepth:            "Volltext",
					CoverageNotes:            "",
					PublisherName:            "via Hein Online",
					InterlibraryRelevance:    "Keine Fernleihe",
					InterlibraryNationwide:   "",
					InterlibraryTransmission: "",
					InterlibraryComment:      "",
					Publisher:                "",
					Anchor:                   "",
					ZDBID:                    "2805467-2",
				}}, err: nil},
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
