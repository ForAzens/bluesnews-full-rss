package main

import (
	// "log"
	// "net/http"

	"log"
	"os"

	"github.com/ForAzens/bluesnews-full-rss/internal/bluenews"
	"github.com/ForAzens/bluesnews-full-rss/internal/feed"
)

var BASE_URL = "https://www.bluesnews.com"

func main() {

	articles := bluenews.FromDate()

	rss := feed.NewRss()
	for i := range articles {
		article := articles[i]
		rss.AddItem(feed.Item{
			Title:   article.Title,
			Content: article.ContentHTML,
		})
	}



	err := rss.EncodeToWriter(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
}
