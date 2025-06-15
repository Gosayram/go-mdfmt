package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewGoldmarkParser(t *testing.T) {
	parser := NewGoldmarkParser()
	if parser == nil {
		t.Fatal("NewGoldmarkParser() returned nil")
	}

	err := parser.Validate()
	if err != nil {
		t.Fatalf("Parser validation failed: %v", err)
	}
}

func TestGoldmarkParser_ParseHeading(t *testing.T) {
	parser := NewGoldmarkParser()
	content := []byte("# Hello World\n\nThis is a test.")

	doc, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(doc.Children) == 0 {
		t.Fatal("No children in document")
	}

	// Check if we have at least one heading
	hasHeading := false
	for _, child := range doc.Children {
		if heading, ok := child.(*Heading); ok {
			hasHeading = true
			if heading.Level != 1 {
				t.Errorf("Expected heading level 1, got %d", heading.Level)
			}
			if !strings.Contains(heading.Text, "Hello World") {
				t.Errorf("Expected heading text to contain 'Hello World', got %q", heading.Text)
			}
		}
	}

	if !hasHeading {
		t.Error("No heading found in parsed document")
	}
}

func TestGoldmarkParser_ParseParagraph(t *testing.T) {
	parser := NewGoldmarkParser()
	content := []byte("This is a simple paragraph.")

	doc, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(doc.Children) == 0 {
		t.Fatal("No children in document")
	}

	// Check if we have at least one paragraph
	hasParagraph := false
	for _, child := range doc.Children {
		if paragraph, ok := child.(*Paragraph); ok {
			hasParagraph = true
			if !strings.Contains(paragraph.Text, "simple paragraph") {
				t.Errorf("Expected paragraph text to contain 'simple paragraph', got %q", paragraph.Text)
			}
		}
	}

	if !hasParagraph {
		t.Error("No paragraph found in parsed document")
	}
}

func TestGoldmarkParser_ParseList(t *testing.T) {
	parser := NewGoldmarkParser()
	content := []byte(`
- Item 1
- Item 2
- Item 3
`)

	doc, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find the list
	var list *List
	for _, child := range doc.Children {
		if l, ok := child.(*List); ok {
			list = l
			break
		}
	}

	if list == nil {
		t.Fatal("No list found in parsed document")
	}

	if list.Ordered {
		t.Error("Expected unordered list, got ordered")
	}

	if len(list.Items) != 3 {
		t.Errorf("Expected 3 list items, got %d", len(list.Items))
	}

	expectedItems := []string{"Item 1", "Item 2", "Item 3"}
	for i, item := range list.Items {
		if i < len(expectedItems) {
			if !strings.Contains(item.Text, expectedItems[i]) {
				t.Errorf("Expected item %d to contain %q, got %q", i, expectedItems[i], item.Text)
			}
		}
	}
}

func TestGoldmarkParser_ParseOrderedList(t *testing.T) {
	parser := NewGoldmarkParser()
	content := []byte(`
1. First item
2. Second item
3. Third item
`)

	doc, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find the list
	var list *List
	for _, child := range doc.Children {
		if l, ok := child.(*List); ok {
			list = l
			break
		}
	}

	if list == nil {
		t.Fatal("No list found in parsed document")
	}

	if !list.Ordered {
		t.Error("Expected ordered list, got unordered")
	}

	if len(list.Items) != 3 {
		t.Errorf("Expected 3 list items, got %d", len(list.Items))
	}
}

func TestGoldmarkParser_ParseCodeBlock(t *testing.T) {
	parser := NewGoldmarkParser()
	content := []byte("```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```")

	doc, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find the code block
	var codeBlock *CodeBlock
	for _, child := range doc.Children {
		if cb, ok := child.(*CodeBlock); ok {
			codeBlock = cb
			break
		}
	}

	if codeBlock == nil {
		t.Fatal("No code block found in parsed document")
	}

	if codeBlock.Language != "go" {
		t.Errorf("Expected language 'go', got %q", codeBlock.Language)
	}

	if !codeBlock.Fenced {
		t.Error("Expected fenced code block")
	}

	if !strings.Contains(codeBlock.Content, "func main") {
		t.Errorf("Expected code content to contain 'func main', got %q", codeBlock.Content)
	}
}

func TestGoldmarkParser_ParseComplexDocument(t *testing.T) {
	parser := NewGoldmarkParser()
	content := []byte(`# Title

This is a paragraph with **bold** and *italic* text.

## Subtitle

Here's a list:
- Item 1
- Item 2

And a code block:
` + "```python\nprint('Hello, World!')\n```")

	doc, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(doc.Children) == 0 {
		t.Fatal("No children in document")
	}

	// Count different node types
	headingCount := 0
	paragraphCount := 0
	listCount := 0
	codeBlockCount := 0

	for _, child := range doc.Children {
		switch child.(type) {
		case *Heading:
			headingCount++
		case *Paragraph:
			paragraphCount++
		case *List:
			listCount++
		case *CodeBlock:
			codeBlockCount++
		}
	}

	if headingCount < 1 {
		t.Errorf("Expected at least 1 heading, got %d", headingCount)
	}
	if paragraphCount < 1 {
		t.Errorf("Expected at least 1 paragraph, got %d", paragraphCount)
	}
	if listCount < 1 {
		t.Errorf("Expected at least 1 list, got %d", listCount)
	}
	if codeBlockCount < 1 {
		t.Errorf("Expected at least 1 code block, got %d", codeBlockCount)
	}
}

func TestGoldmarkParser_EmptyDocument(t *testing.T) {
	parser := NewGoldmarkParser()
	content := []byte("")

	doc, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected non-nil document")
	}

	// Empty document should have an empty children slice
	if doc.Children == nil {
		t.Error("Expected non-nil children slice")
	}
}

func TestGoldmarkParser_Validate(t *testing.T) {
	parser := NewGoldmarkParser()

	err := parser.Validate()
	if err != nil {
		t.Errorf("Expected valid parser, got error: %v", err)
	}

	// Test with nil markdown (should not happen in normal usage)
	invalidParser := &GoldmarkParser{markdown: nil}
	err = invalidParser.Validate()
	if err == nil {
		t.Error("Expected validation error for invalid parser")
	}
}

// Benchmark tests
func BenchmarkGoldmarkParser_ParseSimpleDocument(b *testing.B) {
	parser := NewGoldmarkParser()
	content := `# Heading

This is a paragraph with some text.

- List item 1
- List item 2
- List item 3

` + "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse([]byte(content))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoldmarkParser_ParseComplexDocument(b *testing.B) {
	parser := NewGoldmarkParser()
	content := `# Main Heading

## Section 1

This is a paragraph with **bold** and *italic* text, plus some ` + "`inline code`" + `.

### Subsection

1. First ordered item
2. Second ordered item
   - Nested unordered item
   - Another nested item
3. Third ordered item

## Section 2

Here's a code block:

` + "```javascript\nfunction hello() {\n    console.log('Hello, world!');\n    return true;\n}\n```" + `

And here's a [link](https://example.com) and some more text.

> This is a blockquote
> with multiple lines

## Final Section

- [ ] Todo item 1
- [x] Completed item
- [ ] Todo item 2

Final paragraph with some text.
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse([]byte(content))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoldmarkParser_ParseLargeDocument(b *testing.B) {
	parser := NewGoldmarkParser()

	// Generate a large document
	var content strings.Builder
	content.WriteString("# Large Document\n\n")

	for i := 0; i < 100; i++ {
		content.WriteString(fmt.Sprintf("## Section %d\n\n", i+1))
		content.WriteString("This is a paragraph with some text that describes the section.\n\n")

		content.WriteString("### Subsection\n\n")
		for j := 0; j < 10; j++ {
			content.WriteString(fmt.Sprintf("- List item %d with some descriptive text\n", j+1))
		}
		content.WriteString("\n")

		if i%10 == 0 {
			content.WriteString("```go\n")
			content.WriteString("func example() {\n")
			content.WriteString("    fmt.Println(\"Example code\")\n")
			content.WriteString("}\n")
			content.WriteString("```\n\n")
		}
	}

	contentBytes := []byte(content.String())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(contentBytes)
		if err != nil {
			b.Fatal(err)
		}
	}
}
