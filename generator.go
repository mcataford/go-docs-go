package main

import (
	"fmt"
	"log"
	"strings"
)

type ForEachChild func(Node)

var Traverse func(Node, ForEachChild)

func GenerateMarkdown(ast Node) string {
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

	log.Println(fmt.Sprintf("%+v", ast))

	Traverse(ast, forEachChild)

	return strings.Join(document, "\n")
}
