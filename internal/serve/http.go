package serve

import (
	"log"
	"net/http"

	"github.com/ForAzens/bluesnews-full-rss/internal/feed"
	"github.com/ForAzens/bluesnews-full-rss/internal/persistence"
)

var BASE_URL = "https://www.bluesnews.com"

func CreateAndStartServer(address string, am persistence.ArticleManager) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /rss.xml", func(w http.ResponseWriter, r *http.Request) {
		articles := am.FetchAll()
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
		Addr:    address,
		Handler: mux,
	}

	log.Println("Start server in " + server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
