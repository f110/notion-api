package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.f110.dev/notion-api/v3"
)

func TestRender(t *testing.T) {
	buf, err := Render([]*notion.Block{
		{
			Type: notion.BlockTypeHeading2,
			Heading2: &notion.Heading{
				RichText: []*notion.RichTextObject{
					{
						Type: notion.RichTextObjectTypeText,
						Text: &notion.Text{
							Content: "Summary",
						},
						PlainText: "Summary",
					},
				},
			},
		},
		{
			Type: notion.BlockTypeParagraph,
			Paragraph: &notion.Paragraph{
				RichText: []*notion.RichTextObject{
					{
						Type: notion.RichTextObjectTypeText,
						Text: &notion.Text{
							Content: "Foobar",
						},
						PlainText: "Foobar",
					},
				},
			},
		},
		{
			Type: notion.BlockTypeCode,
			Code: &notion.Code{
				RichText: []*notion.RichTextObject{
					{
						Type: notion.RichTextObjectTypeText,
						Text: &notion.Text{
							Content: "$ uptime",
						},
						PlainText: "$ uptime",
					},
				},
				Language: "console",
			},
		},
		{
			Type: notion.BlockTypeBulletedListItem,
			BulletedListItem: &notion.Paragraph{
				RichText: []*notion.RichTextObject{
					{
						Type: notion.RichTextObjectTypeText,
						Text: &notion.Text{
							Content: "Baz",
						},
						PlainText: "Baz",
					},
				},
			},
		},
		{
			Type: notion.BlockTypeBulletedListItem,
			BulletedListItem: &notion.Paragraph{
				RichText: []*notion.RichTextObject{
					{
						Type: notion.RichTextObjectTypeText,
						Text: &notion.Text{
							Content: "Qux",
							Link: &notion.Link{
								URL: "https://example.com",
							},
						},
						PlainText: "Qux",
					},
				},
			},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "## Summary\n\nFoobar\n\n```console\n$ uptime\n```\n\n- Baz\n- [Qux](https://example.com)\n", string(buf))
}
