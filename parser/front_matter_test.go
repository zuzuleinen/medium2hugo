package parser

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrontMatterString(t *testing.T) {
	date, err := time.Parse(time.RFC3339, "2024-04-10T18:43:25+03:00")
	require.NoError(t, err)

	f := FrontMatter{
		Title: "Fundamentals of I/O in Go",
		Date:  date,
		Draft: false,
		Tags:  []string{"Go", "I/O"},
	}
	expected := `+++
title = 'Fundamentals of I/O in Go'
date = 2024-04-10T18:43:25+03:00
draft = false
tags = ['Go', 'I/O']
+++
`
	assert.Equal(t, expected, f.String())
}

func TestJsonParse(t *testing.T) {
	// mock for https://medium.com/andreiboar/fundamentals-of-i-o-in-go-part-2-e7bb68cd5608?format=json
	mockResponse, err := os.Open("article.json")
	require.NoError(t, err)

	p := JSONParser{}
	c, err := p.Parse(mockResponse)
	require.NoError(t, err)

	assert.Equal(t, "Fundamentals of I/O in Go: Part 2", c.Title)
	assert.Equal(t, []string{"Golang", "Input Output", "Golang Tutorial", "Software Development", "Programming"}, c.Tags)
	assert.Equal(t, "2024-04-10T18:43:25+03:00", c.Date.Format("2006-01-02T15:04:05-07:00"))
}
