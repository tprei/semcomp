package handler

import (
	"strings"

	"github.com/tprei/semcomp/search/internal/database/json"
)

type Song struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`

	Image string `json:"img"`
	Views string `json:"views"`

	Lyrics `json:"lyrics"`
}

type Lyrics []Strophe

type Strophe []string

const songsFile = `songs.json`

func filterSongs(songs []Song, query string) []Song {
	filtered := make([]Song, 0, len(songs))

songs:
	for _, song := range songs {
		for _, strophes := range song.Lyrics {
			for _, verse := range strophes {
				lowerVerse := strings.ToLower(verse)
				lowerQuery := strings.ToLower(query)

				if strings.Contains(lowerVerse, lowerQuery) {
					filtered = append(filtered, song)
					continue songs
				}
			}
		}
	}

	return filtered
}

func Search(query string) ([]Song, error) {
	songs := make([]Song, 0)
	err := json.ReadSongs(songsFile, &songs)
	if err != nil {
		return nil, err
	}

	return filterSongs(songs, query), nil
}
