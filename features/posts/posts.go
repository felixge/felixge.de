package posts

import (
	"cmp"
	"context"
	"log/slog"
	"slices"
	"time"

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
