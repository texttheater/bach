package main

import (
	"encoding/json"
	"log"
	"os"
)

type object = map[string]any

func modifyBookItem(item object) {
	chapter := item["Chapter"].(object)
	content := chapter["content"].(string)
	chapter["content"] = content // TODO modify
	for _, subItem := range chapter["sub_items"].([]any) {
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
