package notion

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListUsers(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/list-users.json")
	require.NoError(t, err)
	rt.RegisterResponder(
		http.MethodGet,
		"/v1/users",
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	users, err := client.ListAllUsers(context.Background())
	require.NoError(t, err)
	require.Len(t, users, 2)

	assert.Equal(t, "person", users[0].Type)
	require.NotNil(t, users[0].Person)
	assert.Equal(t, "foo@example.com", users[0].Person.Email)

	assert.Equal(t, "bot", users[1].Type)
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/get-user.json")
	require.NoError(t, err)
	rt.RegisterRegexpResponder(
		http.MethodGet,
		regexp.MustCompile(`/v1/users/[0-9a-z-]{36}$`),
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	user, err := client.GetUser(context.Background(), "2d2f95c8-c1b6-4ce1-88be-47b5b4e876e7")
	require.NoError(t, err)

	assert.Equal(t, "2d2f95c8-c1b6-4ce1-88be-47b5b4e876e7", user.ID)
	assert.Equal(t, "user", user.Object)
	assert.Equal(t, "Foo Bar", user.Name)
	assert.Equal(t, "https://lh4.googleusercontent.com/baz/AAAAAAAAAAI/AAAAAAAAAAA/foobar/photo.jpg", user.AvatarURL)
	assert.Equal(t, "person", user.Type)
	require.NotNil(t, user.Person)
	assert.Equal(t, "foo@example.com", user.Person.Email)
	assert.Nil(t, user.Bot)
}

func TestListDatabases(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/list-databases.json")
	require.NoError(t, err)
	rt.RegisterResponder(
		http.MethodGet,
		"/v1/databases",
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	databases, err := client.ListDatabases(context.Background())
	require.NoError(t, err)

	require.Len(t, databases, 1)
	db := databases[0]
	assert.Equal(t, "a4f18e20-365d-4fe1-91e8-080381f877d5", db.ID)
	assert.Equal(t, "database", db.Object)
	assert.Equal(t, int64(1621068182), db.CreatedTime.Unix())
	assert.Equal(t, int64(1621068420), db.LastEditedTime.Unix())
	if assert.Len(t, db.Title, 1) {
		assert.Equal(t, "text", db.Title[0].Type)
		if assert.NotNil(t, db.Title[0].Text) {
			assert.Equal(t, "For development", db.Title[0].Text.Content)
			assert.Nil(t, db.Title[0].Text.Link)
			assert.Equal(t, "For development", db.Title[0].PlainText)
		}
	}
	assert.Len(t, db.Properties, 18)

	require.NotNil(t, db.Properties["Name"])
	require.NotNil(t, db.Properties["Tags"])
	require.NotNil(t, db.Properties["Test1"])
	require.NotNil(t, db.Properties["Test2"])
	require.NotNil(t, db.Properties["Test3"])
	require.NotNil(t, db.Properties["Test4"])
	require.NotNil(t, db.Properties["Test5"])
	require.NotNil(t, db.Properties["Test6"])
	require.NotNil(t, db.Properties["Test7"])
	require.NotNil(t, db.Properties["Test8"])
	require.NotNil(t, db.Properties["Test9"])
	require.NotNil(t, db.Properties["Test10"])
	require.NotNil(t, db.Properties["Test11"])
	// TODO: This is probably bug of Notion.
	//require.NotNil(t, db.Properties["Test12"])
	// TODO: This is probably bug of Notion.
	//require.NotNil(t, db.Properties["Test13"])
	require.NotNil(t, db.Properties["Test14"])
	require.NotNil(t, db.Properties["Test15"])
	require.NotNil(t, db.Properties["Test16"])
	require.NotNil(t, db.Properties["Test17"])
	require.NotNil(t, db.Properties["Test18"])

	assert.Equal(t, "title", db.Properties["Name"].Type)
	assert.NotNil(t, db.Properties["Name"].Title)

	assert.Equal(t, "multi_select", db.Properties["Tags"].Type)
	assert.NotNil(t, db.Properties["Tags"].MultiSelect)

	assert.Equal(t, "rich_text", db.Properties["Test1"].Type)
	assert.NotNil(t, db.Properties["Test1"].RichText)

	assert.Equal(t, "number", db.Properties["Test2"].Type)
	if assert.NotNil(t, db.Properties["Test2"].Number) {
		assert.Equal(t, "number", db.Properties["Test2"].Number.Format)
	}

	assert.Equal(t, "select", db.Properties["Test3"].Type)
	if assert.NotNil(t, db.Properties["Test3"].Select) {
		if assert.Len(t, db.Properties["Test3"].Select.Options, 1) {
			assert.Equal(t, "3e3c5d58-4313-439e-a46e-cfaacc843d09", db.Properties["Test3"].Select.Options[0].ID)
			assert.Equal(t, "Baz", db.Properties["Test3"].Select.Options[0].Name)
			assert.Equal(t, "gray", db.Properties["Test3"].Select.Options[0].Color)
		}
	}

	assert.Equal(t, "multi_select", db.Properties["Test4"].Type)
	if assert.NotNil(t, db.Properties["Test4"].MultiSelect) {
		if assert.Len(t, db.Properties["Test4"].MultiSelect.Options, 1) {
			assert.Equal(t, "3fe82728-0646-45db-89cf-025ac1a20f02", db.Properties["Test4"].MultiSelect.Options[0].ID)
			assert.Equal(t, "Notion", db.Properties["Test4"].MultiSelect.Options[0].Name)
			assert.Equal(t, "red", db.Properties["Test4"].MultiSelect.Options[0].Color)
		}
	}

	assert.Equal(t, "date", db.Properties["Test5"].Type)
	assert.NotNil(t, db.Properties["Test5"].Date)

	assert.Equal(t, "people", db.Properties["Test6"].Type)
	assert.NotNil(t, db.Properties["Test6"].People)

	assert.Equal(t, "files", db.Properties["Test7"].Type)
	assert.NotNil(t, db.Properties["Test7"].Files)

	assert.Equal(t, "checkbox", db.Properties["Test8"].Type)
	assert.NotNil(t, db.Properties["Test8"].Checkbox)

	assert.Equal(t, "url", db.Properties["Test9"].Type)
	assert.NotNil(t, db.Properties["Test9"].URL)

	assert.Equal(t, "email", db.Properties["Test10"].Type)
	assert.NotNil(t, db.Properties["Test10"].Email)

	assert.Equal(t, "phone_number", db.Properties["Test11"].Type)
	assert.NotNil(t, db.Properties["Test11"].PhoneNumber)

	assert.Equal(t, "created_time", db.Properties["Test15"].Type)
	assert.NotNil(t, db.Properties["Test15"].CreatedTime)

	assert.Equal(t, "created_by", db.Properties["Test16"].Type)
	assert.NotNil(t, db.Properties["Test16"].CreatedBy)

	assert.Equal(t, "last_edited_time", db.Properties["Test17"].Type)
	assert.NotNil(t, db.Properties["Test17"].LastEditedTime)

	assert.Equal(t, "last_edited_by", db.Properties["Test18"].Type)
	assert.NotNil(t, db.Properties["Test18"].LastEditedBy)
}

func TestGetDatabase(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/get-database.json")
	require.NoError(t, err)
	rt.RegisterRegexpResponder(
		http.MethodGet,
		regexp.MustCompile(`/v1/databases/[a-z0-9-]{36}$`),
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	db, err := client.GetDatabase(context.Background(), "a4f18e20-365d-4fe1-91e8-080381f877d5")
	require.NoError(t, err)

	assert.Equal(t, "ba8e1263-af24-4cd0-87e0-6e2933303b60", db.ID)
	assert.Equal(t, "database", db.Object)
	assert.Equal(t, int64(1621068180), db.CreatedTime.Unix())
	assert.Equal(t, int64(1621079100), db.LastEditedTime.Unix())

	if assert.Len(t, db.Title, 1) {
		assert.Equal(t, "text", db.Title[0].Type)
		if assert.NotNil(t, db.Title[0].Text) {
			assert.Equal(t, "For development", db.Title[0].Text.Content)
			assert.Nil(t, db.Title[0].Text.Link)
			assert.Equal(t, "For development", db.Title[0].PlainText)
		}
	}
	assert.Len(t, db.Properties, 20)

	require.NotNil(t, db.Properties["Name"])
	require.NotNil(t, db.Properties["Tags"])
	require.NotNil(t, db.Properties["Test1"])
	require.NotNil(t, db.Properties["Test2"])
	require.NotNil(t, db.Properties["Test3"])
	require.NotNil(t, db.Properties["Test4"])
	require.NotNil(t, db.Properties["Test5"])
	require.NotNil(t, db.Properties["Test6"])
	require.NotNil(t, db.Properties["Test7"])
	require.NotNil(t, db.Properties["Test8"])
	require.NotNil(t, db.Properties["Test9"])
	require.NotNil(t, db.Properties["Test10"])
	require.NotNil(t, db.Properties["Test11"])
	require.NotNil(t, db.Properties["Test12"])
	require.NotNil(t, db.Properties["Test13"])
	require.NotNil(t, db.Properties["Test14"])
	require.NotNil(t, db.Properties["Test15"])
	require.NotNil(t, db.Properties["Test16"])
	require.NotNil(t, db.Properties["Test17"])
	require.NotNil(t, db.Properties["Test18"])

	assert.Equal(t, "title", db.Properties["Name"].Type)
	assert.NotNil(t, db.Properties["Name"].Title)
	assert.Equal(t, "Name", db.Properties["Name"].Name)

	assert.Equal(t, "multi_select", db.Properties["Tags"].Type)
	assert.NotNil(t, db.Properties["Tags"].MultiSelect)
	assert.Equal(t, "Tags", db.Properties["Tags"].Name)

	assert.Equal(t, "rich_text", db.Properties["Test1"].Type)
	assert.NotNil(t, db.Properties["Test1"].RichText)
	assert.Equal(t, "Test1", db.Properties["Test1"].Name)

	assert.Equal(t, "number", db.Properties["Test2"].Type)
	if assert.NotNil(t, db.Properties["Test2"].Number) {
		assert.Equal(t, "number", db.Properties["Test2"].Number.Format)
	}
	assert.Equal(t, "Test2", db.Properties["Test2"].Name)

	assert.Equal(t, "select", db.Properties["Test3"].Type)
	if assert.NotNil(t, db.Properties["Test3"].Select) {
		if assert.Len(t, db.Properties["Test3"].Select.Options, 1) {
			assert.Equal(t, "3e3c5d58-4313-439e-a46e-cfaacc843d09", db.Properties["Test3"].Select.Options[0].ID)
			assert.Equal(t, "Baz", db.Properties["Test3"].Select.Options[0].Name)
			assert.Equal(t, "gray", db.Properties["Test3"].Select.Options[0].Color)
		}
	}
	assert.Equal(t, "Test3", db.Properties["Test3"].Name)

	assert.Equal(t, "multi_select", db.Properties["Test4"].Type)
	if assert.NotNil(t, db.Properties["Test4"].MultiSelect) {
		if assert.Len(t, db.Properties["Test4"].MultiSelect.Options, 1) {
			assert.Equal(t, "3fe82728-0646-45db-89cf-025ac1a20f02", db.Properties["Test4"].MultiSelect.Options[0].ID)
			assert.Equal(t, "Notion", db.Properties["Test4"].MultiSelect.Options[0].Name)
			assert.Equal(t, "red", db.Properties["Test4"].MultiSelect.Options[0].Color)
		}
	}
	assert.Equal(t, "Test4", db.Properties["Test4"].Name)

	assert.Equal(t, "date", db.Properties["Test5"].Type)
	assert.NotNil(t, db.Properties["Test5"].Date)
	assert.Equal(t, "Test5", db.Properties["Test5"].Name)

	assert.Equal(t, "people", db.Properties["Test6"].Type)
	assert.NotNil(t, db.Properties["Test6"].People)
	assert.Equal(t, "Test6", db.Properties["Test6"].Name)

	assert.Equal(t, "files", db.Properties["Test7"].Type)
	assert.NotNil(t, db.Properties["Test7"].Files)
	assert.Equal(t, "Test7", db.Properties["Test7"].Name)

	assert.Equal(t, "checkbox", db.Properties["Test8"].Type)
	assert.NotNil(t, db.Properties["Test8"].Checkbox)
	assert.Equal(t, "Test8", db.Properties["Test8"].Name)

	assert.Equal(t, "url", db.Properties["Test9"].Type)
	assert.NotNil(t, db.Properties["Test9"].URL)
	assert.Equal(t, "Test9", db.Properties["Test9"].Name)

	assert.Equal(t, "email", db.Properties["Test10"].Type)
	assert.NotNil(t, db.Properties["Test10"].Email)
	assert.Equal(t, "Test10", db.Properties["Test10"].Name)

	assert.Equal(t, "phone_number", db.Properties["Test11"].Type)
	assert.NotNil(t, db.Properties["Test11"].PhoneNumber)
	assert.Equal(t, "Test11", db.Properties["Test11"].Name)

	assert.Equal(t, "created_time", db.Properties["Test15"].Type)
	assert.NotNil(t, db.Properties["Test15"].CreatedTime)
	assert.Equal(t, "Test15", db.Properties["Test15"].Name)

	assert.Equal(t, "rollup", db.Properties["Test14"].Type)
	if assert.NotNil(t, db.Properties["Test14"].Rollup) {
		assert.Equal(t, "Name", db.Properties["Test14"].Rollup.Name)
		assert.Equal(t, "Test13", db.Properties["Test14"].Rollup.Relation)
		assert.Equal(t, "title", db.Properties["Test14"].Rollup.RollupPropertyID)
		assert.Equal(t, "MqLD", db.Properties["Test14"].Rollup.RelationPropertyID)
		assert.Equal(t, RollupFunctionShowOriginal, db.Properties["Test14"].Rollup.Function)
	}
	assert.Equal(t, "Test14", db.Properties["Test14"].Name)

	assert.Equal(t, "created_by", db.Properties["Test16"].Type)
	assert.NotNil(t, db.Properties["Test16"].CreatedBy)
	assert.Equal(t, "Test16", db.Properties["Test16"].Name)

	assert.Equal(t, "last_edited_time", db.Properties["Test17"].Type)
	assert.NotNil(t, db.Properties["Test17"].LastEditedTime)
	assert.Equal(t, "Test17", db.Properties["Test17"].Name)

	assert.Equal(t, "last_edited_by", db.Properties["Test18"].Type)
	assert.NotNil(t, db.Properties["Test18"].LastEditedBy)
	assert.Equal(t, "Test18", db.Properties["Test18"].Name)
}

func TestGetPages(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/post-database-query.json")
	require.NoError(t, err)
	rt.RegisterRegexpResponder(
		http.MethodPost,
		regexp.MustCompile(`/v1/databases/[a-z0-9-]{36}/query`),
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	pages, err := client.GetPages(context.Background(), "a4f18e20-365d-4fe1-91e8-080381f877d5", nil, nil)
	require.NoError(t, err)

	require.Len(t, pages, 1)
	page := pages[0]
	assert.Equal(t, "16493215-50a8-41b8-8b43-0a0c014a7910", page.ID)
	assert.Equal(t, "page", page.Object)
	assert.Equal(t, int64(1621070820), page.CreatedTime.Unix())
	assert.Equal(t, int64(1621079040), page.LastEditedTime.Unix())
	assert.False(t, page.Archived)
	if assert.NotNil(t, page.Parent) {
		assert.Equal(t, page.Parent.Type, "database_id")
		assert.Equal(t, page.Parent.DatabaseID, "a4f18e20-365d-4fe1-91e8-080381f877d5")
	}
	assert.Len(t, page.Properties, 19)

	require.NotNil(t, page.Properties["Name"])
	require.NotNil(t, page.Properties["Tags"])
	require.NotNil(t, page.Properties["Test1"])
	require.NotNil(t, page.Properties["Test2"])
	require.NotNil(t, page.Properties["Test3"])
	require.NotNil(t, page.Properties["Test4"])
	require.NotNil(t, page.Properties["Test5"])
	require.NotNil(t, page.Properties["Test6"])
	require.NotNil(t, page.Properties["Test7"])
	require.NotNil(t, page.Properties["Test8"])
	require.NotNil(t, page.Properties["Test9"])
	require.NotNil(t, page.Properties["Test10"])
	require.NotNil(t, page.Properties["Test11"])
	require.NotNil(t, page.Properties["Test12"])
	require.NotNil(t, page.Properties["Test13"])
	// TODO: This is probably bug of Notion.
	//require.NotNil(t, page.Properties["Test14"])
	require.NotNil(t, page.Properties["Test15"])
	require.NotNil(t, page.Properties["Test16"])
	require.NotNil(t, page.Properties["Test17"])
	require.NotNil(t, page.Properties["Test18"])

	assert.Equal(t, "title", page.Properties["Name"].Type)
	if assert.Len(t, page.Properties["Name"].Title, 1) {
		title := page.Properties["Name"].Title[0]
		assert.Equal(t, "text", title.Type)
		if assert.NotNil(t, title.Text) {
			assert.Equal(t, "Foo", title.Text.Content)
			assert.Nil(t, title.Text.Link)
		}
		assert.Equal(t, "Foo", title.PlainText)
		assert.Empty(t, title.Href)
		if assert.NotNil(t, title.Annotations) {
			assert.False(t, title.Annotations.Bold)
			assert.False(t, title.Annotations.Italic)
			assert.False(t, title.Annotations.Strikethrough)
			assert.False(t, title.Annotations.Underline)
			assert.False(t, title.Annotations.Code)
			assert.Equal(t, "default", title.Annotations.Color)
		}
	}

	assert.Equal(t, "multi_select", page.Properties["Tags"].Type)
	if assert.Len(t, page.Properties["Tags"].MultiSelect, 1) {
		option := page.Properties["Tags"].MultiSelect[0]
		assert.Equal(t, "f7b8bf43-e891-49d4-bde5-9206f68f7ec4", option.ID)
		assert.Equal(t, "Foobar", option.Name)
		assert.Equal(t, "yellow", option.Color)
	}

	assert.Equal(t, "text", page.Properties["Test1"].Type)
	if assert.NotNil(t, page.Properties["Test1"].Text) {
		text := page.Properties["Test1"].Text[0]
		assert.Equal(t, "text", text.Type)
		if assert.NotNil(t, text.Text) {
			assert.Equal(t, "Field", text.Text.Content)
			assert.Nil(t, text.Text.Link)
		}
		assert.Equal(t, "Field", text.PlainText)
		assert.Empty(t, text.Href)
		assert.False(t, text.Annotations.Bold)
		assert.False(t, text.Annotations.Italic)
		assert.False(t, text.Annotations.Strikethrough)
		assert.False(t, text.Annotations.Underline)
		assert.False(t, text.Annotations.Code)
		assert.Equal(t, "default", text.Annotations.Color)
	}

	assert.Equal(t, "number", page.Properties["Test2"].Type)
	assert.Equal(t, 190, page.Properties["Test2"].Number)

	assert.Equal(t, "select", page.Properties["Test3"].Type)
	if assert.NotNil(t, page.Properties["Test3"].Select) {
		sel := page.Properties["Test3"].Select
		assert.Equal(t, "3e3c5d58-4313-439e-a46e-cfaacc843d09", sel.ID)
		assert.Equal(t, "Baz", sel.Name)
		assert.Equal(t, "gray", sel.Color)
	}

	assert.Equal(t, "multi_select", page.Properties["Test4"].Type)
	if assert.Len(t, page.Properties["Test4"].MultiSelect, 1) {
		option := page.Properties["Test4"].MultiSelect[0]
		assert.Equal(t, "3fe82728-0646-45db-89cf-025ac1a20f02", option.ID)
		assert.Equal(t, "Notion", option.Name)
		assert.Equal(t, "red", option.Color)
	}

	assert.Equal(t, "date", page.Properties["Test5"].Type)
	if assert.NotNil(t, page.Properties["Test5"].Date) {
		assert.Equal(t, "2021-05-15", page.Properties["Test5"].Date.Start.Format("2006-01-02"))
	}

	assert.Equal(t, "people", page.Properties["Test6"].Type)
	if assert.Len(t, page.Properties["Test6"].People, 1) {
		people := page.Properties["Test6"].People[0]
		assert.Equal(t, "user", people.Object)
		assert.Equal(t, "2d2f95c8-c1b6-4ce1-88be-47b5b4e876e7", people.ID)
		assert.Equal(t, "Foo Bar", people.Name)
		assert.Equal(t, "https://lh4.googleusercontent.com/baz/AAAAAAAAAAI/AAAAAAAAAAA/foobar/photo.jpg", people.AvatarURL)
		assert.Equal(t, "person", people.Type)
		if assert.NotNil(t, people.Person) {
			assert.Equal(t, "foo@example.com", people.Person.Email)
		}
	}

	assert.Equal(t, "files", page.Properties["Test7"].Type)
	if assert.Len(t, page.Properties["Test7"].Files, 1) {
		assert.Equal(t, "test.txt", page.Properties["Test7"].Files[0].Name)
	}

	assert.Equal(t, "checkbox", page.Properties["Test8"].Type)
	assert.True(t, page.Properties["Test8"].Checkbox)

	assert.Equal(t, "url", page.Properties["Test9"].Type)
	assert.Equal(t, "https://example.com", page.Properties["Test9"].URL)

	assert.Equal(t, "email", page.Properties["Test10"].Type)
	assert.Equal(t, "foo@example.com", page.Properties["Test10"].Email)

	assert.Equal(t, "phone_number", page.Properties["Test11"].Type)
	assert.Equal(t, "+81-23-4567-8901", page.Properties["Test11"].PhoneNumber)

	assert.Equal(t, "formula", page.Properties["Test12"].Type)
	if assert.NotNil(t, page.Properties["Test12"].Formula) {
		assert.Equal(t, "string", page.Properties["Test12"].Formula.Type)
		assert.Equal(t, "Foo", page.Properties["Test12"].Formula.String)
	}

	assert.Equal(t, "relation", page.Properties["Test13"].Type)
	assert.Len(t, page.Properties["Test13"].Relation, 0)

	assert.Equal(t, "created_time", page.Properties["Test15"].Type)
	assert.Equal(t, int64(1621070820), page.Properties["Test15"].CreatedTime.Unix())

	assert.Equal(t, "created_by", page.Properties["Test16"].Type)
	if assert.NotNil(t, page.Properties["Test16"]) {
		assert.Equal(t, "2d2f95c8-c1b6-4ce1-88be-47b5b4e876e7", page.Properties["Test16"].CreatedBy.ID)
	}

	assert.Equal(t, "last_edited_time", page.Properties["Test17"].Type)
	assert.Equal(t, int64(1621079040), page.Properties["Test17"].LastEditedTime.Unix())

	assert.Equal(t, "last_edited_by", page.Properties["Test18"].Type)
	assert.Equal(t, "2d2f95c8-c1b6-4ce1-88be-47b5b4e876e7", page.Properties["Test18"].LastEditedBy.ID)
}

func TestGetPage(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/get-page.json")
	require.NoError(t, err)
	rt.RegisterRegexpResponder(
		http.MethodGet,
		regexp.MustCompile(`/v1/pages/[a-z0-9-]{36}`),
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	page, err := client.GetPage(context.Background(), "16493215-50a8-41b8-8b43-0a0c014a7910")
	require.NoError(t, err)

	assert.Equal(t, "page", page.Object)
	assert.Equal(t, "16493215-50a8-41b8-8b43-0a0c014a7910", page.ID)
	assert.Equal(t, int64(1621070820), page.CreatedTime.Unix())
	assert.Equal(t, int64(1621079040), page.LastEditedTime.Unix())
	if assert.NotNil(t, page.Parent) {
		assert.Equal(t, "ba8e1263-af24-4cd0-87e0-6e2933303b60", page.Parent.DatabaseID)
	}
	assert.False(t, page.Archived)
	assert.Equal(t, "https://www.notion.so/Foo-1649321550a841b88b430a0c014a7910", page.URL)
	assert.Len(t, page.Properties, 20)

	require.NotNil(t, page.Properties["Name"])
	require.NotNil(t, page.Properties["Tags"])
	require.NotNil(t, page.Properties["Test1"])
	require.NotNil(t, page.Properties["Test2"])
	require.NotNil(t, page.Properties["Test3"])
	require.NotNil(t, page.Properties["Test4"])
	require.NotNil(t, page.Properties["Test5"])
	require.NotNil(t, page.Properties["Test6"])
	require.NotNil(t, page.Properties["Test7"])
	require.NotNil(t, page.Properties["Test8"])
	require.NotNil(t, page.Properties["Test9"])
	require.NotNil(t, page.Properties["Test10"])
	require.NotNil(t, page.Properties["Test11"])
	require.NotNil(t, page.Properties["Test12"])
	require.NotNil(t, page.Properties["Test13"])
	require.NotNil(t, page.Properties["Test14"])
	require.NotNil(t, page.Properties["Test15"])
	require.NotNil(t, page.Properties["Test16"])
	require.NotNil(t, page.Properties["Test17"])
	require.NotNil(t, page.Properties["Test18"])

	assert.Equal(t, "title", page.Properties["Name"].Type)
	if assert.Len(t, page.Properties["Name"].Title, 1) {
		title := page.Properties["Name"].Title[0]
		assert.Equal(t, "text", title.Type)
		if assert.NotNil(t, title.Text) {
			assert.Equal(t, "Foo", title.Text.Content)
			assert.Nil(t, title.Text.Link)
		}
		assert.Equal(t, "Foo", title.PlainText)
		assert.Empty(t, title.Href)
		if assert.NotNil(t, title.Annotations) {
			assert.False(t, title.Annotations.Bold)
			assert.False(t, title.Annotations.Italic)
			assert.False(t, title.Annotations.Strikethrough)
			assert.False(t, title.Annotations.Underline)
			assert.False(t, title.Annotations.Code)
			assert.Equal(t, "default", title.Annotations.Color)
		}
	}

	assert.Equal(t, "multi_select", page.Properties["Tags"].Type)
	if assert.Len(t, page.Properties["Tags"].MultiSelect, 1) {
		option := page.Properties["Tags"].MultiSelect[0]
		assert.Equal(t, "f7b8bf43-e891-49d4-bde5-9206f68f7ec4", option.ID)
		assert.Equal(t, "Foobar", option.Name)
		assert.Equal(t, "yellow", option.Color)
	}

	assert.Equal(t, "text", page.Properties["Test1"].Type)
	if assert.NotNil(t, page.Properties["Test1"].Text) {
		text := page.Properties["Test1"].Text[0]
		assert.Equal(t, "text", text.Type)
		if assert.NotNil(t, text.Text) {
			assert.Equal(t, "Field", text.Text.Content)
			assert.Nil(t, text.Text.Link)
		}
		assert.Equal(t, "Field", text.PlainText)
		assert.Empty(t, text.Href)
		assert.False(t, text.Annotations.Bold)
		assert.False(t, text.Annotations.Italic)
		assert.False(t, text.Annotations.Strikethrough)
		assert.False(t, text.Annotations.Underline)
		assert.False(t, text.Annotations.Code)
		assert.Equal(t, "default", text.Annotations.Color)
	}

	assert.Equal(t, "number", page.Properties["Test2"].Type)
	assert.Equal(t, 190, page.Properties["Test2"].Number)

	assert.Equal(t, "select", page.Properties["Test3"].Type)
	if assert.NotNil(t, page.Properties["Test3"].Select) {
		sel := page.Properties["Test3"].Select
		assert.Equal(t, "3e3c5d58-4313-439e-a46e-cfaacc843d09", sel.ID)
		assert.Equal(t, "Baz", sel.Name)
		assert.Equal(t, "gray", sel.Color)
	}

	assert.Equal(t, "multi_select", page.Properties["Test4"].Type)
	if assert.Len(t, page.Properties["Test4"].MultiSelect, 1) {
		option := page.Properties["Test4"].MultiSelect[0]
		assert.Equal(t, "3fe82728-0646-45db-89cf-025ac1a20f02", option.ID)
		assert.Equal(t, "Notion", option.Name)
		assert.Equal(t, "red", option.Color)
	}

	assert.Equal(t, "date", page.Properties["Test5"].Type)
	if assert.NotNil(t, page.Properties["Test5"].Date) {
		assert.Equal(t, "2021-05-15", page.Properties["Test5"].Date.Start.Format("2006-01-02"))
	}

	assert.Equal(t, "people", page.Properties["Test6"].Type)
	if assert.Len(t, page.Properties["Test6"].People, 1) {
		people := page.Properties["Test6"].People[0]
		assert.Equal(t, "user", people.Object)
		assert.Equal(t, "2d2f95c8-c1b6-4ce1-88be-47b5b4e876e7", people.ID)
		assert.Equal(t, "Foo Bar", people.Name)
		assert.Equal(t, "https://lh4.googleusercontent.com/baz/AAAAAAAAAAI/AAAAAAAAAAA/foobar/photo.jpg", people.AvatarURL)
		assert.Equal(t, "person", people.Type)
		if assert.NotNil(t, people.Person) {
			assert.Equal(t, "foo@example.com", people.Person.Email)
		}
	}

	assert.Equal(t, "files", page.Properties["Test7"].Type)
	if assert.Len(t, page.Properties["Test7"].Files, 1) {
		assert.Equal(t, "test.txt", page.Properties["Test7"].Files[0].Name)
	}

	assert.Equal(t, "checkbox", page.Properties["Test8"].Type)
	assert.True(t, page.Properties["Test8"].Checkbox)

	assert.Equal(t, "url", page.Properties["Test9"].Type)
	assert.Equal(t, "https://example.com", page.Properties["Test9"].URL)

	assert.Equal(t, "email", page.Properties["Test10"].Type)
	assert.Equal(t, "foo@example.com", page.Properties["Test10"].Email)

	assert.Equal(t, "phone_number", page.Properties["Test11"].Type)
	assert.Equal(t, "+81-23-4567-8901", page.Properties["Test11"].PhoneNumber)

	assert.Equal(t, "formula", page.Properties["Test12"].Type)
	if assert.NotNil(t, page.Properties["Test12"].Formula) {
		assert.Equal(t, "string", page.Properties["Test12"].Formula.Type)
		assert.Equal(t, "Foo", page.Properties["Test12"].Formula.String)
	}

	assert.Equal(t, "relation", page.Properties["Test13"].Type)
	assert.Len(t, page.Properties["Test13"].Relation, 0)

	assert.Equal(t, "created_time", page.Properties["Test15"].Type)
	assert.Equal(t, int64(1621070820), page.Properties["Test15"].CreatedTime.Unix())

	assert.Equal(t, "created_by", page.Properties["Test16"].Type)
	if assert.NotNil(t, page.Properties["Test16"]) {
		assert.Equal(t, "2d2f95c8-c1b6-4ce1-88be-47b5b4e876e7", page.Properties["Test16"].CreatedBy.ID)
	}

	assert.Equal(t, "last_edited_time", page.Properties["Test17"].Type)
	assert.Equal(t, int64(1621079040), page.Properties["Test17"].LastEditedTime.Unix())

	assert.Equal(t, "last_edited_by", page.Properties["Test18"].Type)
	assert.Equal(t, "2d2f95c8-c1b6-4ce1-88be-47b5b4e876e7", page.Properties["Test18"].LastEditedBy.ID)
}

func TestGetBlock(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/get-block.json")
	require.NoError(t, err)
	rt.RegisterRegexpResponder(
		http.MethodGet,
		regexp.MustCompile(`/v1/blocks/[a-z0-9-]{36}`),
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	block, err := client.GetBlock(context.Background(), "6fbac55c-9e74-4489-b386-0c88d5aa54dd")
	require.NoError(t, err)

	assert.Equal(t, "block", block.Object)
	assert.Equal(t, "6fbac55c-9e74-4489-b386-0c88d5aa54dd", block.ID)
	assert.Equal(t, int64(1621500000), block.CreatedTime.Unix())
	assert.Equal(t, int64(1621500000), block.LastEditedTime.Unix())
	assert.False(t, block.HasChildren)
	assert.False(t, block.Archived)
	assert.Equal(t, "paragraph", block.Type)
	if assert.NotNil(t, block.Paragraph) {
		if assert.Len(t, block.Paragraph.Text, 1) {
			assert.Equal(t, "text", block.Paragraph.Text[0].Type)
			assert.Equal(t, "development", block.Paragraph.Text[0].Text.Content)
			assert.Equal(t, "development", block.Paragraph.Text[0].PlainText)
		}
	}
}

func TestGetBlocks(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/get-block-children.json")
	require.NoError(t, err)
	rt.RegisterRegexpResponder(
		http.MethodGet,
		regexp.MustCompile(`/v1/blocks/[a-z0-9-]{36}/children`),
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	blocks, err := client.GetBlocks(context.Background(), "16493215-50a8-41b8-8b43-0a0c014a7910")
	require.NoError(t, err)

	require.Len(t, blocks, 10)

	assert.Equal(t, "paragraph", blocks[0].Type)
	if assert.NotNil(t, blocks[0].Paragraph) {
		if assert.Len(t, blocks[0].Paragraph.Text, 1) {
			assert.Equal(t, "text", blocks[0].Paragraph.Text[0].Type)
			assert.Equal(t, "Body text", blocks[0].Paragraph.Text[0].PlainText)
		}
	}

	assert.Equal(t, "heading_1", blocks[1].Type)
	if assert.NotNil(t, blocks[1].Heading1) {
		if assert.Len(t, blocks[1].Heading1.Text, 1) {
			assert.Equal(t, "Title1", blocks[1].Heading1.Text[0].PlainText)
		}
	}

	assert.Equal(t, "heading_2", blocks[2].Type)
	if assert.NotNil(t, blocks[2].Heading2) {
		if assert.Len(t, blocks[2].Heading2.Text, 1) {
			assert.Equal(t, "Title2", blocks[2].Heading2.Text[0].PlainText)
		}
	}

	assert.Equal(t, "heading_3", blocks[3].Type)
	if assert.NotNil(t, blocks[3].Heading3) {
		if assert.Len(t, blocks[3].Heading3.Text, 1) {
			assert.Equal(t, "Title3", blocks[3].Heading3.Text[0].PlainText)
		}
	}

	assert.Equal(t, "paragraph", blocks[4].Type)
	if assert.NotNil(t, blocks[4].Paragraph) {
		if assert.Len(t, blocks[4].Paragraph.Text, 3) {
			assert.True(t, blocks[4].Paragraph.Text[0].Annotations.Bold)
			assert.True(t, blocks[4].Paragraph.Text[1].Annotations.Italic)
		}
	}

	// Actually, this block is a divider.
	assert.Equal(t, "unsupported", blocks[5].Type)

	assert.Equal(t, "bulleted_list_item", blocks[6].Type)
	if assert.NotNil(t, blocks[6].BulletedListItem) {
		if assert.Len(t, blocks[6].BulletedListItem.Text, 1) {
			assert.Equal(t, "Bullet1", blocks[6].BulletedListItem.Text[0].PlainText)
		}
	}

	assert.Equal(t, "bulleted_list_item", blocks[7].Type)
	if assert.NotNil(t, blocks[7].BulletedListItem) {
		if assert.Len(t, blocks[7].BulletedListItem.Text, 1) {
			assert.Equal(t, "Bullet2", blocks[7].BulletedListItem.Text[0].PlainText)
		}
	}

	// Actually, this block is a code block.
	assert.Equal(t, "unsupported", blocks[8].Type)

	assert.Equal(t, "paragraph", blocks[9].Type)
	if assert.NotNil(t, blocks[9].Paragraph) {
		if assert.Len(t, blocks[9].Paragraph.Text, 2) {
			assert.True(t, blocks[9].Paragraph.Text[0].Annotations.Code)
			assert.Equal(t, " foobar", blocks[9].Paragraph.Text[1].PlainText)
		}
	}
}

func TestUpdateBlock(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/update-block.json")
	require.NoError(t, err)
	rt.RegisterRegexpResponder(
		http.MethodPatch,
		regexp.MustCompile(`/v1/blocks/[a-z0-9-]{36}`),
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	block, err := client.UpdateBlock(context.Background(), &Block{Meta: &Meta{ID: "cdfb0555-29e4-4bad-baaa-240a0097c77d"}})
	require.NoError(t, err)

	assert.Equal(t, "block", block.Object)
	assert.Equal(t, "cdfb0555-29e4-4bad-baaa-240a0097c77d", block.ID)
}

func TestCreatePage(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/post-page.json")
	require.NoError(t, err)
	rt.RegisterRegexpResponder(
		http.MethodPost,
		regexp.MustCompile(`/v1/pages$`),
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	page, err := client.CreatePage(context.Background(), &Page{})
	require.NoError(t, err)

	assert.Equal(t, "9585d9b5-ad82-4221-9f82-a3a4767d5b92", page.ID)
	assert.Equal(t, int64(1621158331), page.CreatedTime.Unix())
	assert.Equal(t, int64(1621158331), page.LastEditedTime.Unix())
}

func TestPatchPage(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/patch-page.json")
	require.NoError(t, err)
	rt.RegisterRegexpResponder(
		http.MethodPatch,
		regexp.MustCompile(`/v1/pages/[0-9a-z-]{36}$`),
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	properties := map[string]*PropertyData{
		"Test1": {
			Type: "rich_text",
			RichText: []*RichTextObject{
				{
					Type: "text",
					Text: &Text{
						Content: "Update property",
					},
				},
			},
		},
	}
	page, err := client.UpdateProperties(context.Background(), "9585d9b5-ad82-4221-9f82-a3a4767d5b92", properties)
	require.NoError(t, err)

	assert.Equal(t, "9585d9b5-ad82-4221-9f82-a3a4767d5b92", page.ID)
	if assert.NotNil(t, page.Properties["Test1"]) &&
		assert.Len(t, page.Properties["Test1"].RichText, 1) &&
		assert.NotNil(t, page.Properties["Test1"].RichText[0].Text) {
		assert.Equal(t, "Update property", page.Properties["Test1"].RichText[0].Text.Content)
	}
}

func TestPatchBlockChildren(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/patch-block-children.json")
	require.NoError(t, err)
	rt.RegisterRegexpResponder(
		http.MethodPatch,
		regexp.MustCompile(`/v1/blocks/[0-9a-z-]{36}/children$`),
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	block, err := client.AppendBlock(context.Background(), "9585d9b5-ad82-4221-9f82-a3a4767d5b92", []*Block{
		{
			Meta: &Meta{
				Object: "block",
			},
			Type: "paragraph",
			Paragraph: &Paragraph{
				Text: []*RichTextObject{
					{Type: "text", Text: &Text{Content: "Good"}},
				},
			},
		},
	})
	require.NoError(t, err)

	assert.Equal(t, "80434f64-2d08-414d-83b0-5fca519794bd", block.ID)
	assert.True(t, block.HasChildren)
}

func TestPostSearch(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/post-search.json")
	require.NoError(t, err)
	rt.RegisterRegexpResponder(
		http.MethodPost,
		regexp.MustCompile(`/v1/search$`),
		httpmock.NewStringResponder(
			http.StatusOK,
			string(res),
		),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	results, err := client.Search(context.Background(), "q", nil)
	require.NoError(t, err)

	require.Len(t, results, 2)
	if assert.IsType(t, &Page{}, results[0]) {
		page := results[0].(*Page)
		assert.Equal(t, "page", page.Object)
		assert.Equal(t, "04df63cf-8140-481a-81e0-8f0af77131b2", page.ID)
		if assert.NotNil(t, page.Parent) {
			assert.Equal(t, "ba8e1263-af24-4cd0-87e0-6e2933303b60", page.Parent.DatabaseID)
		}
		assert.Len(t, page.Properties, 13)
	}
	if assert.IsType(t, &Database{}, results[1]) {
		db := results[1].(*Database)
		assert.Equal(t, "database", db.Object)
		assert.Equal(t, "ba8e1263-af24-4cd0-87e0-6e2933303b60", db.ID)
		assert.Len(t, db.Properties, 18)
	}
}

func TestCreateDatabase(t *testing.T) {
	t.Parallel()

	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("testdata/post-databases.json")
	require.NoError(t, err)
	rt.RegisterRegexpResponder(
		http.MethodPost,
		regexp.MustCompile(`/v1/databases$`),
		httpmock.NewStringResponder(http.StatusOK, string(res)),
	)

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)

	db, err := client.CreateDatabase(context.Background(), &Database{})
	require.NoError(t, err)

	assert.Equal(t, "181a22c0-66af-439a-8265-c2473a59fee9", db.ID)
	assert.Equal(t, "database", db.Object)
	assert.Equal(t, int64(1627192200), db.CreatedTime.Unix())
	assert.Equal(t, int64(1627192200), db.LastEditedTime.Unix())

	if assert.Len(t, db.Title, 1) {
		assert.Equal(t, "text", db.Title[0].Type)
		if assert.NotNil(t, db.Title[0].Text) {
			assert.Equal(t, "Create database", db.Title[0].Text.Content)
		}
		assert.Equal(t, "Create database", db.Title[0].PlainText)
	}

	assert.Len(t, db.Properties, 2)

	require.NotNil(t, db.Properties["Name"])
	require.NotNil(t, db.Properties["Test1"])
}
