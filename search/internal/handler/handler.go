package handler

import "github.com/tprei/semcomp/search/internal/database/json"

type Song struct{}
type Lyrics struct{}
type Strophe struct{}

const songsFile = `songs.json`

func filterSongs(songs []Song, query string) []Song {
	return nil
}

func Search(query string) ([]Song, error) {
	songs := make([]Song, 0)
	err := json.ReadSongs(songsFile, &songs)
	if err != nil {
		return nil, err
	}

	return filterSongs(songs, query), nil
}
