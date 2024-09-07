package persistence

import (
	"github.com/ForAzens/bluesnews-full-rss/internal/bluesnews"
)

type ArticleFetcher interface {
	FetchAll() []bluesnews.Article
}

type ArticleSaver interface {
	Save(a bluesnews.Article) error
}

type ArticleManager interface {
	ArticleFetcher
	ArticleSaver
}
