package main

import (
	"fmt"
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
	outFilename := parser.ArticleFilename(mediumURL)
	imgDir := "images"

	if err := parser.ExportToHugo(mediumURL, outFilename, imgDir); err != nil {
		log.Fatalf("error occured while exporting: %s", err)
	}

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("\nArticle saved in: %s\n", yellow(outFilename))
}
