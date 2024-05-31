package markdown

import (
	"bytes"
	"embed"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func NewMarkdown(fs embed.FS) *Markdown {
	return &Markdown{fs: fs}
}

type Markdown struct {
	fs embed.FS
}

func (r *Markdown) Render(path string) (*RenderResult, error) {
	data, err := r.fs.ReadFile(path)
	if err != nil {
		return nil, err
	}
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(),
	)
	p := md.Parser()
	tp := &titleParser{Parser: p}
	md.SetParser(tp)
	var buf bytes.Buffer
	if err := md.Convert(data, &buf); err != nil {
		return nil, err
	}
	return &RenderResult{Title: tp.Title, HTML: buf.Bytes()}, nil
}

type RenderResult struct {
	// Title is the first h1 heading in the markdown file.
	Title string
	// HTML is the rendered markdown content.
	HTML []byte
}

type titleParser struct {
	Title string
	parser.Parser
}

func (p *titleParser) Parse(reader text.Reader, opts ...parser.ParseOption) ast.Node {
	root := p.Parser.Parse(reader, opts...)
	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindHeading && n.(*ast.Heading).Level == 1 {
			p.Title = string(n.Text(reader.Source()))
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})
	return root
}
