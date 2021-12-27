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

func TestPage_New(t *testing.T) {
	rt := httpmock.NewMockTransport()
	res, err := os.ReadFile("./testdata/page-copy.json")
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
	pages, err := client.GetPages(context.Background(), "56f2049d-feb1-4a3f-b227-2fa76ca74d0e", nil, nil)
	require.NoError(t, err)
	page := pages[0]

	newPage := page.New()
	assert.Empty(t, newPage.ID)
	assert.Nil(t, newPage.CreatedTime)
	assert.Nil(t, newPage.LastEditedTime)
	assert.Contains(t, newPage.Properties, "Schedule")
	assert.Contains(t, newPage.Properties, "Status")
	assert.Contains(t, newPage.Properties, "Name")
	assert.NotContains(t, newPage.Properties, "Updated at")
	assert.NotContains(t, newPage.Properties, "Created at")
	assert.NotContains(t, newPage.Properties, "Version")
}
