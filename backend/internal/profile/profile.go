// Package profile loads the structured facts about the site owner (name,
// links, projects, experience) from corpus/profile.json. These back the
// agent's structured tools; free-form knowledge lives in the corpus Markdown.
package profile

import (
	"encoding/json"
	"fmt"
	"os"
)

// Link is a named contact or social link.
type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// Project is a featured project.
type Project struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Stack       []string `json:"stack"`
	URL         string   `json:"url,omitempty"`
}

// Role is one entry in the work history.
type Role struct {
	Title   string `json:"title"`
	Company string `json:"company"`
	Period  string `json:"period"`
	Summary string `json:"summary"`
}

// Profile is the owner's structured data.
type Profile struct {
	Name         string              `json:"name"`
	Headline     string              `json:"headline"`
	Location     string              `json:"location"`
	Availability string              `json:"availability"`
	Email        string              `json:"email"`
	Links        []Link              `json:"links"`
	Skills       map[string][]string `json:"skills"`
	Projects     []Project           `json:"projects"`
	Experience   []Role              `json:"experience"`
}

// Placeholder is used when no profile.json is found, so the server still runs.
var Placeholder = Profile{Name: "the site owner"}

// Load reads a profile JSON file.
func Load(path string) (Profile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Placeholder, err
	}
	var p Profile
	if err := json.Unmarshal(raw, &p); err != nil {
		return Placeholder, fmt.Errorf("parse %s: %w", path, err)
	}
	if p.Name == "" {
		p.Name = Placeholder.Name
	}
	return p, nil
}
