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

type BluesnewsParser struct{}

func (p *BluesnewsParser) ParseHTML(htmlReader io.Reader) (*Article, error) {
	doc, err := goquery.NewDocumentFromReader(htmlReader)

	if err != nil {
		return nil, err
	}

	titleSelection := doc.Find("h1.pill")
	dateSelection := titleSelection.Children().Eq(0)
	title := titleSelection.Text()
	dateStr, _ := strings.CutSuffix(title, dateSelection.Text())
	dateStr = strings.TrimSpace(dateStr)
	pubDate, err := time.Parse("Monday, Jan 02, 2006", dateStr)

	if err != nil {
		return nil, ErrTitleDateNotValid
	}

	content, _ := titleSelection.Next().Html()

	return &Article{Title: title, PubDate: pubDate, ContentHTML: content}, nil
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
