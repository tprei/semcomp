package engine

import (
	"context"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/tprei/semcomp/scraper/pkg/logger"
	"golang.org/x/net/html/charset"
)

var defaultHeaders = map[string]string{
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
}

type leaf struct {
	*http.Request

	ctx    context.Context
	client *http.Client
}

func NewLeaf(ctx context.Context, req *http.Request) *leaf {
	return &leaf{
		ctx:     ctx,
		Request: req,
		client:  http.DefaultClient,
	}
}

func (l *leaf) WithClient(client *http.Client) {
	l.client = client
}

func (l leaf) Do() (*goquery.Document, error) {
	log := logger.GetLogger()
	log = logger.WithFields(log, logrus.Fields{"URL": l.URL.String()})

	for k, v := range defaultHeaders {
		l.Request.Header.Add(k, v)
	}

	resp, err := l.client.Do(l.Request)
	if err != nil {
		log.Error("failed to get leaf node")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("failed to get leaf node")
	}

	reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, http.ErrBodyNotAllowed
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Error("failed to parse html")
		return nil, err
	}

	return doc, nil
}
