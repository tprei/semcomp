package json

import (
	"encoding/json"
	"os"
)

// ReadSongs takes a JSON filename, parses it and unmarshalls it onto obj
//
// obj must be of type pointer, otherwise this functon returns InvalidUnmarshalError
func ReadSongs(filename string, obj any) error {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, obj); err != nil {
		return err
	}

	return nil
}
