package collector

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/tprei/semcomp/search/internal/collector/genre"
	"github.com/tprei/semcomp/search/internal/collector/song"
	"github.com/tprei/semcomp/search/pkg/scraper/engine"
	"github.com/tprei/semcomp/search/pkg/utils"
)

func FetchAll(ctx context.Context) ([]song.Song, error) {
	nodes := make([]*engine.Node, 4)
	for i, g := range []genre.Genre{genre.Pop, genre.Rock, genre.Indie, genre.Funk} {
		nodes[i] = engine.NewNode(
			genre.FetchGenre(
				ctx,
				g.String(),
			),
			engine.NewLoggingField("genre", g.String()),
		)
	}

	results := engine.NewScraper(
		ctx,
		nodes,
		engine.Config.FixedAttempts,
		engine.Config.DelayAttempts,
		utils.EmptyFields,
	).Run()

	songs := make([]song.Song, 0, len(results))
	for i := range results {
		ithSongList, ok := results[i].([]song.Song)
		if !ok {
			log.Errorf("failed to cast %#v to []song.Song", results[i])
			return nil, errors.New("failed to cast to song.Song")
		}

		songs = append(songs, ithSongList...)
	}

	return songs, nil
}
