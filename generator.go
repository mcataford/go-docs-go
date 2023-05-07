package main

import (
	"fmt"
	"strings"
)

type ForEachChild func(Node)

var Traverse func(Node, ForEachChild)

// Generates a Markdown document from the provided
// ast.
func GenerateMarkdown(ast Node) string {
	Traverse = func(ast Node, fn ForEachChild) {
		fn(ast)

		for _, child := range ast.children {
			Traverse(child, fn)
		}
	}

	document := []string{fmt.Sprintf("# %s", ast.identifier)}

	forEachChild := func(node Node) {
		leadingComment := ""

		if len(node.leadingComments) > 0 {
			leadingComment = node.leadingComments[0]
		}

		document = append(document, ([]string{"### " + node.identifier, leadingComment})...)
	}

	for _, child := range ast.children {
		Traverse(child, forEachChild)
	}

	return strings.Join(document, "\n")
}
