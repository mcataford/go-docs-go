package main

import (
	"regexp"
)

type Node struct {
	nodeType NodeType
	raw      string
	start    int
	end      int
	children []Node
}

type NodeType int

const (
	Program NodeType = iota
	FunctionDeclaration
	ClassDeclaration
)

// Parses a text containing Typescript or Javascript code into
// a rudimentary AST. The root node of the returned tree represents
// the full file processed.
func Parse(fullText string) Node {
	currentPosition := 0
	nodes := []Node{}

	root := Node{Program, fullText, 0, len(fullText), []Node{}}

	for currentPosition < (len(fullText) - 1) {
		node, ok := maybeParseFunctionDeclaration(fullText[currentPosition:])

		if ok {
			nodes = append(nodes, node)
			currentPosition = node.end + 1
			continue
		}

		node, ok = maybeParseClassDeclaration(fullText[currentPosition:])

		if ok {
			nodes = append(nodes, node)
			currentPosition = node.end + 1
			continue
		}

		break
	}

	root.children = nodes

	return root
}

// From a text containing source code and a starting point,
// determines the start and end of the next closure.
//
// In this case, we define a closure as a balanced set of
// { } brackets such that we capture the body of blocks like
//
// function myFunction() {
//   // ...Code
// }
func findClosureBoundaries(fullText string, start int) (int, int) {
	stack := []rune{}
	end := -1
	for position, character := range fullText[start:] {
		if character == '{' {
			stack = append(stack, character)
			if end == -1 {
				end = 0
			}
		} else if character == '}' {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 && end != -1 {
			end = position + start
			break
		}
	}

	return start, end
}

// Given the full text of a file or a fragment of a file, builds
// a struct containing the start and end of the function declaration,
// the full text (incl. body) and other metadata about the function block.
//
// See `Node` for more details on how this works.
//
// The function returns a tuple containing the fully-formed node, if
// possible, and a boolean representing whether a node was found or not.
func maybeParseFunctionDeclaration(fullText string) (Node, bool) {
	pattern := `function (?P<functionName>([a-zA-Z_][a-zA-Z0-9_]*))\(.*\)`
	r, _ := regexp.Compile(pattern)

	matches := r.FindStringSubmatchIndex(fullText)

	if len(matches) == 0 {
		return Node{}, false
	}

	start, end := findClosureBoundaries(fullText, matches[0])

	return Node{FunctionDeclaration, fullText[start : end+1], start, end + 1, []Node{}}, true
}

// Given the full text of a file or a fragment of a file, builds
// a struct containing the start and end of the class declaration,
// the full text (incl. body) and other metadata about the class block.
//
// See `Node` for more details on how this works.
//
// The function returns a tuple containing the fully-formed node, if
// possible, and a boolean representing whether a node was found or not.
func maybeParseClassDeclaration(fullText string) (Node, bool) {
	pattern := `class [a-zA-Z_][a-zA-Z0-9_]* ?\{`

	r, _ := regexp.Compile(pattern)

	matches := r.FindStringSubmatchIndex(fullText)

	if len(matches) == 0 {
		return Node{}, false
	}

	start, end := findClosureBoundaries(fullText, matches[0])

	return Node{ClassDeclaration, fullText[start : end+1], start, end + 1, []Node{}}, true
}
