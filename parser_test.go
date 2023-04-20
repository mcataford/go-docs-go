package main

import (
	"golang.org/x/exp/slices"
	"strings"
	"testing"
)

func TestParseEmptyFile(t *testing.T) {
	fullText := ""
	ast := Parse(fullText)

	children_count := len(ast.children)

	if children_count != 0 {
		t.Errorf("Expected no children when parsing empty file, got %q", children_count)
	}

	if ast.nodeType != Program {
		t.Errorf("Expected root node to have type Program, got %q", ast.nodeType)
	}
}

func TestParseClassLeadingComments(t *testing.T) {
	leadingComments := []string{`/*
    * This is a leading comment.
    */`,
		`/* This is another leading comments. */`}
	fullText := strings.Join(leadingComments, "") + `
    class MyClass {
        // Logic.
    }`

	ast := Parse(fullText)
	classDeclarationNode := ast.children[0]

	if classDeclarationNode.nodeType != ClassDeclaration {
		t.Errorf("Expected a class node, got %q", classDeclarationNode.nodeType)
	}

	leadingCommentsFound := classDeclarationNode.leadingComments

	if len(leadingCommentsFound) != 2 {
		t.Errorf("Expected two leading comments, got %q", len(leadingCommentsFound))
	}

	if !slices.Equal(classDeclarationNode.leadingComments, leadingComments) {
		t.Errorf("Didn't find the leading comments expected: %q != %q", classDeclarationNode.leadingComments, leadingComments)
	}

}
