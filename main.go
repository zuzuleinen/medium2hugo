package main

import (
	"log"
	"os"

	"medium2hugo/parser"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) < 2 {
		color.Red("Please add an URL")
		return
	}
	mediumURL := os.Args[1]
	outFilename := "article.md"

	if err := parser.ExportToHugo(mediumURL, outFilename); err != nil {
		log.Fatalf("error occured while exporting: %s", err)
	}

	color.Green("Article saved in: %s\n", outFilename)
}
