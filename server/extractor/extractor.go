// Package extractor provides HTML content extraction for documents.
package extractor

import (
	"bytes"
	"errors"
	"io"
	"net/url"
	"strings"

	readability "codeberg.org/readeck/go-readability/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"

	"github.com/asciimoo/hister/server/document"
)

// Extractor extracts content from a Document.
type Extractor interface {
	Name() string
	Match(*document.Document) bool
	Extract(*document.Document) error
	Preview(*document.Document) (string, error)
}

// ErrNoExtractor is returned when no extractor can handle the document.
var ErrNoExtractor = errors.New("no extractor found")

var extractors = []Extractor{
	&readabilityExtractor{},
	&defaultExtractor{},
}

// Extract tries each registered extractor in order and returns the first
// successful result. Returns ErrNoExtractor if none succeed.
func Extract(d *document.Document) error {
	for _, e := range extractors {
		if e.Match(d) {
			if err := e.Extract(d); err != nil {
				log.Warn().Err(err).Str("URL", d.URL).Str("Extractor", e.Name()).Msg("Failed to extract content")
			} else {
				return nil
			}
		}
	}
	return ErrNoExtractor
}

// Preview returns a rendered preview of the document using the first matching
// extractor. Returns ErrNoExtractor if none match.
func Preview(d *document.Document) (string, error) {
	for _, e := range extractors {
		if e.Match(d) {
			content, err := e.Preview(d)
			if err != nil {
				log.Warn().Err(err).Str("URL", d.URL).Str("Extractor", e.Name()).Msg("Failed to preview content")
			} else {
				return content, nil
			}
		}
	}
	return "", ErrNoExtractor
}

type defaultExtractor struct{}

type readabilityExtractor struct{}

func (e *defaultExtractor) Name() string {
	return "Default"
}

func (e *defaultExtractor) Match(_ *document.Document) bool {
	return true
}

func (e *defaultExtractor) Extract(d *document.Document) error {
	d.Title = ""
	r := bytes.NewReader([]byte(d.HTML))
	doc := html.NewTokenizer(r)
	inBody := false
	skip := false
	var text strings.Builder
	var currentTag string
out:
	for {
		tt := doc.Next()
		switch tt {
		case html.ErrorToken:
			err := doc.Err()
			if errors.Is(err, io.EOF) {
				break out
			}
			return errors.New("failed to parse html: " + err.Error())
		case html.SelfClosingTagToken, html.StartTagToken:
			tn, _ := doc.TagName()
			currentTag = string(tn)
			switch currentTag {
			case "body":
				inBody = true
			case "script", "style", "noscript":
				skip = true
			}
		case html.TextToken:
			if currentTag == "title" {
				d.Title += strings.TrimSpace(string(doc.Text()))
			}
			if inBody && !skip {
				text.Write(doc.Text())
			}
		case html.EndTagToken:
			tn, _ := doc.TagName()
			switch string(tn) {
			case "body":
				inBody = false
			case "script", "style", "noscript":
				skip = false
			}
		}
	}
	d.Text = strings.TrimSpace(text.String())
	if d.Text == "" && d.Title == "" {
		return errors.New("no content found")
	}
	return nil
}

func (e *defaultExtractor) Preview(d *document.Document) (string, error) {
	return d.Text, nil
}

func (e *readabilityExtractor) Name() string {
	return "Readability"
}

func (e *readabilityExtractor) Match(_ *document.Document) bool {
	return true
}

func (e *readabilityExtractor) Extract(d *document.Document) error {
	r := bytes.NewReader([]byte(d.HTML))

	u, err := url.Parse(d.URL)
	if err != nil {
		return err
	}
	a, err := readability.FromReader(r, u)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	if err := a.RenderText(buf); err != nil {
		return err
	}
	d.Text = buf.String()
	d.Title = a.Title()
	d.SetFaviconURL(a.Favicon())
	return nil
}

func (e *readabilityExtractor) Preview(d *document.Document) (string, error) {
	r := bytes.NewReader([]byte(d.HTML))
	u, err := url.Parse(d.URL)
	if err != nil {
		return "", err
	}
	a, err := readability.FromReader(r, u)
	if err != nil {
		return "", err
	}
	var htmlContent strings.Builder
	if err := a.RenderHTML(&htmlContent); err != nil {
		return "", err
	}
	return htmlContent.String(), nil
}
