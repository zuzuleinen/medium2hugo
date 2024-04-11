package main

import (
	"log"
	"net/http"
	"os"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

func main() {
	if len(os.Args) < 2 {
		color.Red("Please add an URL")
		return
	}

	// https://medium.com/@andreiboar/fundamentals-of-i-o-in-go-part-2-e7bb68cd5608
	url := os.Args[1]

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
	selec := doc.Find("article")

	converter := md.NewConverter("", true, nil)
	markdown := converter.Convert(selec)

	outFilename := "result.md"

	f, err := os.Create(outFilename)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString(markdown)
	if err != nil {
		log.Fatal(err)
	}

	color.Green("File saved: %s\n", outFilename)
}
