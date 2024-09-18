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
