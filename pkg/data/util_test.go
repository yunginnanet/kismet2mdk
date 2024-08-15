package data

import (
	"slices"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestGatherSources(t *testing.T) {
	sources := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
	}

	groupedSources, err := gatherSources(sources...)
	if err != nil {
		t.Fatal(err)
	}

	if len(groupedSources) != 2 {
		t.Fatalf("Expected 2 groups of sources but got %d", len(groupedSources))
	}

	spew.Dump(groupedSources)

	groups := make([][]string, 2)

	for _, g := range []int{0, 1} {
		groups[g] = slices.DeleteFunc(groupedSources[g][:], func(val string) bool {
			return val == ""
		})
		if len(groups[g]) != 10 {
			spew.Dump(groups[g])
			t.Fatalf("Expected 10 sources in group but got %d", len(groups[g]))
		}
	}
}
