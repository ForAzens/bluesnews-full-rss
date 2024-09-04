package serve

import (
	"log"
	"net/http"

	"github.com/ForAzens/bluesnews-full-rss/internal/feed"
)

var BASE_URL = "https://www.bluesnews.com"

func CreateAndStartServer(address string, rss feed.Rss) {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /rss.xml", func(w http.ResponseWriter, r *http.Request) {
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
