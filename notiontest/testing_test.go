package notiontest_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.f110.dev/notion-api/v3"
	"go.f110.dev/notion-api/v3/notiontest"
)

func TestMock(t *testing.T) {
	mock := notiontest.NewMock()
	mock.
		User("Alice").
		User("Bob").
		BotUser("Client").
		Database(notiontest.NewDatabase("Sample database")).
		Page(
			notiontest.NewPage("Sample page",
				notiontest.PageProperty("Col 1", &notion.PropertyData{Type: notion.PropertyTypeText, Text: []*notion.RichTextObject{{Type: notion.RichTextObjectTypeText, Text: &notion.Text{Content: "Foo"}}}}),
			),
		)

	tr := httpmock.NewMockTransport()
	mock.RegisterMock(tr)

	tc := mock.AuthenticatedClient("Client")
	client, err := notion.New(tc, "https://example.com")
	require.NoError(t, err)

	t.Run("Error", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "https://example.com/v1/users", nil)
		require.NoError(t, err)
		res, err := tc.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		buf, err := io.ReadAll(res.Body)
		var errResponse notion.Error
		err = json.Unmarshal(buf, &errResponse)
		require.NoError(t, err)
		assert.Equal(t, "missing_version", errResponse.Code)
	})

	t.Run("ListAllUsers", func(t *testing.T) {
		users, err := client.ListAllUsers(context.Background())
		require.NoError(t, err)
		assert.Len(t, users, 3)
		assert.Equal(t, "Alice", users[0].Name)
		assert.Equal(t, "Bob", users[1].Name)
		assert.Equal(t, "Client", users[2].Name)
	})

	t.Run("GetUser", func(t *testing.T) {
		userID := mock.FindUser("Alice").GetID()
		user, err := client.GetUser(context.Background(), userID)
		require.NoError(t, err)
		assert.Equal(t, "Alice", user.Name)
	})

	t.Run("GetMe", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "https://example.com/v1/users/me", nil)
		require.NoError(t, err)
		req.Header.Set("Notion-Version", "2022-06-28")
		res, err := tc.Do(req)
		require.NoError(t, err)

		var botUser notion.User
		err = json.NewDecoder(res.Body).Decode(&botUser)
		require.NoError(t, err)
		assert.Equal(t, "Client", botUser.Name)
	})

	t.Run("CreateDatabase", func(t *testing.T) {
		db, err := client.CreateDatabase(context.Background(), notiontest.NewDatabase("New database"))
		require.NoError(t, err)

		assert.NotEmpty(t, db.ID)
	})

	t.Run("GetDatabase", func(t *testing.T) {
		sampleDB := mock.FindDatabase("Sample database")
		db, err := client.GetDatabase(context.Background(), sampleDB[0].GetID())
		require.NoError(t, err)
		assert.NotNil(t, db)

		_, err = client.GetDatabase(context.Background(), "b3575be4-77e4-429c-a4da-6835721cc2ba") // Unknown ID
		require.Error(t, err)
		var nErr *notion.Error
		assert.ErrorAs(t, err, &nErr)
		assert.Equal(t, 404, nErr.Status)
	})

	t.Run("UpdateDatabase", func(t *testing.T) {
		db := mock.NewDatabase("Sample database for Update")
		db.Properties["Sub title"] = &notion.PropertyMetadata{Name: "Sub title"}
		updatedDB, err := client.UpdateDatabase(context.Background(), db)
		require.NoError(t, err)

		assert.Equal(t, db.GetID(), updatedDB.GetID())
		assert.Contains(t, updatedDB.Properties, "Sub title")
	})

	t.Run("CreatePage", func(t *testing.T) {
		page, err := client.CreatePage(context.Background(), &notion.Page{})
		require.NoError(t, err)

		assert.NotEmpty(t, page.ID)
	})

	t.Run("GetPage", func(t *testing.T) {
		samplePage := mock.FindPage("Sample page")
		page, err := client.GetPage(context.Background(), samplePage[0].GetID())
		require.NoError(t, err)
		assert.NotNil(t, page)

		_, err = client.GetPage(context.Background(), "b3575be4-77e4-429c-a4da-6835721cc2ba") // Unknown ID
		require.Error(t, err)
		var nErr *notion.Error
		assert.ErrorAs(t, err, &nErr)
		assert.Equal(t, 404, nErr.Status)
	})

	t.Run("GetPageProperty", func(t *testing.T) {
		samplePage := mock.FindPage("Sample page")
		property, err := client.GetPageProperty(context.Background(), samplePage[0].GetID(), samplePage[0].Properties["Col 1"].ID)
		require.NoError(t, err)
		assert.NotNil(t, property)

		require.NotEmpty(t, property.Text)
		require.NotNil(t, property.Text[0].Text)
		assert.Equal(t, "Foo", property.Text[0].Text.Content)
	})
}
