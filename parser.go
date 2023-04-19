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

// Given the full text of a source file and a starting position
// marking where a function declaration's signature starts, this
// scans the characters that follow until it can find the full
// closure of the function block.
//
// The full function text is returned.
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
