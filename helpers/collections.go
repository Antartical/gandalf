package helpers

import (
	"github.com/lib/pq"
)

// Check if the given element is present on the given array
func PqStringArrayContains(pqArray pq.StringArray, element interface{}) bool {
	for _, accessedElement := range pqArray {
		if accessedElement == element {
			return true
		}
	}
	return false
}
