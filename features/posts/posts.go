package posts

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/antchfx/xmlquery"
)

type Post struct {
	Title   string
	PubDate time.Time
	URL     string
}

func Posts(ctx context.Context) ([]*Post, error) {
	return RSSPosts(ctx, "https://blog.felixge.de/rss/")
}

type feed struct {
	URL string
}

var feeds = []feed{
	{"https://blog.felixge.de/rss/"},
}

func RSSPosts(ctx context.Context, url string) ([]*Post, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return parseRss(res.Body)
}

func parseRss(r io.Reader) ([]*Post, error) {
	const timeLayout = "Mon, 02 Jan 2006 15:04:05 MST"
	doc, err := xmlquery.Parse(r)
	if err != nil {
		return nil, err
	}
	var posts []*Post
	for _, item := range xmlquery.Find(doc, "/rss/channel/item") {
		post := &Post{
			Title: item.SelectElement("title").InnerText(),
			URL:   item.SelectElement("link").InnerText(),
		}
		pubDateText := item.SelectElement("pubDate").InnerText()
		post.PubDate, err = time.Parse(timeLayout, pubDateText)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}
