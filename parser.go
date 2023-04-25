package main

import (
	"regexp"
)

type Node struct {
	nodeType        NodeType
	raw             string
	identifier      string
	start           int
	end             int
	children        []Node
	leadingComments []string
}

type NodeType int

const (
	Program NodeType = iota
	FunctionDeclaration
	ClassDeclaration
	ClassMethod
)

// Raw patterns used when searching for syntactic elements in source.
var LeadingCommentPattern = `(\/\*(.|\s)*\*\/)*\s*`
var FunctionDeclarationCommentPattern = LeadingCommentPattern + `(?P<functionDeclaration>(function (?P<functionName>([a-zA-Z_][a-zA-Z0-9_]*))\(.*\)))`
var ClassDeclarationPattern = LeadingCommentPattern + `(?P<classDeclaration>(class (?P<className>([a-zA-Z_][a-zA-Z0-9_]*))))`
var TypedArgumentPattern = `[a-zA-Z_][a-zA-Z0-9_]*\s*(:\s*[a-zA-Z_][a-zA-Z0-9_])?`

// Compiled regexp patterns.
var rFunctionDeclaration = regexp.MustCompile(FunctionDeclarationCommentPattern)
var rClassDeclaration = regexp.MustCompile(ClassDeclarationPattern)
var rClassMethod = regexp.MustCompile(LeadingCommentPattern + `(?P<classMethod>([a-zA-Z_][a-zA-Z0-9_]*))\((` + TypedArgumentPattern + `\s*)*\)`)

// Parses a text containing Typescript or Javascript code into
// a rudimentary AST. The root node of the returned tree represents
// the full file processed.
func Parse(fullText string) Node {
	currentPosition := 0
	nodes := []Node{}

	root := Node{Program, fullText, "", 0, len(fullText), []Node{}, []string{}}

	for currentPosition < (len(fullText) - 1) {
		node, ok := maybeParseFunctionDeclaration(fullText[currentPosition:], currentPosition)

		if ok {
			nodes = append(nodes, node)
			currentPosition = node.end + 1
			continue
		}

		node, ok = maybeParseClassDeclaration(fullText[currentPosition:], currentPosition)

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
func findClosureBoundaries(fullText string, scanStart int) (int, int) {
	depth, start, end := 0, -1, -1

	for position, character := range fullText[scanStart:] {
		if character == '{' {
			depth = depth + 1
			if start == -1 {
				start = position + scanStart
			}
			if end == -1 {
				end = 0
			}
		} else if character == '}' {
			depth = depth - 1
		}
		if depth == 0 && end != -1 {
			end = position + scanStart
			break
		}
	}

	return start, end
}

// Extracts any block comments of the form /* ... */ from
// the given text.
//
// This always returns an array of strings representing
// the comments found, and always scans the full provided text.
func findBlockComments(fullText string) []string {
	comments := []string{}

	inComment := false
	previousRune := ' '
	end := -1

	for position, character := range fullText {
		if character == '*' && previousRune == '/' {
			inComment = true
			if end == -1 {
				end = position - 1
			}
		} else if character == '/' && previousRune == '*' && inComment {
			inComment = false
			comments = append(comments, fullText[end:position+1])
			end = position + 1
		}

		previousRune = character
	}

	return comments
}

// Given the full text of a file or a fragment of a file, builds
// a struct containing the start and end of the function declaration,
// the full text (incl. body) and other metadata about the function block.
//
// See `Node` for more details on how this works.
//
// The function returns a tuple containing the fully-formed node, if
// possible, and a boolean representing whether a node was found or not.
func maybeParseFunctionDeclaration(fullText string, offset int) (Node, bool) {
	matches := rFunctionDeclaration.FindStringSubmatchIndex(fullText)
	m := rFunctionDeclaration.FindStringSubmatch(fullText)

	if len(matches) == 0 {
		return Node{}, false
	}

	start, end := findClosureBoundaries(fullText, matches[0])

	fnStart := rFunctionDeclaration.SubexpIndex("functionDeclaration")
	fnName := m[rFunctionDeclaration.SubexpIndex("functionName")]
	leadingComments := findBlockComments(fullText[:matches[fnStart]+offset])
	return Node{FunctionDeclaration, fullText[start : end+1], fnName, start + offset, end + 1 + offset, []Node{}, leadingComments}, true
}

// Given the full text of a file or a fragment of a file, builds
// a struct containing the start and end of the class declaration,
// the full text (incl. body) and other metadata about the class block.
//
// See `Node` for more details on how this works.
//
// The function returns a tuple containing the fully-formed node, if
// possible, and a boolean representing whether a node was found or not.
func maybeParseClassDeclaration(fullText string, offset int) (Node, bool) {
	matches := rClassDeclaration.FindStringSubmatchIndex(fullText)
	m := rClassDeclaration.FindStringSubmatch(fullText)

	if len(matches) == 0 {
		return Node{}, false
	}

	declarationStart := matches[0]
	start, end := findClosureBoundaries(fullText, declarationStart)

	clsStart := rClassDeclaration.SubexpIndex("classDeclaration")
	clsName := m[rClassDeclaration.SubexpIndex("className")]
	leadingComments := findBlockComments(fullText[:matches[clsStart]+offset])

	children := parseClassBody(fullText[start:end+1], offset)

	return Node{ClassDeclaration, fullText[declarationStart : end+1], clsName, declarationStart + offset, end + 1 + offset, children, leadingComments}, true
}

// Parses a class body to extract children methods.
//
// The children method are returned as an array of Node structs.
func parseClassBody(fullText string, offset int) []Node {
	children := []Node{}

	indexMatches := rClassMethod.FindAllStringIndex(fullText, -1)
	fullMatches := rClassMethod.FindAllStringSubmatch(fullText, -1)

	classMethodNameIndex := rClassMethod.SubexpIndex("classMethod")

	lastBlockEnd := 0
	for position, matchIndices := range indexMatches {
		declarationStart := matchIndices[0]

		methodName := fullMatches[position][classMethodNameIndex]

		// Depending on syntax, another block might have been
		// identified as "method-like" within the last block.
		// In this case, let's skip ahead since method cannot
		// contain other methods.
		if declarationStart < lastBlockEnd {
			continue
		}

		_, end := findClosureBoundaries(fullText, declarationStart)
		lastBlockEnd = end
		leadingComments := findBlockComments(fullText[declarationStart:indexMatches[position][1]])
		children = append(children, Node{ClassMethod, fullText[declarationStart : end+1], methodName, declarationStart, end + 1, []Node{}, leadingComments})
	}

	return children
}
