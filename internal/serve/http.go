package serve

import (
	"log"
	"net/http"

	"github.com/ForAzens/bluesnews-full-rss/internal/bluenews"
	"github.com/ForAzens/bluesnews-full-rss/internal/feed"
)

var BASE_URL = "https://www.bluesnews.com"

func CreateAndStartServer(address string, fetcher bluenews.ArticleFetcher) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /rss.xml", func(w http.ResponseWriter, r *http.Request) {
		articles := fetcher()
		rss := feed.NewRss()

		for i := range articles {
			article := articles[i]
			rss.AddItem(feed.Item{
				Title:   article.Title,
				Content: article.ContentHTML,
			})
		}

		if err := rss.EncodeToWriter(w); err != nil {
			log.Fatalln(err)
		}
	})

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	log.Println("Start server in " + server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
