package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/texttheater/bach/docutil"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/shapes"
)

type object = map[string]any

func processExampleSet(scanner *bufio.Scanner, buffer *bytes.Buffer) error {
	currentExample := shapes.Example{}
	examples := []shapes.Example{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "```" {
			if len(currentExample.Program) > 0 {
				examples = append(examples, currentExample)
			}
			docutil.PrintExamplesTable(buffer, examples)
			for _, x := range examples {
				interpreter.TestExample(x)
			}
			return nil
		} else if line == "" {
			examples = append(examples, currentExample)
			currentExample = shapes.Example{}
		} else if strings.HasPrefix(line, "P ") {
			currentExample.Program = line[2:]
		} else if strings.HasPrefix(line, "T ") {
			currentExample.OutputType = line[2:]
		} else if strings.HasPrefix(line, "V ") {
			currentExample.OutputValue = line[2:]
		} else if strings.HasPrefix(line, "E ") {
			e, err := docutil.ParseError(line[2:])
			if err != nil {
				return err
			}
			currentExample.Error = e
		}
	}
	return nil
}

func processExamples(content string) (string, error) {
	reader := strings.NewReader(content)
	scanner := bufio.NewScanner(reader)
	var buffer bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()
		if line == "```bachdoc" {
			err := processExampleSet(scanner, &buffer)
			if err != nil {
				return "", err
			}
		} else {
			buffer.WriteString(line)
		}
		buffer.WriteString("\n")
	}
	return buffer.String(), nil
}

func modifyBookItem(item object) error {
	chapter := item["Chapter"].(object)
	content := chapter["content"].(string)
	var err error
	chapter["content"], err = processExamples(content)
	if err != nil {
		return err
	}
	subItems := chapter["sub_items"].([]any)
	for _, subItem := range subItems {
		err := modifyBookItem(subItem.(object))
		if err != nil {
			return err
		}
	}
	return nil
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
		err := modifyBookItem(section.(object))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = json.NewEncoder(os.Stdout).Encode(book)
	if err != nil {
		log.Fatal(err)
	}
}
