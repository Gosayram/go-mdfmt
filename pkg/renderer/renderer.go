// Package renderer provides markdown rendering functionality.
package renderer

import (
	"io"
	"regexp"
	"strings"

	"github.com/Gosayram/go-mdfmt/pkg/config"
	"github.com/Gosayram/go-mdfmt/pkg/parser"
)

const (
	// SecondHeadingLevel represents heading level 2
	SecondHeadingLevel = 2
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

// Render renders the AST to markdown string with whitespace normalization.
func (r *MarkdownRenderer) Render(doc *parser.Document, cfg *config.Config) (string, error) {
	r.output.Reset()
	r.config = cfg

	if err := r.renderDocument(doc, 0); err != nil {
		return "", err
	}

	result := r.output.String()

	// Apply document-level whitespace rules
	result = r.normalizeBlankLines(result, cfg.Whitespace.MaxBlankLines)

	// Ensure final newline if configured
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
func (r *MarkdownRenderer) renderHeading(heading *parser.Heading, _ int) error {
	if heading.Style == "setext" && heading.Level <= SecondHeadingLevel {
		// Setext-style heading
		r.output.WriteString(heading.Text)
		r.output.WriteString("\n")

		marker := "="
		if heading.Level == SecondHeadingLevel {
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
func (r *MarkdownRenderer) renderParagraph(para *parser.Paragraph, _ int) error {
	content := para.Text

	// Fix broken markdown links first
	content = r.fixBrokenLinks(content)

	// Apply line width wrapping only if no markdown links are present
	if r.config.LineWidth > 0 && !r.containsMarkdownLinks(content) {
		content = r.wrapText(content, r.config.LineWidth)
	}

	r.output.WriteString(content)
	r.output.WriteString("\n\n")

	return nil
}

// containsMarkdownLinks checks if text contains markdown links
func (r *MarkdownRenderer) containsMarkdownLinks(text string) bool {
	linkPattern := `\[[^\]]*\]\([^)]*\)`
	matched, _ := regexp.MatchString(linkPattern, text)
	// Debug: uncomment to see what's happening
	// fmt.Printf("DEBUG: Text: %q, Contains links: %v\n", text, matched)
	return matched
}

// fixBrokenLinks repairs markdown links that have been broken across lines
func (r *MarkdownRenderer) fixBrokenLinks(text string) string {
	const (
		// Expected parts count for broken link pattern: [text1]\n[text2](url)
		brokenLinkPartsCount = 4
		// Expected parts count for multi-break pattern: [text...](url)
		multiBreakPartsCount = 3
		// Indices for regex match parts
		linkTextPart1Index = 1
		linkTextPart2Index = 2
		urlPartIndex       = 3
		linkTextIndex      = 1
		urlIndex           = 2
	)

	// Pattern to match broken links: [text\nmore text](url)
	brokenLinkPattern := `\[([^\]]*)\n([^\]]*)\]\(([^)]*)\)`
	re := regexp.MustCompile(brokenLinkPattern)

	// Replace broken links with fixed ones
	fixed := re.ReplaceAllStringFunc(text, func(match string) string {
		// Extract parts
		parts := re.FindStringSubmatch(match)
		if len(parts) == brokenLinkPartsCount {
			linkText := parts[linkTextPart1Index] + " " + parts[linkTextPart2Index] // Join with space
			url := parts[urlPartIndex]
			return "[" + linkText + "](" + url + ")"
		}
		return match
	})

	// Also handle cases where the link text itself has multiple line breaks
	multiBreakPattern := `\[([^\]]*(?:\n[^\]]*)*)\]\(([^)]*)\)`
	re2 := regexp.MustCompile(multiBreakPattern)

	fixed = re2.ReplaceAllStringFunc(fixed, func(match string) string {
		parts := re2.FindStringSubmatch(match)
		if len(parts) == multiBreakPartsCount {
			linkText := strings.ReplaceAll(parts[linkTextIndex], "\n", " ")
			url := parts[urlIndex]
			return "[" + linkText + "](" + url + ")"
		}
		return match
	})

	return fixed
}

// renderList renders a list node
func (r *MarkdownRenderer) renderList(list *parser.List, depth int) error {
	for _, item := range list.Items {
		if err := r.renderListItem(item, depth+1); err != nil {
			return err
		}
	}

	r.output.WriteString("\n")
	return nil
}

// renderListItem renders a list item node
func (r *MarkdownRenderer) renderListItem(item *parser.ListItem, depth int) error {
	// Use proper indentation for nested lists only (depth > 1)
	indent := ""
	if depth > 1 {
		indent = strings.Repeat("  ", depth-1)
	}

	// Determine marker
	marker := item.Marker
	if marker == "" {
		marker = r.config.List.BulletStyle
	}

	r.output.WriteString(indent)
	r.output.WriteString(marker)
	r.output.WriteString(" ")
	r.output.WriteString(item.Text)

	// Render nested elements
	if len(item.Children) > 0 {
		r.output.WriteString("\n")
		for _, child := range item.Children {
			if err := r.renderNode(child, depth); err != nil {
				return err
			}
		}
	} else {
		// Add newline after list item text if no nested elements
		r.output.WriteString("\n")
	}

	return nil
}

// renderCodeBlock renders a code block node
func (r *MarkdownRenderer) renderCodeBlock(code *parser.CodeBlock, _ int) error {
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
func (r *MarkdownRenderer) renderText(text *parser.Text, _ int) error {
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

// wrapText wraps text to the specified line width, preserving markdown links
func (r *MarkdownRenderer) wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	// Split text into tokens, preserving markdown links as single units
	tokens := r.tokenizeWithLinks(text)
	if len(tokens) == 0 {
		return text
	}

	var lines []string
	var currentLine strings.Builder

	for i, token := range tokens {
		// Check if adding this token would exceed the line width
		if currentLine.Len() > 0 && currentLine.Len()+1+len(token) > width {
			// Always start new line when width exceeded (for both links and regular words)
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(token)

		// If this is the last token, add the current line
		if i == len(tokens)-1 {
			lines = append(lines, currentLine.String())
		}
	}

	return strings.Join(lines, "\n")
}

// tokenizeWithLinks splits text into words while keeping markdown links intact
func (r *MarkdownRenderer) tokenizeWithLinks(text string) []string {
	// Simple regex-based approach to find markdown links
	linkPattern := `\[[^\]]*\]\([^)]*\)`
	re := regexp.MustCompile(linkPattern)

	var tokens []string
	lastEnd := 0

	// Find all links and process text between them
	matches := re.FindAllStringIndex(text, -1)

	for _, match := range matches {
		start, end := match[0], match[1]

		// Add words before the link
		if start > lastEnd {
			beforeLink := text[lastEnd:start]
			words := strings.Fields(beforeLink)
			tokens = append(tokens, words...)
		}

		// Add the link as a single token
		link := text[start:end]
		tokens = append(tokens, link)

		lastEnd = end
	}

	// Add remaining words after the last link
	if lastEnd < len(text) {
		afterLinks := text[lastEnd:]
		words := strings.Fields(afterLinks)
		tokens = append(tokens, words...)
	}

	// If no links found, just split into words
	if len(matches) == 0 {
		tokens = strings.Fields(text)
	}

	return tokens
}

// normalizeBlankLines limits consecutive blank lines to the configured maximum
func (r *MarkdownRenderer) normalizeBlankLines(text string, maxBlankLines int) string {
	if maxBlankLines < 0 {
		return text
	}

	lines := strings.Split(text, "\n")
	var result []string
	consecutiveEmpty := 0

	for _, line := range lines {
		isEmpty := strings.TrimSpace(line) == ""

		if isEmpty {
			consecutiveEmpty++
			// Only add empty line if we haven't exceeded the limit
			if consecutiveEmpty <= maxBlankLines {
				result = append(result, line)
			}
		} else {
			consecutiveEmpty = 0
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
