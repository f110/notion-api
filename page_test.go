package notion

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPage_New(t *testing.T) {
	rt := mockTransport(t, http.MethodGet, `/v1/pages/[a-z0-9-]{36}`, http.StatusOK, "./testdata/page-copy.json")

	client, err := New(&http.Client{Transport: rt}, "https://example.com")
	require.NoError(t, err)
	page, err := client.GetPage(context.Background(), "56f2049d-feb1-4a3f-b227-2fa76ca74d0e")
	require.NoError(t, err)

	newPage := page.New()
	assert.Empty(t, newPage.ID)
	assert.Nil(t, newPage.CreatedTime)
	assert.Nil(t, newPage.LastEditedTime)
	assert.Contains(t, newPage.Properties, "Name")
	assert.Contains(t, newPage.Properties, "Test1")
	assert.Contains(t, newPage.Properties, "Test12")
	assert.NotContains(t, newPage.Properties, "Tags")
	assert.NotContains(t, newPage.Properties, "Test2")
	assert.NotContains(t, newPage.Properties, "Test3")
	assert.NotContains(t, newPage.Properties, "Test4")
	assert.NotContains(t, newPage.Properties, "Test5")
	assert.NotContains(t, newPage.Properties, "Test6")
	assert.NotContains(t, newPage.Properties, "Test8")
	assert.NotContains(t, newPage.Properties, "Test9")
	assert.NotContains(t, newPage.Properties, "Test10")
	assert.NotContains(t, newPage.Properties, "Test11")
	assert.NotContains(t, newPage.Properties, "Test13")
	assert.NotContains(t, newPage.Properties, "Test14")
	assert.NotContains(t, newPage.Properties, "Test15")
	assert.NotContains(t, newPage.Properties, "Test16")
	assert.NotContains(t, newPage.Properties, "Test17")
	assert.NotContains(t, newPage.Properties, "Test18")
	assert.Empty(t, newPage.URL)
}
