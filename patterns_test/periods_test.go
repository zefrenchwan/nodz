package patterns_test

import (
	"testing"
	"time"

	"github.com/zefrenchwan/nodz.git/patterns"
)

func TestPeriodAddRemoveIntervals(t *testing.T) {
	comparator := patterns.NewTimeComparator()
	now := time.Now().UTC()
	before := now.AddDate(-1, 0, 0)
	after := now.AddDate(1, 0, 0)

	afterInterval := comparator.NewLeftInfiniteInterval(after, true)
	beforeInterval := comparator.NewLeftInfiniteInterval(before, false)

	pAfter := patterns.NewPeriod(afterInterval)
	pBefore := patterns.NewPeriod(beforeInterval)

	// remove same interval
	pBefore.Remove(pBefore)
	if !pBefore.IsEmptyPeriod() {
		t.Error("period minus itself should be empty")
	}

	pAfter.Remove(pAfter)
	if !pAfter.IsEmptyPeriod() {
		t.Error("period minus itself should be empty")
	}

	// reset
	pAfter = patterns.NewPeriod(afterInterval)
	pBefore = patterns.NewPeriod(beforeInterval)

	// remove when other contains receiver
	pBefore.Remove(pAfter)
	if !pBefore.IsEmptyPeriod() {
		t.Error("before included in after, so before - after should be empty")
	}

	// test when period is larger that removed part
	pAfter = patterns.NewPeriod(afterInterval)
	pBefore = patterns.NewPeriod(beforeInterval)
	pAfter.Remove(pBefore)
	expected := comparator.Remove(afterInterval, beforeInterval)[0]
	result := pAfter.AsIntervals()
	if len(result) != 1 || comparator.CompareInterval(expected, result[0]) != 0 {
		t.Error("failed to remive a single interval in a single interval")
	}
}

func TestPeriodRemoveManyIntervals(t *testing.T) {
	comparator := patterns.NewTimeComparator()
	now := time.Now().UTC()
	before := now.AddDate(-1, 0, 0)
	longAgo := before.AddDate(-2, 0, 0)
	after := now.AddDate(1, 0, 0)
	// longAgo < before < now < after
	longAgoInterval := comparator.NewLeftInfiniteInterval(longAgo, true)
	beforeToNow, _ := comparator.NewFiniteInterval(before, now, false, false)
	nowToAfter, _ := comparator.NewFiniteInterval(now, after, true, true)

	// period is ]-oo, longAgo ] union ]before, now[
	period := patterns.NewPeriod(longAgoInterval)
	period.AddInterval(beforeToNow)
	// otherPeriod is [now, after]
	otherPeriod := patterns.NewPeriod(nowToAfter)
	period.Remove(otherPeriod)
	// period should be the same as before
	result := period.AsIntervals()
	if len(result) != 2 ||
		comparator.CompareInterval(result[0], longAgoInterval) != 0 ||
		comparator.CompareInterval(result[1], beforeToNow) != 0 {
		t.Error("removing other elements should not change current period")
	}
}
