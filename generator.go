package main

import (
	"os"
	"strings"
)

type ForEachChild func(Node)

var Traverse func(Node, ForEachChild)

func Generate(ast Node) {
	Traverse = func(ast Node, fn ForEachChild) {
		fn(ast)

		for _, child := range ast.children {
			Traverse(child, fn)
		}
	}

	document := []string{"# API Documentation"}

	forEachChild := func(node Node) {
		leadingComment := ""

		if len(node.leadingComments) > 0 {
			leadingComment = node.leadingComments[0]
		}

		document = append(document, ([]string{"### " + node.identifier, leadingComment})...)
	}

	Traverse(ast, forEachChild)

	file, err := os.Create("API.md")
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(strings.Join(document, "\n"))
}
