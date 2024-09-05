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
		BotUser("Client")

	tr := httpmock.NewMockTransport()
	mock.RegisterMock(tr)

	//ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: mock.GenerateBotToken("Client")})
	//tc := oauth2.NewClient(context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: tr}), ts)
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
}
