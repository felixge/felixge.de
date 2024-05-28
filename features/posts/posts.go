package posts

import (
	"cmp"
	"context"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/sourcegraph/conc/pool"
)

type Post struct {
	Title   string
	PubDate time.Time
	URL     string
}

func Posts(ctx context.Context, log *slog.Logger) ([]*Post, error) {
	p := pool.NewWithResults[[]*Post]().WithErrors().WithContext(ctx)
	for _, feed := range feeds {
		p.Go(func(ctx context.Context) ([]*Post, error) {
			// TODO: add a timeout and fallback to a cached version (in the repo?)
			posts, err := fetchRSS(ctx, feed.URL)
			if err != nil {
				log.Error("failed to fetch feed", "feed", feed.URL, "err", err)
			}
			// don't return err, we want to continue fetching other feeds
			return posts, nil
		})
	}
	results, err := p.Wait()
	if err != nil {
		return nil, err
	}
	var posts []*Post
	for _, result := range results {
		posts = append(posts, result...)
	}
	slices.SortFunc(posts, func(a, b *Post) int {
		return cmp.Compare(b.PubDate.Unix(), a.PubDate.Unix())
	})
	return posts, nil
}

type feed struct {
	URL string
}

var feeds = []feed{
	{URL: "https://blog.felixge.de/rss/"},
	{URL: "https://felixge.de/posts.atom.xml"},
	{URL: "https://medium.com/feed/@felixge"},
}

func fetchRSS(ctx context.Context, url string) ([]*Post, error) {
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
