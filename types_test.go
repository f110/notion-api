package notion

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	t.Run("MarshalAndUnmarshal", func(t *testing.T) {
		tm := Time{Time: time.Date(2020, 5, 3, 14, 15, 30, 0, time.Local)}

		b, err := json.Marshal(tm)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.HasPrefix(b, []byte("\"2020-05-03T14:15:30")) {
			t.Errorf("Marshal: Invalid format: %s", string(b))
		}

		u := &Time{}
		err = json.Unmarshal(b, u)
		if err != nil {
			t.Fatal(err)
		}
		if !tm.Time.Equal(u.Time) {
			t.Errorf("Unmarshal: Not equal value with before marshal")
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		u := &Time{}
		err := json.Unmarshal([]byte("{}"), u)
		if err != nil {
			t.Fatal(err)
		}
		if !u.Time.IsZero() {
			t.Errorf("\"{}\" has to be parsed Zero: %s", u.Time.Format(time.RFC3339))
		}
	})
}

func TestPropertyData_String(t *testing.T) {
	cases := []struct {
		In     *PropertyData
		Expect string
	}{
		{
			In: &PropertyData{
				Type: "title",
				Title: []*RichTextObject{
					{Type: "text", Text: &Text{Content: "Foo"}, PlainText: "foo"},
				},
			},
			Expect: "foo",
		},
		{
			In: &PropertyData{
				Type: "multi_select",
				MultiSelect: []*Option{
					{Name: "Foo"},
					{Name: "Bar"},
				},
			},
			Expect: "Foo, Bar",
		},
		{
			In: &PropertyData{
				Type: "text",
				Text: []*RichTextObject{
					{Type: "text", Text: &Text{Content: "Foo"}, PlainText: "Foo"},
					{Type: "text", Text: &Text{Content: "Bar"}, PlainText: "Bar"},
				},
			},
			Expect: "FooBar",
		},
		{
			In: &PropertyData{
				Type: "rich_text",
				RichText: []*RichTextObject{
					{Type: "text", Text: &Text{Content: "Foo"}, Annotations: &TextAnnotation{Bold: true}, PlainText: "Foo"},
					{Type: "text", Text: &Text{Content: "Bar"}, Annotations: &TextAnnotation{Italic: true}, PlainText: "Bar"},
				},
			},
			Expect: "FooBar",
		},
		{
			In: &PropertyData{
				Type:   "number",
				Number: 10,
			},
			Expect: "10",
		},
		{
			In: &PropertyData{
				Type:   "select",
				Select: &Option{Name: "foo"},
			},
			Expect: "foo",
		},
		{
			In: &PropertyData{
				Type: "date",
				Date: &DateProperty{
					Start: &Date{Time: MustParse(time.Parse("2006-01-2", "2021-12-1"))},
				},
			},
			Expect: "2021-12-01",
		},
		{
			In: &PropertyData{
				Type: "people",
				People: []*User{
					{
						Meta: &Meta{Object: "user"},
						Name: "Postmaster",
						Type: "person",
						Person: &Person{
							Email: "postmaster@example.com",
						},
					},
				},
			},
			Expect: "Postmaster <postmaster@example.com>",
		},
		{
			In: &PropertyData{
				Type: "files",
				Files: []*File{
					{Name: "foo.txt"},
				},
			},
			Expect: "foo.txt",
		},
		{
			In: &PropertyData{
				Type:     "checkbox",
				Checkbox: true,
			},
			Expect: "checked",
		},
		{
			In: &PropertyData{
				Type:     "checkbox",
				Checkbox: false,
			},
			Expect: "not checked",
		},
		{
			In: &PropertyData{
				Type: "url",
				URL:  "https://example.com/path/to",
			},
			Expect: "https://example.com/path/to",
		},
		{
			In: &PropertyData{
				Type:        "phone_number",
				PhoneNumber: "+81-23-4567-8901",
			},
			Expect: "+81-23-4567-8901",
		},
		{
			In: &PropertyData{
				Type:        "created_time",
				CreatedTime: &Time{Time: MustParse(time.Parse("2006-01-02T15:04:05", "2021-05-15T09:27:00"))},
			},
			Expect: "2021-05-15T09:27:00Z",
		},
		{
			In: &PropertyData{
				Type: "created_by",
				CreatedBy: &User{
					Meta: &Meta{Object: "user"},
					Name: "Postmaster",
					Type: "person",
					Person: &Person{
						Email: "postmaster@example.com",
					},
				},
			},
			Expect: "Postmaster <postmaster@example.com>",
		},
		{
			In: &PropertyData{
				Type:           "last_edited_time",
				LastEditedTime: &Time{Time: MustParse(time.Parse("2006-01-02T15:04:05", "2021-05-15T09:27:00"))},
			},
			Expect: "2021-05-15T09:27:00Z",
		},
		{
			In: &PropertyData{
				Type: "last_edited_by",
				LastEditedBy: &User{
					Meta: &Meta{Object: "user"},
					Name: "Postmaster",
					Type: "person",
					Person: &Person{
						Email: "postmaster@example.com",
					},
				},
			},
			Expect: "Postmaster <postmaster@example.com>",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Expect, func(t *testing.T) {
			assert.Equal(t, tc.Expect, tc.In.String())
		})
	}
}

func MustParse(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}
