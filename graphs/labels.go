package graphs

import (
	"slices"
	"strings"
)

// WithLabels defines basic operation to deal with labels
type WithLabels interface {
	// AddLabel appends a label to the set of labels (no duplicate)
	AddLabel(label string)
	// Labels returns the labels as a sorted slice and nil for nil
	// Sorting allows to compare labels using the order.
	Labels() []string
	// RemoveLabel removes, if any, the label key
	RemoveLabel(key string)
}

// JoinLabels should be the standard way to serialize labels to a string
func JoinLabels(labels WithLabels) string {
	currentLabels := labels.Labels()
	copyLabels := make([]string, len(currentLabels))
	copy(copyLabels, currentLabels)

	slices.Sort(copyLabels)
	return strings.Join(copyLabels, ",")
}
