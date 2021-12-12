package parser

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAlbumMainPage(t *testing.T) {
	expected := 16

	b, err := ioutil.ReadFile("./testdata/album.html")
	assert.NoError(t, err)
	body := strings.NewReader(string(b))
	actual, err := ParseAlbumMainPage(body)
	assert.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestParseAlbumPage(t *testing.T) {
	expected := []string{
		"1546342",
		"1546343",
		"1546344",
		"1546345",
		"1546346",
		"1546347",
		"1546348",
		"1546349",
		"1546350",
		"1546351",
		"1546352",
		"1546353",
		"1546354",
		"1546356",
		"1546358",
		"1546359",
		"1546360",
		"1546361",
		"1546362",
		"1546363",
		"1546364",
		"1546365",
		"1546366",
		"1546367",
	}

	b, err := ioutil.ReadFile("./testdata/album.html")
	assert.NoError(t, err)
	body := strings.NewReader(string(b))
	actual, err := ParseAlbumPage(body)
	assert.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}
