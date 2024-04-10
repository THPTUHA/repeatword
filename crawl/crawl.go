package crawl

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Crawler interface {
	Crawl() error
}

type BaseCrawler struct {
	BaseUrl      string
	AudioBaseUrl string
}

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
)

func (b *BaseCrawler) fetchData(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	return doc, err
}

func (cr *BaseCrawler) Crawl() {
}

func (cr *BaseCrawler) saveData() {

}
