package parser

import (
	"bytes"
	"time"

	"github.com/pelletier/go-toml/v2"
)

var delimiter = []byte("+++\n")

// FrontMatter is used to add metadata to a Hugo article
// see https://gohugo.io/content-management/front-matter/
type FrontMatter struct {
	Title string    `toml:"title"`
	Date  time.Time `toml:"date"`
	Draft bool      `toml:"draft"`
	Tags  []string  `toml:"tags"`
}

func (f FrontMatter) String() string {
	content, _ := toml.Marshal(f)

	return string(bytes.Join([][]byte{delimiter, delimiter}, content))
}
