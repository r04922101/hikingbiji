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

// ParseAlbumMainPage parser album main max page number
func ParseAlbumMainPage(body io.Reader) (int, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return 0, errors.Errorf("failed to new goquery document: %v", err)
	}

	maxPage := 0
	doc.Find(".page-item").Each(func(i int, s *goquery.Selection) {
		pageLink, ok := s.Attr("href")
		if !ok {
			return
		}
		url.Parse(pageLink)
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
			maxPage = int(math.Max(float64(maxPage), float64(p)))
		}
	})

	return maxPage, nil
}

// ParseAlbumPage parser an album page to get photo IDs
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
		url.Parse(photoLink)
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
