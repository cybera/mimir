package datasets

import (
	"path/filepath"
	"strings"
)

// CanonicalName returns the cleaned version of a dataset name, stripped of
// any extraneous parts like file extension.
func canonicalName(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

// CanonicalNames returns the cleaned versions of a list of dataset names.
func canonicalNames(names []string) []string {
	var canonicalNames []string

	for _, name := range names {
		canonicalNames = append(canonicalNames, canonicalName(name))
	}

	return canonicalNames
}
