package markdown

import (
	"bytes"
	"fmt"

	"go.f110.dev/notion-api/v3"
)

func Render(blocks []*notion.Block) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, b := range blocks {
		switch b.Type {
		case notion.BlockTypeHeading1:
			buf.WriteString("# ")
			buf.WriteString(b.Heading1.RichText[0].Text.Content)
			buf.WriteString("\n\n")
		case notion.BlockTypeHeading2:
			buf.WriteString("## ")
			buf.WriteString(b.Heading2.RichText[0].Text.Content)
			buf.WriteString("\n\n")
		case notion.BlockTypeHeading3:
			buf.WriteString("### ")
			buf.WriteString(b.Heading3.RichText[0].Text.Content)
			buf.WriteString("\n\n")
		case notion.BlockTypeParagraph:
			for _, v := range b.Paragraph.RichText {
				if v.Annotations != nil {
					if v.Annotations.Code {
						buf.WriteRune('`')
					}
				}

				if v.Text.Link != nil {
					buf.WriteString(fmt.Sprintf("[%s](%s)", v.Text.Content, v.Text.Link.URL))
				} else {
					buf.WriteString(v.Text.Content)
				}

				if v.Annotations != nil {
					if v.Annotations.Code {
						buf.WriteRune('`')
					}
				}
			}
			buf.WriteString("\n\n")
		case notion.BlockTypeCode:
			buf.WriteString("```")
			buf.WriteString(b.Code.Language)
			buf.WriteRune('\n')
			for _, v := range b.Code.RichText {
				buf.WriteString(v.Text.Content)
			}
			buf.WriteRune('\n')
			buf.WriteString("```\n\n")
		case notion.BlockTypeBulletedListItem:
			buf.WriteString("- ")
			buf.WriteString(b.BulletedListItem.RichText[0].Text.Content)
			buf.WriteRune('\n')
		default:
			return nil, fmt.Errorf("not supported block type: %s", b.Type)
		}
	}

	return buf.Bytes(), nil
}
