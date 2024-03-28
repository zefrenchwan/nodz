package patterns

import (
	"slices"
	"testing"
	"time"

	"github.com/zefrenchwan/nodz.git/patterns"
)

func TestInstanceNoInstantiatingClass(t *testing.T) {
	if _, err := patterns.NewFormalInstance(nil); err == nil {
		t.Fail()
	}
}

func TestInstanceNoAttribute(t *testing.T) {
	currentClass := patterns.NewFormalClass("test")
	instance, _ := patterns.NewFormalInstance(&currentClass)

	if err := instance.SetValue("attr", "value"); err == nil {
		t.Fail()
	}

	if v, err := instance.Attributes(); v != nil || err != nil {
		t.Fail()
	}
}

func TestInstanceSetAttribute(t *testing.T) {
	currentClass := patterns.NewFormalClass("test")
	currentClass.AddAttribute("attr")
	currentClass.AddAttribute("otherAttr")

	instance, _ := patterns.NewFormalInstance(&currentClass)

	if attr, err := instance.Attributes(); err != nil {
		t.Fail()
	} else if attr == nil {
		t.Fail()
	} else if slices.Compare(attr, []string{"attr", "otherAttr"}) != 0 {
		t.Fail()
	}

	if err := instance.SetValue("i don't exist", "value"); err == nil {
		t.Fail()
	}

	if err := instance.SetValue("attr", "a value"); err != nil {
		t.Fail()
	}

	if v, err := instance.ValuesForAttribute("otherAttr"); v != nil || err != nil {
		t.Error("failed test when loading values in unused attribute")
	}

	instance.SetValue("attr", "final value")

	if values, err := instance.ValuesForAttribute("attr"); values == nil || err != nil {
		t.Error("no value found for existing attribute")
	} else if len(values) != 1 {
		t.Error("error when finding existing attribute")
	} else if values[0] != "final value" {
		t.Error("error when finding value")
	}

	// final value should be set for full period
	if values, err := instance.TimeValuesForAttribute("attr"); err != nil {
		t.Fail()
	} else if len(values) != 1 {
		t.Error("map for attributes error")
	} else if !values["final value"][0].IsFull() {
		t.Error("setting value should make full")
	}
}

func TestInstanceAddAttribute(t *testing.T) {
	currentClass := patterns.NewFormalClass("test")
	currentClass.AddAttribute("attr")

	now := time.Now().UTC()
	beforeNow := patterns.NewLeftInfiniteTimeInterval(now, false)
	afterNow := patterns.NewRightInfiniteTimeInterval(now, true)

	instance, _ := patterns.NewFormalInstance(&currentClass)
	instance.AddValue("attr", "before", patterns.NewPeriod(beforeNow))
	instance.AddValue("attr", "after", patterns.NewPeriod(afterNow))

	// test values, not periods
	var values []string
	if v, err := instance.ValuesForAttribute("attr"); err != nil {
		t.Fail()
	} else if len(v) != 2 {
		t.Error("missing values when many values")
	} else {
		slices.Sort(v)
		values = v
	}

	if slices.Compare([]string{"after", "before"}, values) != 0 {
		t.Error("no match for many values")
	}

	// test periods
	if valuesMap, err := instance.TimeValuesForAttribute("attr"); err != nil {
		t.Fail()
	} else if len(valuesMap) != 2 {
		t.Error("missing values in map of values")
	} else if beforeValue := valuesMap["before"]; len(beforeValue) != 1 {
		t.Error("intervals test failed")
	} else if patterns.TimeIntervalsCompare(beforeNow, beforeValue[0]) != 0 {
		t.Error("intervals test failed")
	} else if afterValue := valuesMap["after"]; len(afterValue) != 1 {
		t.Error("intervals test failed")
	} else if patterns.TimeIntervalsCompare(afterValue[0], afterNow) != 0 {
		t.Error("intervals test failed")
	}
}

func TestInstancePeriodChange(t *testing.T) {
	currentClass := patterns.NewFormalClass("test")
	currentClass.AddAttribute("attr")

	now := time.Now().UTC()
	beforeNow := patterns.NewLeftInfiniteTimeInterval(now, false)
	afterNow := patterns.NewRightInfiniteTimeInterval(now, true)

	instance, _ := patterns.NewFormalInstance(&currentClass)
	instance.SetValue("attr", "before")
	instance.AddValue("attr", "after", patterns.NewPeriod(afterNow))

	// test values, not periods
	var values []string
	if v, err := instance.ValuesForAttribute("attr"); err != nil {
		t.Fail()
	} else if len(v) != 2 {
		t.Error("missing values when many values")
	} else {
		slices.Sort(v)
		values = v
	}

	if slices.Compare([]string{"after", "before"}, values) != 0 {
		t.Error("no match for many values")
	}

	// test periods
	if valuesMap, err := instance.TimeValuesForAttribute("attr"); err != nil {
		t.Fail()
	} else if len(valuesMap) != 2 {
		t.Error("missing values in map of values")
	} else if beforeValue := valuesMap["before"]; len(beforeValue) != 1 {
		t.Error("intervals test failed")
	} else if patterns.TimeIntervalsCompare(beforeNow, beforeValue[0]) != 0 {
		t.Error("intervals test failed")
	} else if afterValue := valuesMap["after"]; len(afterValue) != 1 {
		t.Error("intervals test failed")
	} else if patterns.TimeIntervalsCompare(afterValue[0], afterNow) != 0 {
		t.Error("intervals test failed")
	}
}
