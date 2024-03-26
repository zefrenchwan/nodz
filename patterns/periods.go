package patterns

import (
	"errors"
	"slices"
	"time"
)

// default comparator for time operations
var periodComparator = NewTimeComparator()

// Period is a set of moments, a moment being a time interval.
// It is neither a duration, nor a set of duration.
// For instance, a person lived in a country from 1999 to 2021 and since 2023.
type Period SeparatedIntervals[time.Time]

// NewPeriod returns a period that contains base exactly
func NewPeriod(base Interval[time.Time]) Period {
	var period Period
	period.elements = []Interval[time.Time]{base}
	return period
}

// NewEmptyPeriod returns an empty period
func NewEmptyPeriod() Period {
	var period Period
	period.elements = make([]Interval[time.Time], 0)
	return period
}

// NewFullPeriod returns a full period
func NewFullPeriod() Period {
	var period Period
	period.elements = []Interval[time.Time]{NewTimeComparator().NewFullInterval()}
	return period
}

// IsEmptyPeriod returns true for an empty period or nil (assumed then to be empty)
func (p *Period) IsEmptyPeriod() bool {
	return p == nil || len(p.elements) == 0 || p.elements[0].IsEmpty()
}

// AsIntervals returns the period as a sorted set of separated intervals
func (p *Period) AsIntervals() []Interval[time.Time] {
	if p == nil {
		return nil
	}

	result := make([]Interval[time.Time], len(p.elements))
	copy(result, p.elements)
	slices.SortFunc(result, periodComparator.CompareInterval)
	return result
}

// AddInterval adds an interval to the period, but ensures invariant that all elements are separated
func (p *Period) AddInterval(i Interval[time.Time]) error {
	if p == nil {
		return errors.New("nil period")
	} else if len(p.elements) == 0 {
		p.elements = []Interval[time.Time]{i}
	} else if i.IsEmpty() || p.elements[0].IsEmpty() {
		return nil
	} else if p.elements[0].IsFull() {
		// already full
		return nil
	}

	// period is not empty, interval to add is not empty
	p.elements = periodComparator.Union(i, p.elements...)

	return nil
}

// Add is the union of periods.
// It returns an error if the receiver is nil
func (p *Period) Add(other Period) error {
	if p == nil {
		return errors.New("nil period")
	} else if other.IsEmptyPeriod() {
		return nil
	}

	currentSize := len(p.elements)
	otherSize := len(other.elements)
	unionOfElements := make([]Interval[time.Time], currentSize+otherSize)

	for index := 0; index < currentSize; index++ {
		unionOfElements[index] = p.elements[index]
	}

	for index := 0; index < otherSize; index++ {
		unionOfElements[currentSize+index] = other.elements[index]
	}

	if len(unionOfElements) >= 2 {
		p.elements = periodComparator.Union(unionOfElements[0], unionOfElements[1:]...)
	} else {
		p.elements = unionOfElements
	}

	return nil
}

// Remove starts with p and remove all the intervals from other.
// Formally, let p_i be the content of p and o_j be the content of other
// New content for p is Union over i of (intersections over j ( p_i minus o_j ))
func (p *Period) Remove(other Period) {
	if p.IsEmptyPeriod() || other.IsEmptyPeriod() {
		return
	}

	var newElements []Interval[time.Time]
	for _, interval := range p.elements {

		var intersections []Interval[time.Time]
		for _, otherInterval := range other.elements {
			differences := periodComparator.Remove(interval, otherInterval)
			if len(differences) != 0 {
				intersections = append(intersections, differences...)
			}
		}

		// For a given interval, that is, for a given i,
		// intersection is the intersection over j of all (p_i minus o_j)
		intersection := periodComparator.Intersection(interval, intersections...)
		if !intersection.IsEmpty() {
			newElements = append(newElements, intersection)
		}
	}

	switch len(newElements) {
	case 0:
		p.elements = []Interval[time.Time]{periodComparator.NewEmptyInterval()}
	case 1:
		p.elements = newElements
	default:
		p.elements = periodComparator.Union(newElements[0], newElements[1:]...)
	}
}
