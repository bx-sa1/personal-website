package main

import (
	"bufio"
	"bytes"
	"errors"
	"html/template"
	"log"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gopkg.in/yaml.v3"
)

type Markdown struct {
	FrontMatter map[string]interface{}
	Markdown    ast.Node
}

func (md *Markdown) RenderMarkdown() template.HTML {
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	return template.HTML(markdown.Render(md.Markdown, renderer))
}

func parseFrontMatter(md []byte) (string, string, error) {
	reader := bytes.NewReader(md)
	scanner := bufio.NewScanner(reader)
	var _yaml string
	var raw string

	if scanner.Scan() && scanner.Text() == "---" {
		raw += scanner.Text() + "\n"
		for scanner.Scan() && scanner.Text() != "---" {
			line := scanner.Text() + "\n"
			_yaml += line
			raw += line
		}
		raw += scanner.Text() + "\n"
		return _yaml, raw, nil
	} else {
		return "", "", errors.New("no FrontMatter")
	}
}

func ParseMarkdown(file string) (*Markdown, error) {
	md, err := LoadFile(file)
	if err != nil {
		return nil, err
	}

	_yaml, raw, err := parseFrontMatter(md)
	if err != nil {
		return nil, err
	}

	log.Println(raw)

	front_matter := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(_yaml), front_matter)
	if err != nil {
		return nil, err
	}

	ext := parser.CommonExtensions
	p := parser.NewWithExtensions(ext)
	doc := p.Parse(md[len(raw):])

	return &Markdown{front_matter, doc}, nil
}
