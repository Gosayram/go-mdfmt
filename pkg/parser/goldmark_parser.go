package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	gmparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// GoldmarkParser implements the Parser interface using goldmark
type GoldmarkParser struct {
	markdown goldmark.Markdown
}

// NewGoldmarkParser creates a new goldmark-based parser
func NewGoldmarkParser() *GoldmarkParser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,           // GitHub Flavored Markdown
			extension.Table,         // Tables support
			extension.Strikethrough, // Strikethrough support
			extension.TaskList,      // Task lists support
		),
		goldmark.WithParserOptions(
			gmparser.WithAutoHeadingID(), // Auto-generate heading IDs
		),
	)

	return &GoldmarkParser{
		markdown: md,
	}
}

// Parse parses the given markdown content and returns an AST
func (p *GoldmarkParser) Parse(content []byte) (*Document, error) {
	// Parse with goldmark
	reader := text.NewReader(content)
	doc := p.markdown.Parser().Parse(reader)

	// Convert goldmark AST to our AST
	ourDoc := &Document{
		Children: make([]Node, 0),
	}

	// Walk through goldmark AST and convert nodes
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		ourNode := p.convertNode(n, content)
		if ourNode != nil {
			ourDoc.Children = append(ourDoc.Children, ourNode)
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to convert AST: %w", err)
	}

	return ourDoc, nil
}

// convertNode converts a goldmark AST node to our AST node
func (p *GoldmarkParser) convertNode(n ast.Node, source []byte) Node {
	switch n.Kind() {
	case ast.KindHeading:
		heading := n.(*ast.Heading)
		return &Heading{
			Level: heading.Level,
			Text:  p.extractText(n, source),
			Style: "atx", // Default to ATX style
		}

	case ast.KindParagraph:
		return &Paragraph{
			Text: p.extractText(n, source),
		}

	case ast.KindList:
		list := n.(*ast.List)
		ourList := &List{
			Ordered: list.IsOrdered(),
			Items:   make([]*ListItem, 0),
			Marker:  p.getListMarker(list),
		}

		// Convert list items
		for child := list.FirstChild(); child != nil; child = child.NextSibling() {
			if child.Kind() == ast.KindListItem {
				item := &ListItem{
					Text:   p.extractText(child, source),
					Marker: p.getListItemMarker(child.(*ast.ListItem)),
				}
				ourList.Items = append(ourList.Items, item)
			}
		}

		return ourList

	case ast.KindFencedCodeBlock, ast.KindCodeBlock:
		code := &CodeBlock{
			Content: p.extractText(n, source),
			Fenced:  n.Kind() == ast.KindFencedCodeBlock,
			Fence:   "```", // Default fence
		}

		// Extract language if it's a fenced code block
		if n.Kind() == ast.KindFencedCodeBlock {
			fenced := n.(*ast.FencedCodeBlock)
			if fenced.Language(source) != nil {
				code.Language = string(fenced.Language(source))
			}
			// Get actual fence character
			if fenced.Info != nil {
				info := string(fenced.Info.Value(source))
				if strings.HasPrefix(info, "~~~") {
					code.Fence = "~~~"
				}
			}
		}

		return code

	case ast.KindText, ast.KindString:
		return &Text{
			Content: p.extractText(n, source),
		}

	default:
		// For other node types, create a generic text node
		text := p.extractText(n, source)
		if text != "" {
			return &Text{
				Content: text,
			}
		}
		return nil
	}
}

// getListMarker determines the list marker from a goldmark list
func (p *GoldmarkParser) getListMarker(list *ast.List) string {
	if list.IsOrdered() {
		return "."
	}
	return "-" // Default bullet
}

// getListItemMarker determines the list item marker
func (p *GoldmarkParser) getListItemMarker(_ *ast.ListItem) string {
	// For now, return a default marker
	// In a real implementation, this would examine the source text
	// to determine the actual marker used
	return "-"
}

// extractText extracts the text content from a goldmark AST node
func (p *GoldmarkParser) extractText(n ast.Node, source []byte) string {
	var buf bytes.Buffer

	// Special handling for different node types
	switch n.Kind() {
	case ast.KindText:
		text := n.(*ast.Text)
		buf.Write(text.Segment.Value(source))
		return buf.String()

	case ast.KindFencedCodeBlock:
		fenced := n.(*ast.FencedCodeBlock)
		for i := 0; i < fenced.Lines().Len(); i++ {
			line := fenced.Lines().At(i)
			buf.Write(line.Value(source))
		}
		return buf.String()

	case ast.KindCodeBlock:
		code := n.(*ast.CodeBlock)
		for i := 0; i < code.Lines().Len(); i++ {
			line := code.Lines().At(i)
			buf.Write(line.Value(source))
		}
		return buf.String()

	case ast.KindString:
		str := n.(*ast.String)
		buf.Write(str.Value)
		return buf.String()
	}

	// For container nodes, extract text from all children
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		childText := p.extractText(child, source)
		buf.WriteString(childText)
	}

	return strings.TrimSpace(buf.String())
}

// Validate checks if the parser is properly configured
func (p *GoldmarkParser) Validate() error {
	if p.markdown == nil {
		return fmt.Errorf("goldmark parser is not initialized")
	}
	return nil
}
