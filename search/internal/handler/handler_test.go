package handler

import (
	"testing"
)

type isInsideTestCase struct {
	verse    []string
	query    []string
	expected bool
}

type filterTestCase struct {
	songs []Song
	query string

	expected []Song
}

func TestIsInsideVerse(t *testing.T) {
	cases := []isInsideTestCase{
		{
			verse:    []string{"Psychic", "spies", "from", "China"},
			query:    []string{"Psychic", "spies"},
			expected: true,
		},
		{
			verse:    []string{"Psychic", "spies", "from", "China"},
			query:    []string{"Psychic", "from"},
			expected: false,
		},
	}

	for _, c := range cases {
		if got := isInsideVerse(c.verse, c.query); got != c.expected {
			t.Fatalf("expected %v, got %v", c.expected, got)
		}
	}
}

func TestFilterSongs(t *testing.T) {
	cases := []filterTestCase{
		{
			songs: []Song{
				{
					Lyrics: []Strophe{
						[]string{
							"Mas o seu sorriso vale mais que um diamante",
							"Se você vier comigo, aí nós vamos adiante",
						},
					},
				},
			},
			query: "você vier comigo",
			expected: []Song{
				{
					Lyrics: []Strophe{
						[]string{
							"Mas o seu sorriso vale mais que um diamante",
							"Se você vier comigo, aí nós vamos adiante",
						},
					},
				},
			},
		},
		{
			songs: []Song{
				{
					Lyrics: []Strophe{
						[]string{
							"It's understood that Hollywood sells Californication",
						},
					},
				},
			},
			query: "Californication",
			expected: []Song{
				{
					Lyrics: []Strophe{
						[]string{
							"It's understood that Hollywood sells Californication",
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		if got := filterSongs(c.songs, c.query); len(got) != len(c.expected) {
			t.Fatalf("got len: %v, expected len %v", len(got), len(c.expected))
		}
	}
}
