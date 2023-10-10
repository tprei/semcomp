package song

import (
	"context"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/tprei/semcomp/scraper/pkg/scraper/engine"
)

type Verse string
type Strophe []Verse
type Lyrics []Strophe

type Song struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Lyrics `json:"lyrics"`
	Avatar string `json:"img"`
	Views  string `json:"views"`
}

func parseArtistName(document *goquery.Document) string {
	return document.Find("h1.head-title").Text()
}

func parseSongTitle(document *goquery.Document) string {
	title := document.Find("h2.head-subtitle").Text()
	return strings.TrimSpace(title)
}

func parseArtistAvatar(document *goquery.Document) string {
	const defaultAvatarIcon = "https://akamai.sscdn.co/letras/desktop/static/img/ic_placeholder_artist.svg"
	icon := document.Find("img.head-image").AttrOr("src", defaultAvatarIcon)
	return strings.TrimSpace(icon)
}

func parseSongViews(document *goquery.Document) string {
	views := document.Find("div.head-info-exib > b").Text()
	return strings.TrimSpace(views)
}

func parseLyrics(document *goquery.Document) (lyrics Lyrics, err error) {
	strophes := make([]Strophe, 0)
	document.Find("div.lyric-original").Children().Each(func(i int, stropheNode *goquery.Selection) {
		selection := stropheNode.Not(".lyricAnnotation")

		var html string
		html, err = selection.Html()
		if err != nil {
			return
		}

		verses := strings.Split(html, "<br/>")
		strophe := make([]Verse, len(verses))
		for i := range verses {
			strophe[i] = Verse(verses[i])
		}

		strophes = append(strophes, strophe)
	})

	lyrics = Lyrics(strophes)
	return
}

func FetchSong(ctx context.Context, URL string) func(context.Context) (any, error) {
	return func(_ context.Context) (any, error) {
		req, err := http.NewRequest(http.MethodGet, URL, nil)
		if err != nil {
			return nil, err
		}

		leaf := engine.NewLeaf(ctx, req)
		doc, err := leaf.Do()
		if err != nil {
			return nil, err
		}

		artist := parseArtistName(doc)
		songTitle := parseSongTitle(doc)
		avatar := parseArtistAvatar(doc)
		views := parseSongViews(doc)

		if artist == "" || songTitle == "" {
			log.Warnf("got empty artist or song name in URL %s", URL)
		}

		lyrics, err := parseLyrics(doc)
		if err != nil {
			return nil, err
		}

		return Song{
			Title:  songTitle,
			Artist: artist,
			Lyrics: lyrics,
			Avatar: avatar,
			Views:  views,
		}, nil
	}
}
