package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

type BlogEntry struct {
	id          string
	date        string
	content     string
	attachments []string
}

var mdConverter = md.NewConverter("", true, nil)

func writeEntry(entry BlogEntry, outputFolder string) {
	entryStr := fmt.Sprintf(`
+++
author = "Roberto Perez"
title = "%s"
date = "%s"
+++

%s
`, entry.date, entry.date, entry.content)

	for _, attachment := range entry.attachments {
		entryStr += fmt.Sprintf("![](%s)\n", attachment)
	}

	fileName := fmt.Sprintf("%s/%s.md", outputFolder, entry.id)
	os.WriteFile(fileName, []byte(entryStr), 0644)
}

func fetchEntries() ([]byte, error) {
	res, err := http.Get("https://mstdn.social/users/robjperez/outbox?page=true")
	if err != nil {
		return nil, err
	}

	json, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return json, nil
}

func main() {
	// data, err := os.ReadFile("./entries.json")

	// if err != nil {
	// 	panic(err)
	// }
	data, err := fetchEntries()

	var jsonData map[string]json.RawMessage
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		panic(err)
	}

	var items []map[string]json.RawMessage
	json.Unmarshal(jsonData["orderedItems"], &items)

	for idx, item := range items {
		var itemType string
		json.Unmarshal(item["type"], &itemType)
		if itemType != "Create" {
			continue
		}

		entry := BlogEntry{id: strconv.Itoa(idx)}

		var message map[string]json.RawMessage
		json.Unmarshal(item["object"], &message)

		var content string
		json.Unmarshal(message["content"], &content)

		markdown, err := mdConverter.ConvertString(content)
		if err != nil {
			log.Fatal(err)
		}

		entry.content = markdown

		json.Unmarshal(message["published"], &entry.date)

		var attachments []map[string]json.RawMessage
		json.Unmarshal(message["attachment"], &attachments)
		for _, attachment := range attachments {
			var url string
			json.Unmarshal(attachment["url"], &url)
			entry.attachments = append(entry.attachments, url)
		}

		writeEntry(entry, "./blog/content/blog")
	}
}
