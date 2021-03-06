package parser

import (
	"io"
	"math"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

const (
	// albumSize defines a constant for how many photos should be in an album page
	albumSize = 24
)

// ParseAlbumMainPage parses album main max page number
func ParseAlbumMainPage(body io.Reader) (int64, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return 0, errors.Errorf("failed to new goquery document: %v", err)
	}

	var maxPage int64
	doc.Find(".page-item").Each(func(i int, s *goquery.Selection) {
		pageLink, ok := s.Attr("href")
		if !ok {
			// div.page-item.at-this-page
			p, _ := strconv.ParseInt(s.Text(), 10, 64)
			maxPage = int64(math.Max(float64(maxPage), float64(p)))
			return
		}

		u, err := url.Parse(pageLink)
		if err != nil {
			err = errors.Errorf("failed to parse link %s: %v", pageLink, err)
			return
		}
		vs, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			err = errors.Errorf("failed to parse query %s: %v", u.RawQuery, err)
			return
		}
		if page := vs.Get("page"); page != "" {
			p, _ := strconv.ParseInt(page, 10, 64)
			maxPage = int64(math.Max(float64(maxPage), float64(p)))
		}
	})

	return maxPage, nil
}

// ParseAlbumPage parses an album page to get photo IDs
func ParseAlbumPage(body io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, errors.Errorf("failed to new goquery document: %v", err)
	}

	photoIDs := make([]string, 0, 24)
	doc.Find("a.postMeta-img.img-container").Each(func(i int, s *goquery.Selection) {
		photoLink, ok := s.Attr("href")
		if !ok {
			return
		}
		u, err := url.Parse(photoLink)
		if err != nil {
			err = errors.Errorf("failed to parse link %s: %v", photoLink, err)
			return
		}
		vs, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			err = errors.Errorf("failed to parse query %s: %v", u.RawQuery, err)
			return
		}
		if photoID := vs.Get("ap_id"); photoID != "" {
			photoIDs = append(photoIDs, photoID)
		}
	})
	if err != nil {
		return nil, errors.Errorf("failed to find photo: %v", err)
	}

	return photoIDs, nil
}
