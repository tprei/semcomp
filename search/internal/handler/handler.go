package handler

import "github.com/tprei/semcomp/search/internal/database/json"

type Song struct{}
type Lyrics struct{}
type Strophe struct{}

const songsFile = `songs.json`

func parseSongMap(map[string]string) []Song {
	return nil
}

func filterSongs(songs []Song, query string) []Song {
	return nil
}

func Search(query string) ([]Song, error) {
	songsMap, err := json.ReadSongs(songsFile)
	if err != nil {
		return nil, err
	}

	songs := parseSongMap(songsMap)
	return filterSongs(songs, query), nil
}
