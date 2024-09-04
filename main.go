package main

import (
	"flag"
	"log"

	"github.com/ForAzens/bluesnews-full-rss/internal/bluenews"
	"github.com/ForAzens/bluesnews-full-rss/internal/feed"
	"github.com/ForAzens/bluesnews-full-rss/internal/serve"
)

var BASE_URL = "https://www.bluesnews.com"

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "serve", "Different modes to use: 'serve' or 'fetch'")
	flag.Parse()

	switch mode {
	case "serve":
		articles := bluenews.FromDate()
		log.Println("number of articles")
		log.Println(len(articles))

		rss := feed.NewRss()
		for i := range articles {
			article := articles[i]
			rss.AddItem(feed.Item{
				Title:   article.Title,
				Content: article.ContentHTML,
			})
		}
		serve.CreateAndStartServer("localhost:8080", rss)
	case "fetch":
    // TODO: Fetch articles and save it in local filesystem or database
		log.Println("Fetching articles")
	default:
		log.Fatalln("Unknown mode parameter")

	}

}
