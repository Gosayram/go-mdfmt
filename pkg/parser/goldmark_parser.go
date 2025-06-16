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

const (
	// StrongEmphasisLevel defines the level for strong emphasis (**)
	StrongEmphasisLevel = 2
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

	// Walk through goldmark AST and convert only top-level nodes
	for child := doc.FirstChild(); child != nil; child = child.NextSibling() {
		ourNode := p.convertNode(child, content)
		if ourNode != nil {
			ourDoc.Children = append(ourDoc.Children, ourNode)
		}
	}

	return ourDoc, nil
}

// convertNode converts a goldmark AST node to our AST node
func (p *GoldmarkParser) convertNode(n ast.Node, source []byte) Node {
	switch n.Kind() {
	case ast.KindHeading:
		return p.convertHeading(n, source)
	case ast.KindParagraph:
		return p.convertParagraph(n, source)
	case ast.KindList:
		return p.convertList(n, source)
	case ast.KindFencedCodeBlock, ast.KindCodeBlock:
		return p.convertCodeBlock(n, source)
	case ast.KindText, ast.KindString:
		return p.convertText(n, source)
	default:
		return p.convertGenericNode(n, source)
	}
}

// convertHeading converts a heading node
func (p *GoldmarkParser) convertHeading(n ast.Node, source []byte) Node {
	heading := n.(*ast.Heading)
	headingText := p.extractText(n, source)
	headingText = strings.Join(strings.Fields(headingText), " ")
	return &Heading{
		Level: heading.Level,
		Text:  strings.TrimSpace(headingText),
		Style: "atx",
	}
}

// convertParagraph converts a paragraph node
func (p *GoldmarkParser) convertParagraph(n ast.Node, source []byte) Node {
	return &Paragraph{
		Text: p.extractText(n, source),
	}
}

// convertList converts a list node
func (p *GoldmarkParser) convertList(n ast.Node, source []byte) Node {
	list := n.(*ast.List)
	ourList := &List{
		Ordered: list.IsOrdered(),
		Items:   make([]*ListItem, 0),
		Marker:  p.getListMarker(list),
	}

	for child := list.FirstChild(); child != nil; child = child.NextSibling() {
		if child.Kind() == ast.KindListItem {
			item := p.convertListItem(child, source)
			ourList.Items = append(ourList.Items, item)
		}
	}
	return ourList
}

// convertListItem converts a list item node
func (p *GoldmarkParser) convertListItem(n ast.Node, source []byte) *ListItem {
	item := &ListItem{
		Text:     p.extractText(n, source),
		Marker:   p.getListItemMarker(n.(*ast.ListItem)),
		Children: make([]Node, 0),
	}

	for nestedChild := n.FirstChild(); nestedChild != nil; nestedChild = nestedChild.NextSibling() {
		if nestedChild.Kind() == ast.KindList {
			nestedList := p.convertNode(nestedChild, source)
			if nestedList != nil {
				item.Children = append(item.Children, nestedList)
			}
		}
	}
	return item
}

// convertCodeBlock converts a code block node
func (p *GoldmarkParser) convertCodeBlock(n ast.Node, source []byte) Node {
	code := &CodeBlock{
		Content: p.extractText(n, source),
		Fenced:  n.Kind() == ast.KindFencedCodeBlock,
		Fence:   "```",
	}

	if n.Kind() == ast.KindFencedCodeBlock {
		p.extractCodeBlockInfo(n, source, code)
	}
	return code
}

// extractCodeBlockInfo extracts language and fence info from fenced code block
func (p *GoldmarkParser) extractCodeBlockInfo(n ast.Node, source []byte, code *CodeBlock) {
	fenced := n.(*ast.FencedCodeBlock)
	if fenced.Language(source) != nil {
		code.Language = string(fenced.Language(source))
	}
	if fenced.Info != nil {
		info := string(fenced.Info.Value(source))
		if strings.HasPrefix(info, "~~~") {
			code.Fence = "~~~"
		}
	}
}

// convertText converts a text/string node
func (p *GoldmarkParser) convertText(n ast.Node, source []byte) Node {
	return &Text{
		Content: p.extractText(n, source),
	}
}

// convertGenericNode converts other node types to text
func (p *GoldmarkParser) convertGenericNode(n ast.Node, source []byte) Node {
	content := p.extractText(n, source)
	if content != "" {
		return &Text{
			Content: content,
		}
	}
	return nil
}

// getListMarker determines the list marker from a goldmark list
func (p *GoldmarkParser) getListMarker(list *ast.List) string {
	if list.IsOrdered() {
		return "."
	}
	return "-" // Default bullet
}

// getListItemMarker determines the list item marker
func (p *GoldmarkParser) getListItemMarker(item *ast.ListItem) string {
	// Check if this is part of an ordered list
	if parent := item.Parent(); parent != nil && parent.Kind() == ast.KindList {
		list := parent.(*ast.List)
		if list.IsOrdered() {
			// For ordered lists, we'll let the formatter handle the numbering
			return "1."
		}
	}
	return "-" // Default bullet for unordered lists
}

// extractText extracts the text content from a goldmark AST node
func (p *GoldmarkParser) extractText(n ast.Node, source []byte) string {
	switch n.Kind() {
	case ast.KindText, ast.KindString:
		return p.extractSimpleText(n, source)
	case ast.KindFencedCodeBlock, ast.KindCodeBlock:
		return p.extractCodeBlockText(n, source)
	case ast.KindListItem:
		return p.extractListItemText(n, source)
	case ast.KindList:
		return ""
	case ast.KindParagraph:
		return p.extractParagraphText(n, source)
	default:
		return p.extractGenericText(n, source)
	}
}

// extractSimpleText extracts text from simple text/string nodes
func (p *GoldmarkParser) extractSimpleText(n ast.Node, source []byte) string {
	switch n.Kind() {
	case ast.KindText:
		textNode := n.(*ast.Text)
		return string(textNode.Segment.Value(source))
	case ast.KindString:
		str := n.(*ast.String)
		return string(str.Value)
	}
	return ""
}

// extractCodeBlockText extracts text from code block nodes
func (p *GoldmarkParser) extractCodeBlockText(n ast.Node, source []byte) string {
	var buf bytes.Buffer

	switch n.Kind() {
	case ast.KindFencedCodeBlock:
		fenced := n.(*ast.FencedCodeBlock)
		for i := 0; i < fenced.Lines().Len(); i++ {
			line := fenced.Lines().At(i)
			buf.Write(line.Value(source))
		}
	case ast.KindCodeBlock:
		code := n.(*ast.CodeBlock)
		for i := 0; i < code.Lines().Len(); i++ {
			line := code.Lines().At(i)
			buf.Write(line.Value(source))
		}
	}
	return buf.String()
}

// extractListItemText extracts text from list item nodes
func (p *GoldmarkParser) extractListItemText(n ast.Node, source []byte) string {
	var buf bytes.Buffer

	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if child.Kind() != ast.KindList {
			var childText string
			if child.Kind() == ast.KindParagraph {
				// Use paragraph text extraction to preserve inline formatting
				childText = p.extractParagraphText(child, source)
			} else {
				// For all other nodes, try to extract with inline formatting
				childText = p.extractWithInlineFormatting(child, source)
			}
			if childText != "" {
				if buf.Len() > 0 {
					buf.WriteString(" ")
				}
				buf.WriteString(childText)
			}
		}
	}
	return strings.TrimSpace(buf.String())
}

// extractWithInlineFormatting extracts text preserving inline formatting
func (p *GoldmarkParser) extractWithInlineFormatting(n ast.Node, source []byte) string {
	var buf bytes.Buffer

	switch n.Kind() {
	case ast.KindText:
		textNode := n.(*ast.Text)
		buf.Write(textNode.Segment.Value(source))
	case ast.KindEmphasis:
		p.extractEmphasisText(n, source, &buf)
	case ast.KindCodeSpan:
		p.extractCodeSpanText(n, source, &buf)
	case ast.KindLink:
		p.extractLinkText(n, source, &buf)
	default:
		// For container nodes, process children with inline formatting
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			childText := p.extractWithInlineFormatting(child, source)
			buf.WriteString(childText)
		}
	}

	return buf.String()
}

// extractParagraphText extracts text from paragraph nodes preserving inline formatting
func (p *GoldmarkParser) extractParagraphText(n ast.Node, source []byte) string {
	var buf bytes.Buffer

	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		switch child.Kind() {
		case ast.KindText:
			childText := p.extractText(child, source)
			buf.WriteString(childText)
		case ast.KindEmphasis:
			p.extractEmphasisText(child, source, &buf)
		case ast.KindCodeSpan:
			p.extractCodeSpanText(child, source, &buf)
		case ast.KindLink:
			p.extractLinkText(child, source, &buf)
		default:
			childText := p.extractText(child, source)
			buf.WriteString(childText)
		}
	}
	return strings.TrimSpace(buf.String())
}

// extractEmphasisText extracts text from emphasis nodes with markers
func (p *GoldmarkParser) extractEmphasisText(n ast.Node, source []byte, buf *bytes.Buffer) {
	emph := n.(*ast.Emphasis)
	marker := "*"
	if emph.Level == StrongEmphasisLevel {
		marker = "**"
	}
	buf.WriteString(marker)
	buf.WriteString(p.extractTextRecursive(n, source))
	buf.WriteString(marker)
}

// extractCodeSpanText extracts text from inline code with backticks
func (p *GoldmarkParser) extractCodeSpanText(n ast.Node, source []byte, buf *bytes.Buffer) {
	buf.WriteString("`")
	buf.WriteString(p.extractTextRecursive(n, source))
	buf.WriteString("`")
}

// extractLinkText extracts text from link nodes with markdown syntax
func (p *GoldmarkParser) extractLinkText(n ast.Node, source []byte, buf *bytes.Buffer) {
	link := n.(*ast.Link)
	buf.WriteString("[")
	buf.WriteString(p.extractTextRecursive(n, source))
	buf.WriteString("](")
	buf.Write(link.Destination)
	buf.WriteString(")")
}

// extractGenericText extracts text from other container nodes
func (p *GoldmarkParser) extractGenericText(n ast.Node, source []byte) string {
	var buf bytes.Buffer

	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if child.Kind() == ast.KindText || child.Kind() == ast.KindString {
			childText := p.extractText(child, source)
			buf.WriteString(childText)
		}
	}
	return strings.TrimSpace(buf.String())
}

// extractTextRecursive extracts text content recursively from all children
func (p *GoldmarkParser) extractTextRecursive(n ast.Node, source []byte) string {
	var buf bytes.Buffer

	switch n.Kind() {
	case ast.KindText:
		textNode := n.(*ast.Text)
		buf.Write(textNode.Segment.Value(source))
		return buf.String()
	case ast.KindString:
		str := n.(*ast.String)
		buf.Write(str.Value)
		return buf.String()
	}

	// For container nodes, extract text from all children recursively
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		childText := p.extractTextRecursive(child, source)
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
