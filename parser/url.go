package parser

import (
	"fmt"
	"net/url"
	"strings"
)

// URLForJSON transforms originalURL into a URL from which JSON response can be fetched
// The `@author` becomes `author` and ?format=json is appended to the URL
//
// At the time of implementation http.Get will only fetch JSON response if author is stripped of @
// In the browser it works with @ as well
func URLForJSON(originalURL string) (string, error) {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", fmt.Errorf("error parsing url: %w", err)
	}

	query := parsedURL.Query()
	query.Set("format", "json")
	parsedURL.RawQuery = query.Encode()

	return strings.Replace(parsedURL.String(), "@", "", 1), nil
}
