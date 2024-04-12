package parser

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

// ExportToHugo exports an article from mediumURL into a Hugo compatible Markdown file
func ExportToHugo(mediumURL, outFilename, imgDir string) error {
	frontMatter, err := getFrontMatter(mediumURL)
	if err != nil {
		return fmt.Errorf("could not compute front matter: %w", err)
	}

	markDownContent, err := getMarkdownBody(mediumURL, imgDir)
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
func getMarkdownBody(url, imgDir string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("error getting url: %s", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".pw-subtitle-paragraph").Parent().Remove()
	doc.Find(".pw-post-title").Parent().Remove()
	article := doc.Find("article")

	converter := md.NewConverter("", true, nil)
	converter.Use(MediumImage(imgDir))
	return converter.Convert(article), nil
}

// MediumImage will parse images from Medium story and save them in an ./images folder
func MediumImage(imgDir string) md.Plugin {
	return func(c *md.Converter) []md.Rule {
		result := ""
		return []md.Rule{
			{
				Filter: []string{"picture"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					selec.ChildrenFiltered("source").Each(func(i int, selection *goquery.Selection) {
						if val, ok := selection.Attr("type"); val == "image/webp" && ok {
							if err := os.MkdirAll(imgDir, 0755); err != nil {
								log.Fatal(err)
							}

							if srcSet, hasSrc := selection.Attr("srcset"); hasSrc {
								imgURL := strings.Split(strings.Split(srcSet, ",")[0], " ")[0]
								imgURL = strings.Replace(imgURL, "/format:webp", "", 1)
								res, err := http.Get(imgURL)
								if err != nil {
									log.Fatal(err)
								}

								filename := extractFilename(imgURL)

								f, err := os.Create(fmt.Sprintf("%s/%s", imgDir, filename))
								if err != nil {
									log.Fatal(err)
								}

								bs, err := io.ReadAll(res.Body)
								if err != nil {
									log.Fatal(err)
								}
								if _, err = f.Write(bs); err != nil {
									log.Fatal(err)
								}

								yellow := color.New(color.FgYellow).SprintFunc()
								fmt.Printf("Exporting image: %s\n", yellow(f.Name()))

								result = fmt.Sprintf("![Image Alt](/%s/%s)", imgDir, extractFilename(imgURL))
							}
						}
					})
					return md.String(result)
				},
			},
		}
	}
}

// extractFilename extracts filename with extension from a Medium URL
// such as https://miro.medium.com/v2/resize:fit:640/1*-KZONqGNNwqPJ4Bmf70o-Q.png
func extractFilename(url string) string {
	parts := strings.Split(url, "/")
	filenameWithExt := parts[len(parts)-1]
	return filenameWithExt
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

// ArticleFilename returns the filename of generated Markdown article based on the mediumURL
func ArticleFilename(mediumURL string) string {
	parts := strings.Split(mediumURL, "/")
	return parts[len(parts)-1] + ".md"
}
