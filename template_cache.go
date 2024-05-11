package main

import "html/template"

type TemplateCache struct {
	cache map[string]*template.Template
}

func NewTemplateCache() *TemplateCache {
	tc := new(TemplateCache)
	tc.cache = make(map[string]*template.Template)
	return tc
}

func (self *TemplateCache) Get(name string) (*template.Template, error) {
	tmpl, ok := self.cache[name]
	if !ok {
		tmpl, err := template.ParseFiles("templates/"+name, "templates/base.tmpl")
		if err != nil {
			return nil, err
		}

		self.cache[name] = tmpl
		return tmpl, nil
	}

	return tmpl, nil
}

func (self *TemplateCache) Clear() {
	self.cache = make(map[string]*template.Template)
}
