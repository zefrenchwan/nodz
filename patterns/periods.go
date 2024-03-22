package patterns

import "time"

// Period is a set of time intervals.
// It means the moments (time instances) a time dependant property is true.
// Typical use would be:
// given a person p and a country c, a period would be the time spent by p in c.
// It may be a couple of months every year, for instance.
type Period SeparatedIntervals[time.Time]
