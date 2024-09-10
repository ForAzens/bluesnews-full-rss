package bluesnews

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBluesnewsClient(t *testing.T) {
	fetcher := createSuccessFetchFn()
	client := BluesnewsClient{fetcher: fetcher}
	publishDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	got, _ := client.GetArticleFromDate(publishDate)
	want := Article{Title: "2023-12-31", PubDate: publishDate}

	if got != want {
		t.Errorf("got %+v want %+v", got, want)
	}
}

func TestBluesnewsFetcher(t *testing.T) {
	t.Run("returns no error", func(t *testing.T) {
		fetcher := createSuccessFetchFn()

		publishDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
		_, err := fetcher.FetchArticle(publishDate)

		assertNoError(t, err)
	})

	t.Run("returns error on 400 status", func(t *testing.T) {
		fetcher := createFetchFnWithStatusCode(400)

		publishDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
		got, err := fetcher.FetchArticle(publishDate)
		want := ""

		assertError(t, err, ErrInvalidStatusCode)
		if got != want {
			t.Errorf("content got %s, content wanted %s", got, want)
		}
	})

	t.Run("returns fetchFn error", func(t *testing.T) {
		timeoutError := errors.New("Timeout error")
		fetcher := createFetchFnWithError(timeoutError)

		publishDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
		got, err := fetcher.FetchArticle(publishDate)
		want := ""

		assertError(t, err, timeoutError)
		if got != want {
			t.Errorf("content got %s, content wanted %s", got, want)
		}
	})
}

type MockErrReader struct{ Error error }

func (e MockErrReader) Read(buf []byte) (int, error) {
	return 0, e.Error
}

func TestBluesnewsParser(t *testing.T) {

	cases := []struct {
		html string
		want *Article
	}{
		{`<h1 class="pill">Saturday, Sep 07, 2024</h1>`,
			&Article{Title: "Saturday, Sep 07, 2024", PubDate: createDate(2024, 9, 7)}},
		{`<h1 class="pill">Saturday, Sep 07, 2024 <span class="day-header-text">Some day</span></h1>`,
			&Article{Title: "Saturday, Sep 07, 2024 Some day", PubDate: createDate(2024, 9, 7)}},
		{`<h1 class="pill">Saturday, Sep 07, 2024</h1><div class="row no-gutter">Content <b>here</b></div>`,
			&Article{Title: "Saturday, Sep 07, 2024", PubDate: createDate(2024, 9, 7), ContentHTML: "Content <b>here</b>"}},
	}

	for _, tt := range cases {
		parser := BluesnewsParser{}
		article, err := parser.ParseHTML(strings.NewReader(tt.html))

		assertNoError(t, err)
		assertArticle(t, article, tt.want)
	}

	t.Run("Error parsing HTML", func(t *testing.T) {
		parser := BluesnewsParser{}
		htmlReader := MockErrReader{Error: errors.New("Mock error")}

		content, err := parser.ParseHTML(htmlReader)

		assertError(t, err, htmlReader.Error)
		if content != nil {
			t.Errorf("content should be nil")
		}
	})

	t.Run("Return error parsing date in title", func(t *testing.T) {
		html := `<h1>InvaliDate</h1>`

		parser := BluesnewsParser{}
		content, err := parser.ParseHTML(strings.NewReader(html))

		assertError(t, err, ErrTitleDateNotValid)
		if content != nil {
			t.Errorf("content should be nil")
		}
	})

	t.Run("GetHTMLArticleByDate returns correct HTML code by date", func(t *testing.T) {
		parser := BluesnewsParser{}
		saturdayHtml := `<h1 class="pill">Saturday, Sep 07, 2024</h1><div class="row no-gutter">Saturday content</div>`
		sundayHtml := `<h1 class="pill">Sunday, Sep 08, 2024 <span>Some day</span></h1><div class="row no-gutter">Sunday content</div>`
		html := fmt.Sprintf("<div>%s%s</div>", saturdayHtml, sundayHtml)

		contentSat, err := parser.GetHTMLArticleByDate(createDate(2024, 9, 7), html)

		assertNoError(t, err)
		if contentSat != saturdayHtml {
			t.Errorf("got %s\nwant %s", contentSat, saturdayHtml)
		}

		contentSun, err := parser.GetHTMLArticleByDate(createDate(2024, 9, 8), html)
		assertNoError(t, err)
		if contentSun != sundayHtml {
			t.Errorf("got %s\nwant %s", contentSun, sundayHtml)
		}
	})

	t.Run("GetHTMLArticleByDate returns error when no date is found", func(t *testing.T) {
		parser := BluesnewsParser{}
		saturdayHtml := `<h1 class="pill">Saturday, Sep 07, 2024</h1><div class="row no-gutter">Saturday content</div>`
		sundayHtml := `<h1 class="pill">Sunday, Sep 08, 2024 <span>Some day</span></h1><div class="row no-gutter">Sunday content</div>`
		html := fmt.Sprintf("<div>%s%s</div>", saturdayHtml, sundayHtml)

		_, err := parser.GetHTMLArticleByDate(createDate(2022, 1, 1), html)

		assertError(t, err, ErrArticleNotFound)
	})

}

func createDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func createSuccessFetchFn() BluesnewsFetcher {
	return BluesnewsFetcher{
		fetchFn: func(date time.Time) (*http.Response, error) {
			title := fmt.Sprintf("%d-%02d-%02d", date.Year(), date.Month(), date.Day())
			body := fmt.Sprintf("%s", title)

			respRecorder := httptest.ResponseRecorder{Code: 200, Body: bytes.NewBufferString(body)}
			resp := respRecorder.Result()

			return resp, nil
		},
	}
}

func createFetchFnWithStatusCode(statusCode int) BluesnewsFetcher {
	return BluesnewsFetcher{
		fetchFn: func(date time.Time) (*http.Response, error) {
			respRecorder := httptest.ResponseRecorder{Code: statusCode}
			resp := respRecorder.Result()

			return resp, nil
		},
	}
}

func createFetchFnWithError(err error) BluesnewsFetcher {
	return BluesnewsFetcher{
		fetchFn: func(date time.Time) (*http.Response, error) {
			return nil, err
		},
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error but didn't want one: %v", got)
	}
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func assertArticle(t testing.TB, got, want *Article) {
	t.Helper()

	if got.Title != want.Title {
		t.Errorf("Title: got %s, want %s", got.Title, want.Title)
	}

	if got.PubDate != want.PubDate {
		t.Errorf("PubDate: got %v, want %v", got.PubDate, want.PubDate)
	}

	if got.ContentHTML != want.ContentHTML {
		t.Errorf("ContentHTML: got %v, want %v", got.ContentHTML, want.ContentHTML)
	}

}
