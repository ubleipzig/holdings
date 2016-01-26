package holdingfile

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

var (
	ErrBeforeCoverageInterval = errors.New("before coverage interval")
	ErrAfterCoverageInterval  = errors.New("after coverage interval")
	ErrMissingValues          = errors.New("missing values")
	ErrMovingWall             = errors.New("moving wall")
)

var (
	intPattern = regexp.MustCompile("[0-9]+")
)

// Holdings can return a list of licenses for a given ISSN.
type Holdings interface {
	Licenses(string) []License
}

// License exposes methods to let clients check their own validity in the
// context of this license.
type License interface {
	// Covers looks at the static content only.
	Covers(Signature) error
	// TimeRestricted will report an error, if a moving wall constraint holds.
	TimeRestricted(time.Time) error
}

// Entries holds a list of license entries by ISSN. A simple implementation of
// Holdings.
type Entries map[string][]License

// Licenses make Entries fulfill the holdings interface.
func (e Entries) Licenses(issn string) []License {
	return e[issn]
}

// Entry is a reduced holding file entry. Usually, moving wall allow the
// items, that are earlier then the boundary. If EmbargoDisallowEarlier is
// set, the effect is reversed.
type Entry struct {
	Begin                  Signature
	End                    Signature
	Embargo                time.Duration
	EmbargoDisallowEarlier bool
}

// TimeRestricted returns an error, if the given time falls within the moving
// wall set by the Entry. The embargo is simply added to the current time,
// so it should expressed with negative values.
func (e Entry) TimeRestricted(t time.Time) error {
	var now = time.Now()
	if e.EmbargoDisallowEarlier {
		if t.Before(now.Add(e.Embargo)) {
			return ErrMovingWall
		}
	}
	if t.After(now.Add(e.Embargo)) {
		return ErrMovingWall
	}
	return nil
}

// Covers returns, whether the given signature lies inside the interval
// defined by entry. If there is not comparable date, the volume and issue
// comparisons do not make much sense. However, if there is a date, we are ok
// with just one of volume or issue defined.
func (e Entry) Covers(s Signature) error {
	if err := e.compareDate(s); err != nil {
		return err
	}
	if err := e.compareVolume(s); err != nil {
		return err
	}
	if err := e.compareIssue(s); err != nil {
		return err
	}
	return nil
}

// compareYear returns an error, if both values are defined and disagree, or
// if too few values are defined to do a sane comparison.
func (e Entry) compareDate(s Signature) error {
	if s.Date == "" || (e.Begin.Date == "" && e.End.Date == "") {
		return ErrMissingValues
	}
	if e.Begin.Date != "" {
		if s.Date < e.Begin.Date {
			return ErrBeforeCoverageInterval
		}
	}
	if e.End.Date != "" {
		if s.Date > e.End.Date {
			return ErrAfterCoverageInterval
		}
	}
	return nil
}

// compareVolume returns an error, if both values are defined and disagree,
// otherwise we assume there is no error.
func (e Entry) compareVolume(s Signature) error {
	if s.Volume == "" {
		return nil
	}
	if e.Begin.Volume != "" {
		if s.VolumeInt() < e.Begin.VolumeInt() {
			return ErrBeforeCoverageInterval
		}
	}
	if e.End.Volume != "" {
		if s.VolumeInt() > e.End.VolumeInt() {
			return ErrAfterCoverageInterval
		}
	}
	return nil
}

// compareIssue returns an error, if both values are defined and disagree,
// otherwise we assume there is no error.
func (e Entry) compareIssue(s Signature) error {
	if s.Issue == "" {
		return nil
	}
	if e.Begin.Issue != "" {
		if s.IssueInt() < e.Begin.IssueInt() {
			return ErrBeforeCoverageInterval
		}
	}
	if e.End.Issue != "" {
		if s.IssueInt() > e.End.IssueInt() {
			return ErrAfterCoverageInterval
		}
	}
	return nil
}

// Signature is a bag of information of the record from which coverage can be
// determined. Date should can be a date in with optional year, month, day. The
// volume and issue should be in the best case integers, but sometimes they
// won't.
type Signature struct {
	// Date is often just a year, but sometime also an ISO-8601 date.
	Date   string
	Volume string
	Issue  string
}

// VolumeInt returns the Volume in a best effort manner.
func (s Signature) VolumeInt() int {
	return findInt(s.Volume)
}

// VolumeInt returns the Volume in a best effort manner.
func (s Signature) IssueInt() int {
	return findInt(s.Issue)
}

// findInt return the first int that is found in s or 0 if there is no number.
func findInt(s string) int {
	// we expect to see a number most of the time
	if i, err := strconv.Atoi(s); err == nil {
		return int(i)
	}
	// otherwise try to parse out a number
	if m := intPattern.FindString(s); m == "" {
		return 0
	} else {
		i, _ := strconv.ParseInt(m, 10, 32)
		return int(i)
	}
}
