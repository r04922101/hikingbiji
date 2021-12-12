package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/r04922101/hikingbiji/httpext"
	"github.com/r04922101/hikingbiji/parser"
)

const (
	bijiHost       = "hiking.biji.co"
	bijiPath       = "index.php"
	listAlbumQuery = "q=album&act=photo_list&album_id=%s"
	clapEndpoint   = "https://hiking.biji.co/album/ajax/clap_photo"
)

var (
	cookie = flag.String("cookie", "", "cookie to represent a login user")
	album  = flag.String("album", "", "album ID")
)

func init() {
	flag.Parse()
}

func clapAlbum(ctx context.Context, httpClient *http.Client, albumID string) ([]string, error) {
	// do HTTP GET to get album main page
	albumURL := &url.URL{
		Scheme:   "https",
		Host:     bijiHost,
		Path:     bijiPath,
		RawQuery: fmt.Sprintf(listAlbumQuery, *album),
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, albumURL.String(), nil)
	if err != nil {
		return nil, errors.Errorf("failed to new GET request to %s with context: %v", albumURL, err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Errorf("failed to do HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// parse album main page to get pages
	maxPage, err := parser.ParseAlbumMainPage(resp.Body)
	if err != nil {
		return nil, errors.Errorf("failed to parse album %s main page: %v", albumID, err)
	} else if maxPage == 0 {
		return nil, errors.Errorf("invalid album max page")
	}

	// clap album pages concurrently
	errorgroup, gctx := errgroup.WithContext(ctx)
	for page := 1; page <= maxPage; page++ {
		pageQuery := albumURL.Query()
		pageQuery.Add("page", fmt.Sprint(page))
		albumPageURL := &url.URL{
			Scheme:   "https",
			Host:     bijiHost,
			Path:     bijiPath,
			RawQuery: pageQuery.Encode(),
		}
		page := page
		errorgroup.Go(func() error {
			logrus.Infof("start clapping page %d", page)
			if err := clapAlbumPage(gctx, httpClient, albumPageURL, albumID); err != nil {
				logrus.Warnf("failed to clap alubm page %d: %v", page, err)
			}
			logrus.Infof("finish clapping page %d", page)
			return nil
		})
	}
	errorgroup.Wait()

	return nil, nil
}

func clapAlbumPage(ctx context.Context, httpClient *http.Client, albumPageURL *url.URL, albumID string) error {
	// do HTTP GET to get album page
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, albumPageURL.String(), nil)
	if err != nil {
		return errors.Errorf("failed to new GET request to %s with context: %v", albumPageURL, err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.Errorf("failed to do HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// parse album page to get photo IDs
	photoIDs, err := parser.ParseAlbumPage(resp.Body)
	if err != nil {
		return errors.Errorf("failed to parse album page %s: %v", albumPageURL)
	}
	// clap photos sequentially to avoid network congestion
	// and too many requests to biji server at the same time
	for _, photoID := range photoIDs {
		if err := clapPhoto(ctx, httpClient, albumID, photoID); err != nil {
			logrus.Warnf("failed to clap photo %s: %v", photoID, err)
		}
	}

	return nil
}

func clapPhoto(ctx context.Context, httpClient *http.Client, albumID, photoID string) error {
	req, err := newClapPhotoRequest(ctx, albumID, photoID, *cookie)
	if err != nil {
		return errors.Errorf("failed to new clap photo %s request: %v", photoID, err)
	}

	if resp, err := httpClient.Do(req); err != nil {
		return errors.Errorf("failed to clap photo %s with error: %v", photoID, err)
	} else if statusCode := resp.StatusCode; statusCode != http.StatusOK {
		return errors.Errorf("failed to clap photo %s with status Code: %v", statusCode, err)
	}

	return nil
}

func newClapPhotoRequest(ctx context.Context, albumID, photoID, cookie string) (*http.Request, error) {
	payload := strings.NewReader(fmt.Sprintf(`{"id": "%s"}`, photoID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, clapEndpoint, payload)
	if err != nil {
		return nil, errors.Errorf("failed to new POST request to %s: %v", clapEndpoint, err)
	}

	req.Header.Add("authority", "hiking.biji.co")
	req.Header.Add("sec-ch-ua", "\"Microsoft Edge\";v=\"95\", \"Chromium\";v=\"95\", \";Not A Brand\";v=\"99\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36 Edg/95.0.1020.30")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("content-type", "text/plain;charset=UTF-8")
	req.Header.Add("accept", "*/*")
	req.Header.Add("origin", "https://hiking.biji.co")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", fmt.Sprintf("https://hiking.biji.co/index.php?q=album&act=photo&album_id=%s&ap_id=%s", albumID, photoID))
	req.Header.Add("accept-language", "en-US,en;q=0.9,zh-TW;q=0.8,zh;q=0.7")
	req.Header.Add("cookie", cookie)

	return req, nil
}

func main() {
	ctx := context.Background()
	httpClient := httpext.NewHTTPClient()

	logrus.Infof("start clapping album %s", *album)
	if _, err := clapAlbum(ctx, httpClient, *album); err != nil {
		logrus.Fatalf("failed to clap album %s: %v", *album, err)
	}
	logrus.Infof("finish clapping album %s!", *album)
}
