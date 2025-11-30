package rss

import (
	"encoding/xml"
	"html"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Description string    `xml:"description"`
		Link        string    `xml:"link"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Published   string `xml:"pubDate"`
}

func (r *RSSFeed) Parse(rssFeed []byte) error {
	err := xml.Unmarshal(rssFeed, r)
	if err != nil {
		return err
	}
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)
	for i := range r.Channel.Items {
		r.Channel.Items[i].Title = html.UnescapeString(r.Channel.Items[i].Title)
		r.Channel.Items[i].Description = html.UnescapeString(r.Channel.Items[i].Description)
	}
	return nil
}
