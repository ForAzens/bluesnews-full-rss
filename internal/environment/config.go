package environment

type config struct {
	BaseUrl      string
	ArticlesPath string
}

func NewConfig() *config {
	return &config{
    BaseUrl:      "localhost:8080",
		ArticlesPath: "./articles/",
	}
}

func (c *config) SetBaseUrl(url string) *config {
	if len(url) > 0 {
		c.BaseUrl = url
	}

	return c
}

func (c *config) SetArticlesPath(path string) *config {
	if len(path) > 0 {
		c.ArticlesPath = path
	}

	return c
}
