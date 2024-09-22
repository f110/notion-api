package notiontest

import "go.f110.dev/notion-api/v3"

func NewDatabase(title string) *notion.Database {
	return &notion.Database{
		Meta: &notion.Meta{},
		Title: []*notion.RichTextObject{
			{Type: notion.RichTextObjectTypeText, Text: &notion.Text{Content: title}},
		},
		Properties: make(map[string]*notion.PropertyMetadata),
	}
}

func NewPage(name string, opts ...PageOpt) *notion.Page {
	p := &notion.Page{
		Meta: &notion.Meta{},
		Properties: map[string]*notion.PropertyData{
			"Name": {
				ID:    "title",
				Type:  "title",
				Title: []*notion.RichTextObject{{Type: notion.RichTextObjectTypeText, Text: &notion.Text{Content: name}}},
			},
		},
	}
	for _, v := range opts {
		v(p)
	}
	return p
}

type PageOpt func(*notion.Page)

func PageProperty(name string, data *notion.PropertyData) PageOpt {
	return func(page *notion.Page) {
		page.Properties[name] = data
	}
}
