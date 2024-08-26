package bluenews

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var BASE_URL = "https://www.bluesnews.com"

func removeTreeClasses(s *goquery.Selection) {
	childrenSel := s.Children()

	for i := range childrenSel.Nodes {
		sel := childrenSel.Eq(i)
		sel.RemoveClass()

		removeTreeClasses(sel)
	}
}

func fixLinks(s *goquery.Selection) {
	s.Find("h2 > a").Each(func(i int, sel *goquery.Selection) {
		href, ok := sel.Attr("href")
		if ok {
			sel.SetAttr("href", BASE_URL+href)
		}
	})
}

func parseBody(rc io.ReadCloser) []Article {
	articles := make([]Article, 2)

	doc, err := goquery.NewDocumentFromReader(rc)
	if err != nil {
		log.Fatalln(err)
	}
	articleSelection := doc.Find("h1.pill + div.row.no-gutter")
	removeTreeClasses(articleSelection)
	fixLinks(articleSelection)
	// 	titleSelection := s.Prev()
	// 	title := titleSelection.Text()

	for i := range articleSelection.Nodes {
		s := articleSelection.Eq(i)

		removeTreeClasses(s)
		fixLinks(s)

		titleSelection := s.Prev()
		title := titleSelection.Text()

		content, err := s.Html()
		if err != nil {
			log.Println(err)
			continue
		}

		articles = append(articles, Article{
			Title:       title,
			PubDate:     time.Now(),
			ContentHTML: content,
		})
	}

	return articles
}

func FromDate() []Article {
	resp, err := http.Get("https://www.bluesnews.com/cgi-bin/blammo.pl?mode=archive&display=20240825")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	return parseBody(resp.Body)
}
