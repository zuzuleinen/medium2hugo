package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLForJSON(t *testing.T) {
	url := "https://medium.com/@andreiboar/fundamentals-of-i-o-in-go-part-2-e7bb68cd5608"
	expected := "https://medium.com/andreiboar/fundamentals-of-i-o-in-go-part-2-e7bb68cd5608?format=json"

	finalURL, err := URLForJSON(url)
	assert.NoError(t, err)

	assert.Equal(t, expected, finalURL)
}
