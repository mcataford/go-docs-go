package generator

import (
	parser "github.com/mcataford/docs/internal/parser"
	"testing"
)

func TestMarkdownGeneratorProducesEmptyMarkdownFromEmptyAst(t *testing.T) {
	fullText := ""
	ast := parser.Parse(fullText, "sample.ts")
	markdown := GenerateMarkdown(ast)

	expectedMarkdown := "# sample.ts"

	if markdown != expectedMarkdown {
		t.Errorf("Unexpected markdown (%q != %q)", markdown, expectedMarkdown)
	}
}

func TestMarkdownGeneratorProducesMarkdownForClass(t *testing.T) {
	fullText := `/* Leading Comment */
    class MyClass {
        /* Method Leading Comment */
        method() {}
    }`
	ast := parser.Parse(fullText, "sample.ts")
	markdown := GenerateMarkdown(ast)

	expectedMarkdown := `# sample.ts
### MyClass
/* Leading Comment */
### method
/* Method Leading Comment */`

	if markdown != expectedMarkdown {
		t.Errorf("Unexpected markdown (%q != %q)", markdown, expectedMarkdown)
	}

}
