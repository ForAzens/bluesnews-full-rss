package bluenews

import "time"

type Article struct {
	Title       string
	PubDate     time.Time
	ContentHTML string
}


