package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/tprei/semcomp/scraper/internal/collector"
)

func main() {
	logrus.SetLevel(log.DebugLevel)
	songs, err := collector.FetchAll(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	bytes, err := json.MarshalIndent(songs, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	// save file
	if err := os.WriteFile("songs.json", bytes, 0666); err != nil {
		log.Fatal(err)
	}
}
