package main

import (
	"io/ioutil"
	"os"
)

func main() {
	target := os.Args[1]
	readBytes, _ := ioutil.ReadFile(target)
	fileContent := string(readBytes)

	source := Parse(fileContent)

	generated_markup := GenerateMarkdown(source)

	file, err := os.Create("API.md")
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(generated_markup)
}
