package notion

import (
	"bytes"
	"time"
)

type Meta struct {
	// Type of object
	Object string `json:"object,omitempty"`
	// Unique identifier
	ID string `json:"id"`
}

type ListMeta struct {
	// Type of object
	Object     string `json:"object,omitempty"`
	HasMore    bool   `json:"has_more,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type Time struct {
	time.Time
}

const ISO8601 = "2006-01-02T15:04:05.999-0700"

func (t *Time) UnmarshalJSON(data []byte) error {
	if bytes.HasSuffix(data, []byte("Z\"")) {
		p, err := time.Parse("\"2006-01-02T15:04:05.999Z\"", string(data))
		if err != nil {
			return err
		}
		t.Time = p
		return nil
	}

	p, err := time.Parse("\""+ISO8601+"\"", string(data))
	if err != nil {
		return err
	}
	t.Time = p
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.Time.Format(ISO8601) + "\""), nil
}

type User struct {
	*Meta

	// Type of the user.
	Type string `json:"type"`
	// Displayed name
	Name string `json:"name"`
	// Avatar image url
	AvatarURL string `json:"avatar_url"`

	Person *Person `json:"person"`
	Bot    *Bot    `json:"bot"`
}

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("\"2006-01-02\"", string(data))
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.Time.Format("2006-01-02") + "\""), nil
}

type Person struct {
	// Email address of people
	Email string `json:"email"`
}

type Bot struct{}

type UserList struct {
	*ListMeta
	Results []*User `json:"results"`
}

type RichTextObject struct {
	PlainText   string          `json:"plain_text"`
	Href        string          `json:"href"`
	Annotations *TextAnnotation `json:"annotations"`
	Type        string          `json:"type"`

	Text     *Text     `json:"text"`
	Mention  *Mention  `json:"mention"`
	Equation *Equation `json:"equation"`
}

type TextAnnotation struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

type Text struct {
	Content string `json:"content"`
	Link    *Link  `json:"link"`
}

type Mention struct {
	Type string `json:"type"`

	User     *User `json:"user"`
	Page     *Meta `json:"page"`
	Database *Meta `json:"database"`
}

type Equation struct {
	Expression string `json:"expression"`
}

type Link struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Database struct {
	*Meta

	CreatedTime    Time                 `json:"created_time"`
	LastEditedTime Time                 `json:"last_edited_time"`
	Title          []*RichTextObject    `json:"title"`
	Properties     map[string]*Property `json:"properties"`
}

type DatabaseList struct {
	*ListMeta
	Results []*Database `json:"results"`
}

type Property struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`

	Title          []*RichTextObject      `json:"title,omitempty"`
	RichText       []*RichTextObject      `json:"rich_text,omitempty"`
	Number         *NumberProperty        `json:"number,omitempty"`
	Select         *SelectProperty        `json:"select,omitempty"`
	MultiSelect    []*MultiSelectProperty `json:"multi_select,omitempty"`
	Date           *DateProperty          `json:"date,omitempty"`
	Formula        *FormulaProperty       `json:"formula,omitempty"`
	Relation       *RelationProperty      `json:"relation,omitempty"`
	Checkbox       bool                   `json:"checkbox,omitempty"`
	Rollup         *RollupProperty        `json:"rollup,omitempty"`
	People         *PeopleProperty        `json:"people,omitempty"`
	Files          *FilesProperty         `json:"files,omitempty"`
	URL            string                 `json:"url,omitempty"`
	Email          string                 `json:"email,omitempty"`
	PhoneNumber    string                 `json:"phone_number,omitempty"`
	CreatedTime    *Time                  `json:"created_time,omitempty"`
	LastEditedTime *Time                  `json:"last_edited_time,omitempty"`
	LastEditedBy   *Meta                  `json:"last_edited_by,omitempty"`
}

type NumberProperty struct {
	Number int `json:"number"`
}

type SelectProperty struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type MultiSelectProperty struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type DateProperty struct {
	Start *Date `json:"start,omitempty"`
	End   *Date `json:"end,omitempty"`
}

type FormulaProperty struct {
	Type string `json:"type"`

	String  string        `json:"string,omitempty"`
	Number  int           `json:"number,omitempty"`
	Boolean bool          `json:"boolean,omitempty"`
	Date    *DateProperty `json:"date,omitempty"`
}

type RelationProperty struct {
	Relation []*Meta `json:"relation"`
}

type RollupProperty struct {
	Type string `json:"type"`

	Number *NumberProperty `json:"number,omitempty"`
	Date   *DateProperty   `json:"date,omitempty"`
	Array  []*Property     `json:"array,omitempty"`
}

type PeopleProperty struct {
	People *User `json:"people"`
}

type FilesProperty struct {
	Files []*Meta `json:"files"`
}

// Page is a page object.
// ref: https://developers.notion.com/reference/page
type Page struct {
	*Meta

	CreatedTime    Time                 `json:"created_time"`
	LastEditedTime Time                 `json:"last_edited_time"`
	Archived       bool                 `json:"archived"`
	Properties     map[string]*Property `json:"properties"`
}

type PageList struct {
	*ListMeta
	Results []*Page `json:"results"`
}

type Filter struct {
	// Compound filter
	Or  []*Filter `json:"or,omitempty"`
	And []*Filter `json:"and,omitempty"`

	// Database property filter
	Property    string             `json:"property,omitempty"`
	Text        *TextFilter        `json:"text,omitempty"`
	Number      *NumberFilter      `json:"number,omitempty"`
	Checkbox    *CheckboxFilter    `json:"checkbox,omitempty"`
	Select      *SelectFilter      `json:"select,omitempty"`
	MultiSelect *MultiSelectFilter `json:"multi_select,omitempty"`
	Date        *DateFilter        `json:"date,omitempty"`
	People      *PeopleFilter      `json:"people,omitempty"`
	Files       *FilesFilter       `json:"files,omitempty"`
	Relation    *RelationFilter    `json:"relation,omitempty"`
	Formula     *FormulaFilter     `json:"formula,omitempty"`
}

type TextFilter struct {
	Equals         string `json:"equals,omitempty"`
	DoesNotEqual   string `json:"does_not_equal,omitempty"`
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	StartsWith     string `json:"starts_with,omitempty"`
	EndsWith       string `json:"ends_with,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type NumberFilter struct {
	Equals               int  `json:"equals,omitempty"`
	DoesNotEqual         int  `json:"does_not_equal,omitempty"`
	GreaterThan          int  `json:"greater_than,omitempty"`
	LessThan             int  `json:"less_than,omitempty"`
	GreaterThanOrEqualTo int  `json:"greater_than_or_equal_to,omitempty"`
	LessThanOrEqualTo    int  `json:"less_than_or_equal_to,omitempty"`
	IsEmpty              bool `json:"is_empty,omitempty"`
	IsNotEmpty           bool `json:"is_not_empty,omitempty"`
}

type CheckboxFilter struct {
	Equals       bool `json:"equals,omitempty"`
	DoesNotEqual bool `json:"does_not_equal,omitempty"`
}

type SelectFilter struct {
	Equals       string `json:"equals,omitempty"`
	DoesNotEqual string `json:"does_not_equal,omitempty"`
	IsEmpty      bool   `json:"is_empty,omitempty"`
	IsNotEmpty   bool   `json:"is_not_empty,omitempty"`
}

type MultiSelectFilter struct {
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type DateFilter struct {
	Equals     *Time     `json:"equals,omitempty"`
	Before     *Time     `json:"before,omitempty"`
	After      *Time     `json:"after,omitempty"`
	OnOrBefore *Time     `json:"on_or_before,omitempty"`
	OnOrAfter  *Time     `json:"on_or_after,omitempty"`
	IsEmpty    bool      `json:"is_empty,omitempty"`
	IsNotEmpty bool      `json:"is_not_empty,omitempty"`
	PastWeek   *struct{} `json:"past_week,omitempty"`
	PastMonth  *struct{} `json:"past_month,omitempty"`
	PastYear   *struct{} `json:"past_year,omitempty"`
	NextWeek   *struct{} `json:"next_week,omitempty"`
	NextMonth  *struct{} `json:"next_month,omitempty"`
	NextYear   *struct{} `json:"next_year,omitempty"`
}

type PeopleFilter struct {
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type FilesFilter struct {
	IsEmpty    bool `json:"is_empty,omitempty"`
	IsNotEmpty bool `json:"is_not_empty,omitempty"`
}

type RelationFilter struct {
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type FormulaFilter struct {
	Text     *TextFilter     `json:"text,omitempty"`
	Checkbox *CheckboxFilter `json:"checkbox,omitempty"`
	Number   *NumberFilter   `json:"number,omitempty"`
	Date     *DateFilter     `json:"date,omitempty"`
}

type Sort struct {
	Property  string `json:"property"`
	Timestamp string `json:"timestamp"`
	Direction string `json:"direction"`
}
