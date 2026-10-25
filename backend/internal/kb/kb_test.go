package kb

import "testing"

func TestLoadSplitsSectionsAndRetrieves(t *testing.T) {
	ix, err := Load("testdata")
	if err != nil {
		t.Fatal(err)
	}
	// testdata/sample.md has an intro + 2 sections; README.md is skipped.
	if ix.Len() != 3 {
		t.Fatalf("Len = %d, want 3", ix.Len())
	}
	got := ix.Retrieve("kubernetes", 1)
	if len(got) != 1 || got[0].ID != "sample#infrastructure" {
		t.Fatalf("Retrieve(kubernetes) = %#v", got)
	}
	// Tags are searchable even though they are not part of the text.
	got = ix.Retrieve("gardening", 1)
	if len(got) != 1 || got[0].Title != "Hobbies" {
		t.Fatalf("Retrieve(gardening) = %#v", got)
	}
}

func TestNilIndexIsSafe(t *testing.T) {
	var ix *Index
	if ix.Len() != 0 || len(ix.Retrieve("x", 3)) != 0 {
		t.Fatal("nil index should be empty")
	}
}

func TestShippedSampleCorpusLoads(t *testing.T) {
	ix, err := Load("../../corpus")
	if err != nil {
		t.Fatal(err)
	}
	if ix.Len() == 0 || len(ix.Retrieve("experience", 1)) == 0 {
		t.Fatalf("sample corpus should load and be searchable (chunks=%d)", ix.Len())
	}
}
