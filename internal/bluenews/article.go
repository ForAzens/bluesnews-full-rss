package bluenews

import (
	"log"
	"time"

	"github.com/ForAzens/bluesnews-full-rss/internal/persistence"
)

type ArticleFetcher func() []Article

type Article struct {
	Title       string
	PubDate     time.Time
	ContentHTML string
}

func GetArticlesFromFS() []Article {
	var articles []Article

	entries, err := persistence.ListOfArticlesFromFs()
	if err != nil {
		log.Fatalf("Unable to read articles folder")
	}

	for i := range entries {
		path := entries[i]
		content, err := persistence.ReadArticleFromFs(path)
		if err != nil {
			log.Fatalf("Unable to read file: %s", path.Name())
		}

		title := persistence.TitleFromFileName(path)
		date, err := persistence.DateFromFileName(path)
		if err != nil {
			log.Fatalf("Unable to parse date for: %s", path.Name())
		}

		articles = append(articles, Article{
			Title:       title,
			PubDate:     date,
			ContentHTML: string(content),
		})
	}
	return articles
}
