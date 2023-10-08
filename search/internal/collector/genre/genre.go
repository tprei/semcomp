package genre

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/tprei/semcomp/search/internal/collector/song"
	"github.com/tprei/semcomp/search/pkg/scraper/engine"
	"github.com/tprei/semcomp/search/pkg/utils"
)

const baseURL = "https://www.letras.mus.br"
const genreMask = baseURL + "/mais-acessadas/%s"

type Genre int
type SongRef string

const (
	Pop Genre = iota
	Rock
	Indie
	Funk
)

func (g Genre) String() string {
	switch g {
	case Pop:
		return "pop"
	case Rock:
		return "rock"
	case Indie:
		return "indie"
	case Funk:
		return "funk"
	}

	return "invalid"
}

func parseSongHrefs(document *goquery.Document) (songs []SongRef, parseErr error) {
	songs = make([]SongRef, 0, 1000)
	document.Find("ol.top-list_mus a").Each(func(i int, songNode *goquery.Selection) {
		if songHref, exists := songNode.Attr("href"); exists {
			songs = append(songs, SongRef(songHref))
		} else {
			log.Errorf("could not get href of %v", songNode.Text())
			parseErr = errors.New("failed to get song page URLs")
			return
		}
	})

	if parseErr != nil {
		return nil, parseErr
	}

	return
}

func FetchGenre(ctx context.Context, genre string) func(context.Context) (any, error) {
	return func(_ context.Context) (any, error) {
		URL := fmt.Sprintf(genreMask, genre)
		req, err := http.NewRequest(http.MethodGet, URL, nil)
		if err != nil {
			return nil, err
		}

		leaf := engine.NewLeaf(ctx, req)
		doc, err := leaf.Do()
		if err != nil {
			return nil, err
		}

		hrefs, err := parseSongHrefs(doc)
		if err != nil {
			return nil, err
		}

		nodes := make([]*engine.Node, len(hrefs))
		for i, songHref := range hrefs {
			URL := SongRef(baseURL) + songHref
			nodes[i] = engine.NewNode(
				song.FetchSong(ctx, string(URL)),
				utils.EmptyFields,
			)
		}

		results := engine.NewScraper(
			ctx,
			nodes,
			engine.Config.FixedAttempts,
			engine.Config.DelayAttempts,
			utils.EmptyFields,
		).Run()

		if len(results) != len(hrefs) {
			return nil, errors.New("sth wrong happened, results length differs to songList length")
		}

		songs := make([]song.Song, len(results))
		for i := range results {
			ithSong, ok := results[i].(song.Song)
			if !ok {
				log.Errorf("failed to cast %#v to song.Song", results[i])
				return nil, errors.New("failed to cast to song.Song")
			}

			songs[i] = ithSong
		}

		return songs, nil
	}
}
