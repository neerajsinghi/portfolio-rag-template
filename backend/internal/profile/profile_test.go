package profile

import "testing"

// Guards the shipped sample so the template always boots with valid data.
func TestShippedSampleProfileParses(t *testing.T) {
	p, err := Load("../../corpus/profile.json")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name == "" || len(p.Projects) == 0 || len(p.Experience) == 0 {
		t.Fatalf("sample profile incomplete: %+v", p)
	}
}

func TestMissingProfileFallsBackToPlaceholder(t *testing.T) {
	p, err := Load("does-not-exist.json")
	if err == nil || p.Name != Placeholder.Name {
		t.Fatalf("got %+v, %v", p, err)
	}
}
