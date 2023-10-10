package json

import (
	"encoding/json"
	"os"
)

// ReadSongs takes a JSON filename, parses it and returns it as a map[string]string
func ReadSongs(filename string) (map[string]string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var songs map[string]string
	if err := json.Unmarshal(bytes, &songs); err != nil {
		return nil, err
	}

	return songs, nil
}
