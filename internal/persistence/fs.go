package persistence

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ForAzens/bluesnews-full-rss/internal/bluenews"
)

type FileSystemProvider struct {
	ArticlesFolderPath string
}

func (fsp *FileSystemProvider) FetchAll() []bluenews.Article {
	filesPath, _ := os.ReadDir(fsp.ArticlesFolderPath)
	var articles []bluenews.Article

	for i := range filesPath {
		path := filesPath[i]
		file, err := os.ReadFile(fsp.ArticlesFolderPath + path.Name())

		if err != nil {
			log.Printf("ERROR: File %s couldn't be read.", path.Name())
			continue
		}

		date, err := dateFromFileName(path.Name())

		if err != nil {
			log.Printf("ERROR: File %s couldn't be parsed.", path.Name())
			continue
		}

		articles = append(articles, bluenews.Article{
			Title:       titleFromFileName(path.Name()),
			PubDate:     date,
			ContentHTML: string(file),
		})
	}

	return articles
}

func (fsp *FileSystemProvider) Save(a bluenews.Article) error {
	path := fsp.ArticlesFolderPath + generateFileName(a.Title, a.PubDate)

	return os.WriteFile(path, []byte(a.ContentHTML), 0664)
}

func generateFileName(title string, date time.Time) string {
	stringDate := fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())

	return fmt.Sprintf("%s##%s.html", stringDate, strings.TrimSpace(title))
}

func titleFromFileName(filename string) string {
	substr := strings.Split(filename, "##")

	return strings.ReplaceAll(substr[1], ".html", "")
}

func dateFromFileName(filename string) (time.Time, error) {
	substr := strings.Split(filename, "##")
	return time.Parse("20060102", substr[0])
}
