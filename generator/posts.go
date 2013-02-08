package generator

import (
	"encoding/xml"
	"fmt"
	"github.com/felixge/makefs"
	"io"
	"io/ioutil"
	"os"
	gopath "path"
	"regexp"
	"strings"
	"time"
)

const postsPath = "/posts/"

func makePosts(fs *makefs.Fs) error {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		return fmt.Errorf("no BASE_URL env var")
	}

	postsDir, err := fs.Open(postsPath)
	if err != nil {
		return err
	}

	postStats, err := postsDir.Readdir(0)
	if err != nil {
		return err
	}

	posts := make([]post, 0, len(postStats))
	postSources := make([]string, 0)
	for _, postStat := range postStats {
		name := postStat.Name()
		if name[0:1] == "." {
			continue
		}

		ext := gopath.Ext(name)
		if ext != ".html" {
			continue
		}

		dateFormat := "2006-01-02"
		date, err := time.Parse(dateFormat, name[0:len(dateFormat)])
		if err != nil {
			return err
		}

		url := fmt.Sprintf(
			"/%s/%s",
			strings.ToLower(date.Format("2006/Jan-2")),
			name[len(dateFormat)+1:],
		)

		p := post{
			Url:        url,
			SourcePath: postsPath + name,
			TargetPath: "/public" + url,
			Date:       date,
		}

		fs.Make(
			p.TargetPath,
			[]string{p.SourcePath, "/layouts/default.html"},
			func(t *makefs.Task) error {
				sources := t.Sources()
				return render(
					t.Target(),
					sources[0],
					sources[1],
					nil,
				)
			},
		)

		posts = append(posts, p)
		postSources = append(postSources, p.SourcePath)
	}

	fs.Make("/public/posts.atom.xml", postSources, func(t *makefs.Task) error {
		author := Person{Name: "Felix Geisendörfer"}

		atomFeed := Feed{
			XMLName: xml.Name{"http://www.w3.org/2005/Atom", "feed"},
			Title:   "TODO - Proper title",
			Link: []Link{{Href: baseUrl}},
			Id:      baseUrl,
			Updated: time.Now(),
			Author:  author,
		}

		for i, p := range posts {
			source := t.Sources()[i]
			if err := parsePost(source, &p); err != nil {
				return err
			}

			entry := Entry{
				Author:  author,
				Title:   p.Title,
				Id:      baseUrl+p.Url,
				Updated: time.Now(),
				Link: []Link{
					{Href: baseUrl+p.Url},
				},
			}

			atomFeed.Entry = append(atomFeed.Entry, entry)
		}

		data, err := xml.MarshalIndent(atomFeed, "", "  ")
		if err != nil {
			return err
		}

		if _, err := t.Target().Write(data); err != nil {
			return err
		}

		return nil
	})

	return nil
}

var titleRegexp = regexp.MustCompile("<h1>(.+)</h1>")

func parsePost(r io.Reader, p *post) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	match := titleRegexp.FindStringSubmatch(string(data))
	if len(match) != 2 {
		return fmt.Errorf("could not find title for post: %s", data)
	}

	p.Title = match[1]

	return nil
}

type post struct {
	Title      string
	SourcePath string
	TargetPath string
	Url        string
	Date       time.Time
}

type Feed struct {
	XMLName xml.Name  `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string    `xml:"title"`
	Id      string    `xml:"id"`
	Link    []Link    `xml:"link"`
	Updated time.Time `xml:"updated,attr"`
	Author  Person    `xml:"author"`
	Entry   []Entry   `xml:"entry"`
}

type Entry struct {
	Title   string    `xml:"title"`
	Id      string    `xml:"id"`
	Link    []Link    `xml:"link"`
	Updated time.Time `xml:"updated"`
	Author  Person    `xml:"author"`
	Summary Text      `xml:"summary"`
}

type Link struct {
	Rel  string `xml:"rel,attr,omitempty"`
	Href string `xml:"href,attr"`
}

type Person struct {
	Name     string `xml:"name"`
	URI      string `xml:"uri"`
	Email    string `xml:"email"`
	InnerXML string `xml:",innerxml"`
}

type Text struct {
	Type string `xml:"type,attr,omitempty"`
	Body string `xml:",chardata"`
}
