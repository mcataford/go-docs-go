package main

import "testing"

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
