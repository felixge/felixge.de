package posts

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
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
	fp := gofeed.NewParser()
	feed, err := fp.Parse(r)
	if err != nil {
		return nil, err
	}

	var posts []*Post
	for _, item := range feed.Items {
		post := &Post{
			Title:   item.Title,
			URL:     item.Link,
			PubDate: *item.PublishedParsed,
		}
		posts = append(posts, post)
	}
	return posts, nil
}
