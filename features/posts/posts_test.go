package posts

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseRSS(t *testing.T) {
	feeds := []struct {
		Fixture string
		Verify  func(t *testing.T, posts []*Post)
	}{
		{
			Fixture: "blog.felixge.de.rss.xml",
			Verify: func(t *testing.T, posts []*Post) {
				require.Len(t, posts, 11)
				blazing := posts[0]
				require.Equal(t, "Blazingly Fast Shadow Stacks for Go", blazing.Title)
				require.Equal(t, "2024-05-28 10:19:46 +0000 GMT", blazing.PubDate.String())
				require.Equal(t, "https://blog.felixge.de/blazingly-fast-shadow-stacks-for-go/", blazing.URL)
			},
		},
		{
			Fixture: "felixge.de.atom.xml",
			Verify: func(t *testing.T, posts []*Post) {
				require.Len(t, posts, 11)
				blazing := posts[0]
				require.Equal(t, "Blazingly Fast Shadow Stacks for Go", blazing.Title)
				require.Equal(t, "2024-05-28 10:19:46 +0000 GMT", blazing.PubDate.String())
				require.Equal(t, "https://blog.felixge.de/blazingly-fast-shadow-stacks-for-go/", blazing.URL)
			},
		},
	}

	for _, f := range feeds {
		t.Run(f.Fixture, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", f.Fixture))
			require.NoError(t, err)

			posts, err := parseRss(bytes.NewReader(data))
			require.NoError(t, err)
			f.Verify(t, posts)
		})
	}
}
