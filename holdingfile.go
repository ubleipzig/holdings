package holdingfile

import "time"

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

// Signature is a bag of information of the record from which coverage can be
// determined. Year should be four digits, volume and issue should be in the
// best case integers, but sometimes they won't.
type Signature struct {
	Year   string
	Volume string
	Issue  string
}

type ISSN string
