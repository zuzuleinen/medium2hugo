package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var images = []struct {
	url      string
	expected string
}{
	{
		url:      "https://miro.medium.com/v2/resize:fit:640/1*-KZONqGNNwqPJ4Bmf70o-Q.png",
		expected: "1*-KZONqGNNwqPJ4Bmf70o-Q.png",
	},
	{
		url:      "https://miro.medium.com/v2/resize:fit:640/1*M7AlXdGEfds9uD8fgkU6qw.png",
		expected: "1*M7AlXdGEfds9uD8fgkU6qw.png",
	},
}

func TestExtractFilename(t *testing.T) {
	for _, v := range images {
		assert.Equal(t, v.expected, extractFilename(v.url))
	}
}
