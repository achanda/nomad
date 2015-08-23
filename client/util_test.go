package client

import (
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/nomad/nomad/mock"
	"github.com/hashicorp/nomad/nomad/structs"
)

func TestDiffAllocs(t *testing.T) {
	alloc1 := mock.Alloc() // Ignore
	alloc2 := mock.Alloc() // Update
	alloc2u := new(structs.Allocation)
	*alloc2u = *alloc2
	alloc2u.ModifyIndex += 1
	alloc3 := mock.Alloc() // Remove
	alloc4 := mock.Alloc() // Add

	exist := []*structs.Allocation{
		alloc1,
		alloc2,
		alloc3,
	}
	updated := []*structs.Allocation{
		alloc1,
		alloc2u,
		alloc4,
	}

	result := diffAllocs(exist, updated)

	if len(result.ignore) != 1 || result.ignore[0] != alloc1 {
		t.Fatalf("Bad: %#v", result.ignore)
	}
	if len(result.added) != 1 || result.added[0] != alloc4 {
		t.Fatalf("Bad: %#v", result.added)
	}
	if len(result.removed) != 1 || result.removed[0] != alloc3 {
		t.Fatalf("Bad: %#v", result.removed)
	}
	if len(result.updated) != 1 {
		t.Fatalf("Bad: %#v", result.updated)
	}
	if result.updated[0].exist != alloc2 || result.updated[0].updated != alloc2u {
		t.Fatalf("Bad: %#v", result.updated)
	}
}

func TestGenerateUUID(t *testing.T) {
	prev := generateUUID()
	for i := 0; i < 100; i++ {
		id := generateUUID()
		if prev == id {
			t.Fatalf("Should get a new ID!")
		}

		matched, err := regexp.MatchString(
			"[\\da-f]{8}-[\\da-f]{4}-[\\da-f]{4}-[\\da-f]{4}-[\\da-f]{12}", id)
		if !matched || err != nil {
			t.Fatalf("expected match %s %v %s", id, matched, err)
		}
	}
}

func TestRandomStagger(t *testing.T) {
	intv := time.Minute
	for i := 0; i < 10; i++ {
		stagger := randomStagger(intv)
		if stagger < 0 || stagger >= intv {
			t.Fatalf("Bad: %v", stagger)
		}
	}
}

func TestShuffleStrings(t *testing.T) {
	// Generate input
	inp := make([]string, 10)
	for idx := range inp {
		inp[idx] = generateUUID()
	}

	// Copy the input
	orig := make([]string, len(inp))
	copy(orig, inp)

	// Shuffle
	shuffleStrings(inp)

	// Ensure order is not the same
	if reflect.DeepEqual(inp, orig) {
		t.Fatalf("shuffle failed")
	}
}
