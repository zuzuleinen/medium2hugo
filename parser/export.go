package parser

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

// ExportToHugo exports an article from mediumURL into a Hugo compatible Markdown file
func ExportToHugo(mediumURL, outFilename string) error {
	frontMatter, err := getFrontMatter(mediumURL)
	if err != nil {
		return fmt.Errorf("could not compute front matter: %w", err)
	}

	markDownContent, err := getMarkdownBody(mediumURL)
	if err != nil {
		return fmt.Errorf("could not compute markdown content: %w", err)
	}

	var buf bytes.Buffer
	if _, err = buf.WriteString(frontMatter.String()); err != nil {
		return fmt.Errorf("could not write front matter to file: %w", err)
	}
	if _, err = buf.WriteString(markDownContent); err != nil {
		return fmt.Errorf("could not write markdown content to file: %w", err)
	}

	f, err := os.Create(outFilename)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	if _, err = f.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

// getMarkdownBody parses HTML and create article body into Markdown
func getMarkdownBody(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("error getting url: %s", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".pw-post-title").Parent().Remove()
	article := doc.Find("article")

	converter := md.NewConverter("", true, nil)

	return converter.Convert(article), nil
}

// getMarkdownBody parses JSON response and creates article Front Matter
func getFrontMatter(originalURl string) (FrontMatter, error) {
	url, err := URLForJSON(originalURl)
	if err != nil {
		return FrontMatter{}, fmt.Errorf("could not compute URL for json: %w", err)
	}
	res, err := http.Get(url)
	if err != nil {
		return FrontMatter{}, fmt.Errorf("could not get JSON response: %w", err)
	}
	defer res.Body.Close()

	p := JSONParser{}
	c, err := p.Parse(res.Body)
	if err != nil {
		return FrontMatter{}, fmt.Errorf("could not parse response body: %w", err)
	}
	return FrontMatter{
		Title: c.Title,
		Date:  c.Date,
		Tags:  c.Tags,
	}, nil
}
