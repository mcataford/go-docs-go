package parser

import (
	"golang.org/x/exp/slices"
	"strings"
	"testing"
)

func TestParseEmptyFile(t *testing.T) {
	fullText := ""
	ast := Parse(fullText, "")

	children_count := len(ast.Children)

	if children_count != 0 {
		t.Errorf("Expected no children when parsing empty file, got %q", children_count)
	}

	if ast.NodeType != Program {
		t.Errorf("Expected root node to have type Program, got %q", ast.NodeType)
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

	ast := Parse(fullText, "")
	classDeclarationNode := ast.Children[0]

	if classDeclarationNode.NodeType != ClassDeclaration {
		t.Errorf("Expected a class node, got %q", classDeclarationNode.NodeType)
	}

	leadingCommentsFound := classDeclarationNode.LeadingComments

	if len(leadingCommentsFound) != 2 {
		t.Errorf("Expected two leading comments, got %q", len(leadingCommentsFound))
	}

	if !slices.Equal(classDeclarationNode.LeadingComments, leadingComments) {
		t.Errorf("Didn't find the leading comments expected: %q != %q", classDeclarationNode.LeadingComments, leadingComments)
	}

}

func TestParseClassMethodLeadingComments(t *testing.T) {
	fullText := `class MyClass {
        myMethod() {
            // wo.
        }

        /*
        * Comment
        */
        myCommentedMethod() {
            // wo.
        }
    }`

	ast := Parse(fullText, "")
	classDeclarationNode := ast.Children[0]

	uncommentedMethod := classDeclarationNode.Children[0]

	if len(uncommentedMethod.LeadingComments) != 0 {
		t.Errorf("Unexpected leading comment on method that does not have one.")
	}

	commentedMethod := classDeclarationNode.Children[1]

	if len(commentedMethod.LeadingComments) != 1 {
		t.Errorf("Expected comment on method but found none.")
	}

	leadingComment := commentedMethod.LeadingComments[0]

	expectedComment := `/*
        * Comment
        */`

	if leadingComment != expectedComment {
		t.Errorf("Comment does not match expectation (%q != %q)", leadingComment, expectedComment)
	}
}
