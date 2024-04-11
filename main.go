package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please add an URL to the Medium article!")
		os.Exit(1)
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

	fmt.Printf("File saved in %s\n", outFilename)
}
