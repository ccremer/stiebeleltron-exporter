package stiebeleltron

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubProperty struct {
	value        float64
	group        string
	searchString string
}

func (t stubProperty) GetGroup() string {
	return t.group
}

func (t stubProperty) GetSearchString() string {
	return t.searchString
}

func (t *stubProperty) SetValue(v float64) {
	t.value = v
}

func TestISGClient_ParsePage(t *testing.T) {
	server := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer server.Close()

	client, err := NewISGClient(ClientOptions{
		BaseURL: server.URL,
		Headers: nil,
	})
	require.NoError(t, err)
	prop := &stubProperty{
		group:        "RUNTIME",
		searchString: "RNT COMP 1 DHW",
	}
	_, err = client.ParsePage("/heatpumpinfo_1.html", []Property{prop})
	require.NoError(t, err)
	assert.Equal(t, float64(1771), prop.value)
}
