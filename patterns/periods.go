package patterns

import "time"

// Period is a set of moments, a moment being a time interval.
// It is neither a duration, nor a set of duration.
// For instance, a person lived in a country from 1999 to 2021 and since 2023.
type Period SeparatedIntervals[time.Time]

// NewFullPeriod returns a full period
func NewFullPeriod() Period {
	var period Period
	period.elements = []Interval[time.Time]{NewTimeComparator().NewFullInterval()}
	return period
}
