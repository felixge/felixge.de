package posts

import (
	"context"
	"io"
	"net/http"

	"github.com/felixge/felixge.de/commons/appcontext"
	"github.com/mmcdole/gofeed"
)

/*

Feeds:

1. At build time, fetch feeds and store them in the file system.
2. At run time, read feeds every 1min and update the file system.
3. At request time, read feeds from the file system.
*/

func WriteFeeds(ctx context.Context, a *appcontext.AppContext) error {
	return nil
	// p := pool.NewWithResults[[]*Post]().WithErrors().WithContext(ctx)
	// for _, feed := range feeds {
	// 	p.Go(func(ctx context.Context) ([]*Post, error) {
	// 		// TODO: add a timeout and fallback to a cached version (in the repo?)
	// 		posts, err := fetchRSS(ctx, feed.URL)
	// 		if err != nil {
	// 			a.Log.Error("failed to fetch feed", "feed", feed.URL, "err", err)
	// 		}
	// 		// don't return err, we want to continue fetching other feeds
	// 		return posts, nil
	// 	})
	// }
	// results, err := p.Wait()
	// if err != nil {
	// 	return nil, err
	// }
	// var posts []*Post
	// for _, result := range results {
	// 	posts = append(posts, result...)
	// }
	// slices.SortFunc(posts, func(a, b *Post) int {
	// 	return cmp.Compare(b.PubDate.Unix(), a.PubDate.Unix())
	// })
	// return posts, nil
}

type feed struct {
	Name string
	URL  string
}

var feeds = []feed{
	{
		Name: "ghost",
		URL:  "https://blog.felixge.de/rss/",
	},
	{
		Name: "hugo",
		URL:  "https://felixge.de/posts.atom.xml",
	},
	{
		Name: "medium",
		URL:  "https://medium.com/feed/@felixge",
	},
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
