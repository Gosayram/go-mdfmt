package renderer

import (
	"io"
	"strings"

	"github.com/Gosayram/go-mdfmt/pkg/config"
	"github.com/Gosayram/go-mdfmt/pkg/parser"
)

// Renderer represents a renderer that converts AST back to markdown
type Renderer interface {
	// Render renders the AST to markdown
	Render(doc *parser.Document, cfg *config.Config) (string, error)
	// RenderTo renders the AST to a writer
	RenderTo(w io.Writer, doc *parser.Document, cfg *config.Config) error
}

// MarkdownRenderer renders AST back to markdown format
type MarkdownRenderer struct {
	output strings.Builder
	config *config.Config
}

// New creates a new markdown renderer
func New() *MarkdownRenderer {
	return &MarkdownRenderer{}
}

// Render renders the AST to markdown string
func (r *MarkdownRenderer) Render(doc *parser.Document, cfg *config.Config) (string, error) {
	r.output.Reset()
	r.config = cfg

	if err := r.renderDocument(doc, 0); err != nil {
		return "", err
	}

	result := r.output.String()

	// Apply final whitespace normalization
	if cfg.Whitespace.EnsureFinalNewline && !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	return result, nil
}

// RenderTo renders the AST to a writer
func (r *MarkdownRenderer) RenderTo(w io.Writer, doc *parser.Document, cfg *config.Config) error {
	content, err := r.Render(doc, cfg)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(content))
	return err
}

// renderDocument renders a document node
func (r *MarkdownRenderer) renderDocument(doc *parser.Document, depth int) error {
	for _, child := range doc.Children {
		if err := r.renderNode(child, depth); err != nil {
			return err
		}
	}
	return nil
}

// renderNode renders a single node
func (r *MarkdownRenderer) renderNode(node parser.Node, depth int) error {
	switch n := node.(type) {
	case *parser.Heading:
		return r.renderHeading(n, depth)
	case *parser.Paragraph:
		return r.renderParagraph(n, depth)
	case *parser.List:
		return r.renderList(n, depth)
	case *parser.ListItem:
		return r.renderListItem(n, depth)
	case *parser.CodeBlock:
		return r.renderCodeBlock(n, depth)
	case *parser.Text:
		return r.renderText(n, depth)
	default:
		// Unknown node type, skip
		return nil
	}
}

// renderHeading renders a heading node
func (r *MarkdownRenderer) renderHeading(heading *parser.Heading, depth int) error {
	if heading.Style == "setext" && heading.Level <= 2 {
		// Setext-style heading
		r.output.WriteString(heading.Text)
		r.output.WriteString("\n")

		marker := "="
		if heading.Level == 2 {
			marker = "-"
		}

		textLength := len(strings.TrimSpace(heading.Text))
		if textLength == 0 {
			textLength = 3 // minimum length
		}

		r.output.WriteString(strings.Repeat(marker, textLength))
		r.output.WriteString("\n\n")
	} else {
		// ATX-style heading
		r.output.WriteString(strings.Repeat("#", heading.Level))
		r.output.WriteString(" ")
		r.output.WriteString(heading.Text)
		r.output.WriteString("\n\n")
	}

	return nil
}

// renderParagraph renders a paragraph node
func (r *MarkdownRenderer) renderParagraph(para *parser.Paragraph, depth int) error {
	content := para.Text

	// Apply line width wrapping
	if r.config.LineWidth > 0 {
		content = r.wrapText(content, r.config.LineWidth)
	}

	r.output.WriteString(content)
	r.output.WriteString("\n\n")

	return nil
}

// renderList renders a list node
func (r *MarkdownRenderer) renderList(list *parser.List, depth int) error {
	for i, item := range list.Items {
		if i > 0 {
			r.output.WriteString("\n")
		}
		if err := r.renderListItem(item, depth+1); err != nil {
			return err
		}
	}

	r.output.WriteString("\n")
	return nil
}

// renderListItem renders a list item node
func (r *MarkdownRenderer) renderListItem(item *parser.ListItem, depth int) error {
	indent := strings.Repeat("  ", depth)

	// Determine marker
	marker := item.Marker
	if marker == "" {
		marker = r.config.List.BulletStyle
	}

	r.output.WriteString(indent)
	r.output.WriteString(marker)
	r.output.WriteString(" ")
	r.output.WriteString(item.Text)

	return nil
}

// renderCodeBlock renders a code block node
func (r *MarkdownRenderer) renderCodeBlock(code *parser.CodeBlock, depth int) error {
	if code.Fenced {
		r.output.WriteString(code.Fence)
		if code.Language != "" {
			r.output.WriteString(code.Language)
		}
		r.output.WriteString("\n")
		r.output.WriteString(code.Content)
		if !strings.HasSuffix(code.Content, "\n") {
			r.output.WriteString("\n")
		}
		r.output.WriteString(code.Fence)
		r.output.WriteString("\n\n")
	} else {
		// Indented code block
		lines := strings.Split(code.Content, "\n")
		for _, line := range lines {
			r.output.WriteString("    ")
			r.output.WriteString(line)
			r.output.WriteString("\n")
		}
		r.output.WriteString("\n")
	}

	return nil
}

// renderText renders a text node
func (r *MarkdownRenderer) renderText(text *parser.Text, depth int) error {
	content := text.Content

	// Apply whitespace normalization
	if r.config.Whitespace.TrimTrailingSpaces {
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			lines[i] = strings.TrimRight(line, " \t")
		}
		content = strings.Join(lines, "\n")
	}

	r.output.WriteString(content)
	return nil
}

// wrapText wraps text to the specified line width
func (r *MarkdownRenderer) wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine strings.Builder

	for i, word := range words {
		// Check if adding this word would exceed the line width
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
			// Start a new line
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)

		// If this is the last word, add the current line
		if i == len(words)-1 {
			lines = append(lines, currentLine.String())
		}
	}

	return strings.Join(lines, "\n")
}
