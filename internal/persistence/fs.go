package persistence

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var DEFAULT_FILE_PATH = "./articles/"

func generateFileName(title string, date time.Time) string {
	stringDate := fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())

	return fmt.Sprintf("%s##%s.html", stringDate, strings.TrimSpace(title))
}

func WriteToFile(title string, date time.Time, content []byte) error {
	path := DEFAULT_FILE_PATH + generateFileName(title, date)

	return os.WriteFile(path, content, 0664)
}

func ListOfArticlesFromFs() ([]os.DirEntry, error) {
	return os.ReadDir(DEFAULT_FILE_PATH)
}

func ReadArticleFromFs(dirEntry os.DirEntry) ([]byte, error) {
	return os.ReadFile(DEFAULT_FILE_PATH + dirEntry.Name())
}

func TitleFromFileName(dirEntry os.DirEntry) string {
	substr := strings.Split(dirEntry.Name(), "##")

	return strings.ReplaceAll(substr[1], ".html", "")
}

func DateFromFileName(dirEntry os.DirEntry) (time.Time, error) {
	substr := strings.Split(dirEntry.Name(), "##")
	return time.Parse("20060102", substr[0])
}
