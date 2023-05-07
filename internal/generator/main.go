package generator

import (
	"fmt"
	parser "github.com/mcataford/docs/internal/parser"
	"strings"
)

type ForEachChild func(parser.Node)

var Traverse func(parser.Node, ForEachChild)

// Generates a Markdown document from the provided
// ast.
func GenerateMarkdown(ast parser.Node) string {
	Traverse = func(ast parser.Node, fn ForEachChild) {
		fn(ast)

		for _, child := range ast.Children {
			Traverse(child, fn)
		}
	}

	document := []string{fmt.Sprintf("# %s", ast.Identifier)}

	forEachChild := func(node parser.Node) {
		leadingComment := ""

		if len(node.LeadingComments) > 0 {
			leadingComment = node.LeadingComments[0]
		}

		document = append(document, ([]string{"### " + node.Identifier, leadingComment})...)
	}

	for _, child := range ast.Children {
		Traverse(child, forEachChild)
	}

	return strings.Join(document, "\n")
}
