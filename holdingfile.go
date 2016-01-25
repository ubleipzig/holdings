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
	if err := e.compareIsse(s); err != nil {
		return err
	}
	return nil
}

// compareYear returns an error, if both values are defined and disagree,
// otherwise we assume there is no error.
func (e Entry) compareYear(s Signature) error {
	if e.Begin.Year != "" && s.Year != "" {
		if s.Year < e.Begin.Year {
			return ErrBeforeCoverageInterval
		}
	}
	if e.End.Year != "" && s.Year != "" {
		if s.Year > e.End.Year {
			return ErrAfterCoverageInterval
		}
	}
	return nil
}

// compareVolume returns an error, if both values are defined and disagree,
// otherwise we assume there is no error.
func (e Entry) compareVolume(s Signature) error {
	if e.Begin.Volume == "" || s.Volume == "" {
		return ErrMissingValues
	}

	if s.Volume < e.Begin.Volume {
		return ErrBeforeCoverageInterval
	} else {
		if s.Volume > e.End.Volume {
			return ErrAfterCoverageInterval
		}
	}
	return nil
}

// compareIssue returns an error, if both values are defined and disagree,
// otherwise we assume there is no error.
func (e Entry) compareIssue(s Signature) error {
	if e.Begin.Issue != "" && s.Issue != "" {
		if s.Issue < e.Begin.Issue {
			return ErrBeforeCoverageInterval
		}
	}
	if e.End.Issue != "" && s.Issue != "" {
		if s.Issue > e.End.Issue {
			return ErrAfterCoverageInterval
		}
	}
	return nil
}

// Signature is a bag of information of the record from which coverage can be
// determined. Year should be four digits, volume and issue should be in the
// best case integers, but sometimes they won't.
type Signature struct {
	Year   string
	Volume string
	Issue  string
}

type ISSN string
