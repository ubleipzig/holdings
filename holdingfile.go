package holdingfile

import (
	"errors"
	"time"
)

var (
	ErrBeforeCoverageInterval = errors.New("before coverage interval")
	ErrAfterCoverageInterval  = errors.New("after coverage interval")
	ErrMissingValues          = errors.New("missing values")
)

// Holdings can return a list of licenses for a given ISSN.
type Holdings interface {
	Licenses(ISSN) []License
}

// License exposes methods to let clients check their own validity in the
// context of this license.
type License interface {
	// Covers looks at the static content only.
	Covers(Signature) error
	// TimeRestricted will report an error, if a moving wall constraint holds.
	TimeRestricted(time.Time) error
}

// Entry is a reduced holding file entry.
type Entry struct {
	Begin   Signature
	End     Signature
	Embargo time.Duration
}

func (e Entry) Covers(s Signature) error {
	if err := e.compareYear(s); err != nil {
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

// compareYear returns an error, if both values are defined and disagree,
// otherwise we assume there is no error.
func (e Entry) compareDate(s Signature) error {
	if s.Date == "" {
		return nil
	}
	if e.Begin.Date != "" {
		if s.Date < e.Begin.Date {
			return ErrBeforeCoverageInterval
		}
	}
	if e.End.Year != "" {
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
		if s.Volume < e.Begin.Volume {
			return ErrBeforeCoverageInterval
		}
	}
	if e.End.Volume != "" {
		if s.Volume > e.End.Volume {
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
		if s.Issue < e.Begin.Issue {
			return ErrBeforeCoverageInterval
		}
	}
	if e.End.Issue != "" {
		if s.Issue > e.End.Issue {
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
	Date   string
	Volume string
	Issue  string
}

type ISSN string
