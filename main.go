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

	Generate(source)
}
