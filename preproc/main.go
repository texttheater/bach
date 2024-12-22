package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strings"
)

type object = map[string]any

func processExampleSet(scanner *bufio.Scanner, buffer *bytes.Buffer) {
	buffer.WriteString("(example set goes here)\n")
	for scanner.Scan() {
		line := scanner.Text()
		if line == "```" {
			return
		}
	}
}

func processExamples(content string) string {
	reader := strings.NewReader(content)
	scanner := bufio.NewScanner(reader)
	var buffer bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()
		if line == "```bachdoc" {
			processExampleSet(scanner, &buffer)
		} else {
			buffer.WriteString(line)
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func modifyBookItem(item object) {
	chapter := item["Chapter"].(object)
	content := chapter["content"].(string)
	chapter["content"] = processExamples(content)
	subItems := chapter["sub_items"].([]any)
	for _, subItem := range subItems {
		modifyBookItem(subItem.(object))
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "supports" {
		os.Exit(0)
	}
	decoder := json.NewDecoder(os.Stdin)
	var data []any
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	book := data[1].(object)
	sections := book["sections"].([]any)
	for _, section := range sections {
		modifyBookItem(section.(object))
	}
	err = json.NewEncoder(os.Stdout).Encode(book)
	if err != nil {
		log.Fatal(err)
	}
}
