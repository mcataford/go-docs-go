package generator

import (
	"fmt"
	parser "github.com/mcataford/docs/internal/parser"
	"strings"
)

type ForEachChild func(parser.Node, int)

var Traverse func(parser.Node, ForEachChild, int)

// Generates a Markdown document from the provided
// ast.
func GenerateMarkdown(ast parser.Node) string {
	Traverse = func(ast parser.Node, fn ForEachChild, depth int) {
		fn(ast, depth)

		for _, child := range ast.Children {
			Traverse(child, fn, depth+1)
		}
	}

	document := []string{fmt.Sprintf("# %s", ast.Identifier)}

	forEachChild := func(node parser.Node, depth int) {
		leadingComment := ""

		if len(node.LeadingComments) > 0 {
			leadingComment = node.LeadingComments[0]
		}

		var headingPrefix string

		if depth < 3 {
			headingPrefix = strings.Repeat("#", 1+depth)
		} else {
			headingPrefix = "###"
		}

		document = append(document, ([]string{fmt.Sprintf("%s %s", headingPrefix, node.Identifier), leadingComment})...)
	}

	for _, child := range ast.Children {
		Traverse(child, forEachChild, 1)
	}

	return strings.Join(document, "\n")
}
