package markdown

import (
	"bytes"
	"fmt"
	"io"

	"go.f110.dev/notion-api/v3"
)

func Render(blocks []*notion.Block) ([]byte, error) {
	buf := new(bytes.Buffer)
	for i, b := range blocks {
		switch b.Type {
		case notion.BlockTypeHeading1:
			buf.WriteString("# ")
			for _, v := range b.Heading1.RichText {
				renderRichText(buf, v)
			}
			buf.WriteString("\n\n")
		case notion.BlockTypeHeading2:
			buf.WriteString("## ")
			for _, v := range b.Heading2.RichText {
				renderRichText(buf, v)
			}
			buf.WriteString("\n\n")
		case notion.BlockTypeHeading3:
			buf.WriteString("### ")
			for _, v := range b.Heading3.RichText {
				renderRichText(buf, v)
			}
			buf.WriteString("\n\n")
		case notion.BlockTypeParagraph:
			for _, v := range b.Paragraph.RichText {
				renderRichText(buf, v)
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
			buf.WriteString("```\n")
		case notion.BlockTypeBulletedListItem:
			if i > 0 && blocks[i-1].Type != notion.BlockTypeBulletedListItem {
				buf.WriteRune('\n')
			}

			buf.WriteString("- ")
			for _, v := range b.BulletedListItem.RichText {
				renderRichText(buf, v)
			}
			buf.WriteRune('\n')
		case notion.BlockTypeNumberedListItem:
			if i > 0 && blocks[i-1].Type != notion.BlockTypeNumberedListItem {
				buf.WriteRune('\n')
			}

			buf.WriteString("1. ")
			for _, v := range b.NumberedListItem.RichText {
				renderRichText(buf, v)
			}
			buf.WriteRune('\n')
		case notion.BlockTypeDivider:
			buf.WriteString("\n---\n")
		default:
			return nil, fmt.Errorf("not supported block type: %s", b.Type)
		}
	}

	return buf.Bytes(), nil
}

func renderRichText(w io.Writer, b *notion.RichTextObject) {
	if b.Annotations != nil {
		if b.Annotations.Code {
			w.Write([]byte{'`'})
		}
	}

	if b.Text.Link != nil {
		io.WriteString(w, fmt.Sprintf("[%s](%s)", b.Text.Content, b.Text.Link.URL))
	} else {
		io.WriteString(w, b.Text.Content)
	}

	if b.Annotations != nil {
		if b.Annotations.Code {
			w.Write([]byte{'`'})
		}
	}
}
