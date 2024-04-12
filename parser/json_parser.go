package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type JSONParser struct {
}

type Content struct {
	Title string
	Tags  []string
	Date  time.Time
}

type JSONResponse struct {
	Success bool    `json:"success"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Value Value `json:"value"`
}

type Value struct {
	Title             string   `json:"title"`
	Virtuals          Virtuals `json:"virtuals"`
	LatestPublishedAt int      `json:"latestPublishedAt"`
}

type Virtuals struct {
	Tags []Tag `json:"tags"`
}

type Tag struct {
	Name string `json:"name"`
}

func (h JSONParser) Parse(r io.Reader) (Content, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Content{}, fmt.Errorf("error reading contents: %w", err)
	}

	// json starts at the first occurrence of {
	jsonData := data[bytes.Index(data, []byte("{")):]

	var resp JSONResponse
	if err = json.Unmarshal(jsonData, &resp); err != nil {
		return Content{}, fmt.Errorf("error on JSON Unmarshal: %s", err)
	}

	tags := make([]string, len(resp.Payload.Value.Virtuals.Tags))
	for k, v := range resp.Payload.Value.Virtuals.Tags {
		tags[k] = v.Name
	}

	timestamp := int64(resp.Payload.Value.LatestPublishedAt)
	seconds := timestamp / 1000
	nanos := (timestamp % 1000) * 1000000
	t := time.Unix(seconds, nanos)

	return Content{
		Title: resp.Payload.Value.Title,
		Tags:  tags,
		Date:  t,
	}, nil
}
