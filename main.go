package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ForAzens/bluesnews-full-rss/internal/bluesnews"
	"github.com/ForAzens/bluesnews-full-rss/internal/environment"
	"github.com/ForAzens/bluesnews-full-rss/internal/persistence"
	"github.com/ForAzens/bluesnews-full-rss/internal/serve"
)

var BASE_URL = os.Getenv("BLUENEWS_RSS_BASE_URL")
var ARTICLES_PATH = os.Getenv("BLUENEWS_RSS_ARTICLES_PATH")

func main() {
	var mode string
	var lastDays int

	flag.StringVar(&mode, "mode", "serve", "Different modes to use: 'serve' or 'fetch'")
	flag.IntVar(&lastDays, "lastDays", 7, "To retrieve the articles of the last X days. Default: 7")

	flag.Parse()

	config := environment.NewConfig()
	config.SetBaseUrl(BASE_URL)
	config.SetArticlesPath(ARTICLES_PATH)

	persistenceProvider := persistence.FileSystemProvider{
		ArticlesFolderPath: config.ArticlesPath,
	}

	switch mode {
	case "serve":
		serve.CreateAndStartServer(config.BaseUrl, &persistenceProvider)
	case "fetch":
		for i := lastDays; i >= 0; i-- {
			date := time.Now().AddDate(0, 0, -i)
			log.Printf("Fetching articles for date: %v", date)
			client := bluesnews.NewBluesnewsClient()

			article, err := client.GetArticleFromDate(date)

			if err != nil {
				log.Printf("Failed to get article from date %v: %v", date, err)
				continue
			}

			err = persistenceProvider.Save(*article)
			if err != nil {
				log.Printf("Failed to save article for %v", article.PubDate)
				log.Println(err)
			}
		}

		log.Println("Fetching articles...")
	default:
		log.Fatalln("Unknown mode parameter")

	}

}
