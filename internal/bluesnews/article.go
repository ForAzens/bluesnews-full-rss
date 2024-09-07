package bluesnews

import (
	"time"
)

type ArticleFetcher func() []Article

type Article struct {
	Title       string
	PubDate     time.Time
	ContentHTML string
}
