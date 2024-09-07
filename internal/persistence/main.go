package persistence

import (
	"github.com/ForAzens/bluesnews-full-rss/internal/bluenews"
)

type ArticleFetcher interface {
	FetchAll() []bluenews.Article
}

type ArticleSaver interface {
	Save(a bluenews.Article) error
}

type ArticleManager interface {
	ArticleFetcher
	ArticleSaver
}
