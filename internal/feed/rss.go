package feed

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel
}

type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Items       []Item
}

type Item struct {
	XMLName        xml.Name `xml:"item"`
	Title          string   `xml:"title"`
	Content        string   `xml:"-"`
	ContentEncoded ContentEncoded
}

// Struct needed because `xml:"content:encoded,cdata"` doesn't work
type ContentEncoded struct {
	XMLName xml.Name `xml:"content:encoded"`
	Content string   `xml:",cdata"`
}

func NewRssTest() Rss {
	return Rss{
		Version: "2.0",
		Channel: Channel{
			Title: "Pepe",
			Items: []Item{
				{Title: "Hey"},
			},
		},
	}
}

func NewRss() Rss {
	return Rss{
		Version: "2.0",
		Channel: Channel{
			Title:       "",
			Link:        "",
			Description: "",
			Items:       []Item{},
		},
	}
}

func (r *Rss) UpdateChannel(title string, link string, description string) {
	r.Channel.Title = title
	r.Channel.Link = link
	r.Channel.Description = description
}

func (r *Rss) AddItem(item Item) {
	item.ContentEncoded = ContentEncoded{
		Content: item.Content,
	}
	log.Printf("Item added")
	r.Channel.Items = append(r.Channel.Items, item)
}

func (r *Rss) EncodeToWriter(w io.Writer) error {
	enc := xml.NewEncoder(w)
	defer enc.Close()
	enc.Indent("  ", "    ")

	return enc.Encode(r)
}

func (r *Rss) EncodeToStr() (string, error) {
	buf := new(bytes.Buffer)
	err := r.EncodeToWriter(buf)

	return buf.String(), err
}
