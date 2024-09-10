package bluesnews

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var ErrInvalidStatusCode = errors.New("Bluesnews returned a non-200 status code.")

type BluesnewsFetchArticleFn func(date time.Time) (*http.Response, error)

type BluesnewsFetcher struct {
	fetchFn BluesnewsFetchArticleFn
}

func GetBluenewsHTTPResponse(date time.Time) (*http.Response, error) {
	stringDate := fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())
	url := fmt.Sprintf("https://www.bluesnews.com/cgi-bin/blammo.pl?mode=archive&display=%s", stringDate)

	return http.Get(url)
}

func (f BluesnewsFetcher) FetchArticle(date time.Time) (string, error) {
	resp, err := f.fetchFn(date)

	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", ErrInvalidStatusCode
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

var ErrTitleDateNotValid = errors.New("The date in the title is not valid.")
var ErrArticleNotFound = errors.New("The article with the specific date was not found.")

type BluesnewsParser struct{}

func (p *BluesnewsParser) ParseHTML(htmlReader io.Reader) (*Article, error) {
	doc, err := goquery.NewDocumentFromReader(htmlReader)

	if err != nil {
		return nil, err
	}

	titleSelection := doc.Find("h1.pill")
	title := titleSelection.Text()
  pubDate, err := time.Parse("Monday, Jan 02, 2006", extractDateString(title))

	if err != nil {
		return nil, ErrTitleDateNotValid
	}

	content, _ := titleSelection.Next().Html()

	return &Article{Title: title, PubDate: pubDate, ContentHTML: content}, nil
}

func (p *BluesnewsParser) GetHTMLArticleByDate(date time.Time, html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		return "", err
	}

	titleSelection := doc.Find("h1.pill")

	for i := range titleSelection.Nodes {
		sel := titleSelection.Eq(i)
		title := sel.Text()
		pubDate, err := time.Parse("Monday, Jan 02, 2006", extractDateString(title))

		if err != nil {
			continue
		}

		if pubDate.Format("2006-01-02") == date.Format("2006-01-02") {
			titleHtml, _ := goquery.OuterHtml(sel)
			contentHtml, _ := goquery.OuterHtml(sel.Next())
			return titleHtml + contentHtml, nil
		}
	}

	return "", ErrArticleNotFound

}

// This function will only works with title like this example:
// Input: Saturday, Sep 07, 2024 Some day Blablabla
// Output: Saturday, Sep 07, 2024
func extractDateString(s string) string {
	parts := strings.SplitN(s, ",", 3)
	if len(parts) < 3 {
		return s
	}

	yearCleaned := strings.Split(parts[2], " ")[1]
	return fmt.Sprintf("%s,%s, %s", parts[0], parts[1], yearCleaned)
}

type BluesnewsClient struct {
	fetcher BluesnewsFetcher
}

func (bc BluesnewsClient) GetArticleFromDate(date time.Time) (Article, error) {
	content, err := bc.fetcher.FetchArticle(date)

	if err != nil {
		return Article{}, err
	}

	return Article{Title: content, PubDate: date}, nil
}
