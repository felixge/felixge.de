package generator

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/felixge/makefs"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var baseUrl = os.Getenv("BASE_URL")
var author = Person{Name: "Felix Geisendörfer", URI: baseUrl + "/"}

func makePosts(fs *makefs.Fs) error {
	if baseUrl == "" {
		return fmt.Errorf("no BASE_URL env var")
	}

	fs.Make("/public/posts/atom.xml", []string{"/posts/*/*/*/*.md"}, func(t *makefs.Task) error {
		atomFeed := Feed{
			XMLName: xml.Name{"http://www.w3.org/2005/Atom", "feed"},
			Title:   "Felix Geisendörfer",
			Link:    []Link{{Href: baseUrl}},
			Id:      "urn:uuid:d306c549-0ae3-47d3-8814-3286ee297933",
			Updated: time.Now(),
			Author:  author,
		}

		for _, source := range t.Sources() {
			post, err := parsePost(source)
			if err != nil {
				return err
			}

			html, err := post.Html()
			if err != nil {
				return err
			}

			entry := Entry{
				Author:    author,
				Title:     post.Title,
				Id:        baseUrl + post.Url,
				Updated:   post.Updated,
				Published: post.Published,
				Link:      []Link{{Href: baseUrl + post.Url}},
				Content:   Text{Type: "html", Body: html},
			}

			atomFeed.Entry = append(atomFeed.Entry, entry)
		}

		data, err := xml.MarshalIndent(atomFeed, "", "  ")
		if err != nil {
			return err
		}

		if _, err := io.WriteString(t.Target(), xml.Header); err != nil {
			return err
		}

		if _, err := t.Target().Write(data); err != nil {
			return err
		}

		return nil
	})

	fs.Make("/public/%.html", []string{"/posts/%.md", "/templates/post.html", "/templates/layout.html"}, func(t *makefs.Task) error {
		post, err := parsePost(t.Source())
		if err != nil {
			return err
		}

		html, err := post.Html()
		if err != nil {
			return err
		}

		return render(
			t.Target(),
			t.Sources()[1],
			t.Sources()[2],
			map[string]interface{}{
				"Title": post.Title,
				"Post":  post,
				"Html":  template.HTML(html),
				"BaseUrl": baseUrl,
			},
		)
	})

	return nil
}

var metaRegexp = regexp.MustCompile("(?s:{.+}\n)")

func parsePost(s *makefs.Source) (*post, error) {
	data, err := ioutil.ReadAll(s)
	if err != nil {
		return nil, err
	}

	meta := metaRegexp.Find(data)
	if len(meta) == 0 {
		return nil, fmt.Errorf("generator: could not meta info for: %s", s.Path())
	}

	p := &post{}
	if err := json.Unmarshal(meta, p); err != nil {
		fmt.Printf("meta: %s\n", err)
		return nil, fmt.Errorf("generator: %s: %s", s.Path(), err)
	}

	p.Url = s.Path()[len("/posts"):]
	p.Url = p.Url[0:len(p.Url)-len(".md")] + ".html"
	p.Markdown = string(data[len(meta):])

	return p, nil
}

type post struct {
	Title     string
	Url       string
	Updated   time.Time
	Published time.Time
	Markdown  string
}

func (p *post) Html() (string, error) {
	buf := bytes.NewBuffer(nil)

	cmd := exec.Command(__dirname + "/processors/bin/markdown.js")
	cmd.Stdin = strings.NewReader(p.Markdown)
	cmd.Stdout = buf
	cmd.Stderr = buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("generator: %s: %s", buf, err)
	}

	return buf.String(), nil
}

type Feed struct {
	XMLName xml.Name  `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string    `xml:"title"`
	Id      string    `xml:"id"`
	Link    []Link    `xml:"link"`
	Updated time.Time `xml:"updated"`
	Author  Person    `xml:"author"`
	Entry   []Entry   `xml:"entry"`
}

type Entry struct {
	Title     string    `xml:"title"`
	Id        string    `xml:"id"`
	Link      []Link    `xml:"link"`
	Updated   time.Time `xml:"updated"`
	Published time.Time `xml:"published"`
	Author    Person    `xml:"author"`
	Content   Text   `xml:"content"`
}

type Link struct {
	Rel  string `xml:"rel,attr,omitempty"`
	Href string `xml:"href,attr"`
}

type Person struct {
	Name     string `xml:"name"`
	URI      string `xml:"uri"`
	InnerXML string `xml:",innerxml"`
}

type Text struct {
	Type string `xml:"type,attr,omitempty"`
	Body string `xml:",chardata"`
}
