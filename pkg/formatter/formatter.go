// Package formatter provides formatting functionality for markdown nodes.
package formatter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Gosayram/go-mdfmt/pkg/config"
	"github.com/Gosayram/go-mdfmt/pkg/parser"
)

// Constants
const (
	// HeadingFormatterPriority defines the priority for heading formatting (higher runs first)
	HeadingFormatterPriority = 100
	// ParagraphFormatterPriority defines the priority for paragraph formatting
	ParagraphFormatterPriority = 90
	// ListFormatterPriority defines the priority for list formatting
	ListFormatterPriority = 80
	// CodeFormatterPriority defines the priority for code block formatting
	CodeFormatterPriority = 70
	// WhitespaceFormatterPriority defines the priority for whitespace formatting (lowest)
	WhitespaceFormatterPriority = 10
	// InlineFormatterPriority defines the priority for inline formatting
	InlineFormatterPriority = 60

	// AtxHeadingStyle represents ATX-style heading format (# ## ###)
	AtxHeadingStyle = "atx"
	// SetextHeadingStyle represents setext-style heading format (underlined with = or -)
	SetextHeadingStyle = "setext"

	// MinHeadingLevel defines the minimum allowed heading level
	MinHeadingLevel = 1
	// MaxHeadingLevel defines the maximum allowed heading level
	MaxHeadingLevel = 6
	// SetextMaxLevel defines the maximum level for setext-style headings
	SetextMaxLevel = 2
)

// Formatter represents a markdown formatter
type Formatter interface {
	// Format formats the given AST according to configuration
	Format(root parser.Node, cfg *config.Config) error
}

// NodeFormatter represents a formatter for specific node types
type NodeFormatter interface {
	// Name returns the name of the formatter
	Name() string
	// CanFormat returns true if this formatter can handle the given node type
	CanFormat(nodeType parser.NodeType) bool
	// Format formats a specific node
	Format(node parser.Node, cfg *config.Config) error
	// Priority returns the priority of this formatter (higher = earlier)
	Priority() int
}

// Engine represents the main formatting engine
type Engine struct {
	formatters []NodeFormatter
}

// New creates a new formatting engine with default formatters
func New() *Engine {
	engine := &Engine{}
	engine.RegisterDefaults()
	return engine
}

// RegisterDefaults registers the default formatters
func (e *Engine) RegisterDefaults() {
	e.Register(&HeadingFormatter{})
	e.Register(&ParagraphFormatter{})
	e.Register(&ListFormatter{})
	e.Register(&CodeBlockFormatter{})
	e.Register(&InlineFormatter{})
	e.Register(&WhitespaceFormatter{})
}

// Register registers a new node formatter
func (e *Engine) Register(formatter NodeFormatter) {
	e.formatters = append(e.formatters, formatter)
	// Sort by priority
	for i := len(e.formatters) - 1; i > 0; i-- {
		if e.formatters[i].Priority() > e.formatters[i-1].Priority() {
			e.formatters[i], e.formatters[i-1] = e.formatters[i-1], e.formatters[i]
		} else {
			break
		}
	}
}

// Format formats the given AST according to configuration
func (e *Engine) Format(doc *parser.Document, cfg *config.Config) error {
	walker := parser.NewWalker(doc)

	for node, ok := walker.Next(); ok; node, ok = walker.Next() {
		for _, formatter := range e.formatters {
			if formatter.CanFormat(node.Type()) {
				if err := formatter.Format(node, cfg); err != nil {
					return err
				}
				break // Only apply first matching formatter
			}
		}
	}

	return nil
}

// BaseFormatter provides common functionality for formatters
type BaseFormatter struct {
	name     string
	priority int
}

// Name returns the formatter name
func (f *BaseFormatter) Name() string {
	return f.name
}

// Priority returns the formatter priority
func (f *BaseFormatter) Priority() int {
	return f.priority
}

// HeadingFormatter formats heading nodes
type HeadingFormatter struct {
	BaseFormatter
}

// NewHeadingFormatter creates a new heading formatter
func NewHeadingFormatter() *HeadingFormatter {
	return &HeadingFormatter{
		BaseFormatter: BaseFormatter{
			name:     "heading",
			priority: HeadingFormatterPriority,
		},
	}
}

// CanFormat returns true if this formatter can handle headings
func (f *HeadingFormatter) CanFormat(nodeType parser.NodeType) bool {
	return nodeType == parser.NodeHeading
}

// Format applies heading formatting rules.
func (f *HeadingFormatter) Format(node parser.Node, cfg *config.Config) error {
	heading, ok := node.(*parser.Heading)
	if !ok {
		return nil
	}

	// Apply heading style preferences
	switch cfg.Heading.Style {
	case AtxHeadingStyle:
		// Ensure ATX-style headers (#, ##, ###, etc.)
		heading.Style = AtxHeadingStyle
	case SetextHeadingStyle:
		// Use setext style for levels 1 and 2, ATX for others
		if heading.Level <= SetextMaxLevel {
			heading.Style = SetextHeadingStyle
		} else {
			heading.Style = AtxHeadingStyle
		}
	}

	// Apply heading level normalization if enabled
	if cfg.Heading.NormalizeLevels {
		// Ensure heading level doesn't exceed HTML limit
		if heading.Level > MaxHeadingLevel {
			heading.Level = MaxHeadingLevel
		}
		// Ensure heading level is at least minimum
		if heading.Level < MinHeadingLevel {
			heading.Level = MinHeadingLevel
		}
	}

	// Clean up heading text (trim whitespace)
	heading.Text = strings.TrimSpace(heading.Text)

	return nil
}

// ParagraphFormatter formats paragraph nodes
type ParagraphFormatter struct {
	BaseFormatter
}

// CanFormat returns true if this formatter can handle paragraphs
func (f *ParagraphFormatter) CanFormat(nodeType parser.NodeType) bool {
	return nodeType == parser.NodeParagraph
}

// Format applies paragraph formatting rules with text reflow.
func (f *ParagraphFormatter) Format(node parser.Node, cfg *config.Config) error {
	paragraph, ok := node.(*parser.Paragraph)
	if !ok {
		return nil
	}

	// Apply text reflow if line width is configured
	if cfg.LineWidth > 0 {
		paragraph.Text = f.wrapText(paragraph.Text, cfg.LineWidth)
	}

	// Clean up excessive whitespace
	paragraph.Text = strings.TrimSpace(paragraph.Text)
	// Replace multiple spaces with single space
	paragraph.Text = normalizeWhitespace(paragraph.Text)

	return nil
}

// wrapText wraps text to the specified line width
func (f *ParagraphFormatter) wrapText(text string, width int) string {
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

// normalizeWhitespace replaces multiple consecutive spaces with single spaces
func normalizeWhitespace(text string) string {
	// Replace multiple spaces/tabs with single space
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		// Replace multiple whitespace characters with single space
		fields := strings.Fields(line)
		lines[i] = strings.Join(fields, " ")
	}
	return strings.Join(lines, "\n")
}

// ListFormatter formats list nodes
type ListFormatter struct {
	BaseFormatter
}

// CanFormat returns true if this formatter can handle lists
func (f *ListFormatter) CanFormat(nodeType parser.NodeType) bool {
	return nodeType == parser.NodeList || nodeType == parser.NodeListItem
}

// Format applies list formatting rules.
func (f *ListFormatter) Format(node parser.Node, cfg *config.Config) error {
	switch n := node.(type) {
	case *parser.List:
		return f.formatList(n, cfg)
	case *parser.ListItem:
		return f.formatListItem(n, cfg)
	}
	return nil
}

// formatList handles formatting of list nodes
func (f *ListFormatter) formatList(list *parser.List, cfg *config.Config) error {
	if !list.Ordered {
		f.formatUnorderedList(list, cfg)
	} else {
		f.formatOrderedList(list, cfg)
	}

	return f.processListItems(list, cfg)
}

// formatUnorderedList sets consistent bullet style for unordered lists
func (f *ListFormatter) formatUnorderedList(list *parser.List, cfg *config.Config) {
	list.Marker = cfg.List.BulletStyle
	// Apply the same marker to all items
	for _, item := range list.Items {
		item.Marker = cfg.List.BulletStyle
	}
}

// formatOrderedList sets consistent numbering for ordered lists
func (f *ListFormatter) formatOrderedList(list *parser.List, cfg *config.Config) {
	for i, item := range list.Items {
		switch cfg.List.NumberStyle {
		case ".":
			item.Marker = fmt.Sprintf("%d.", i+1)
		case ")":
			item.Marker = fmt.Sprintf("%d)", i+1)
		default:
			item.Marker = fmt.Sprintf("%d.", i+1)
		}
	}
}

// processListItems handles list item processing and nested lists
func (f *ListFormatter) processListItems(list *parser.List, cfg *config.Config) error {
	for _, item := range list.Items {
		if cfg.List.ConsistentIndentation {
			// Normalize list item text (trim and clean whitespace)
			item.Text = strings.TrimSpace(item.Text)
			item.Text = normalizeWhitespace(item.Text)
		}

		// Process nested lists recursively
		if err := f.processNestedLists(item, cfg); err != nil {
			return err
		}
	}
	return nil
}

// processNestedLists handles nested lists within list items
func (f *ListFormatter) processNestedLists(item *parser.ListItem, cfg *config.Config) error {
	for _, child := range item.Children {
		if childList, ok := child.(*parser.List); ok {
			if err := f.Format(childList, cfg); err != nil {
				return err
			}
		}
	}
	return nil
}

// formatListItem handles formatting of individual list items
func (f *ListFormatter) formatListItem(item *parser.ListItem, cfg *config.Config) error {
	// Individual list item formatting
	item.Text = strings.TrimSpace(item.Text)
	item.Text = normalizeWhitespace(item.Text)

	// Process nested lists in this item
	return f.processNestedLists(item, cfg)
}

// CodeBlockFormatter formats code block nodes
type CodeBlockFormatter struct {
	BaseFormatter
}

// CanFormat returns true if this formatter can handle code blocks
func (f *CodeBlockFormatter) CanFormat(nodeType parser.NodeType) bool {
	return nodeType == parser.NodeCodeBlock
}

// Format formats code block nodes
func (f *CodeBlockFormatter) Format(node parser.Node, cfg *config.Config) error {
	code, ok := node.(*parser.CodeBlock)
	if !ok {
		return nil
	}

	// Apply fence style preferences
	switch cfg.Code.FenceStyle {
	case "```":
		code.Fence = "```"
	case "~~~":
		code.Fence = "~~~"
	}

	// Language detection is not implemented yet
	_ = cfg.Code.LanguageDetection

	return nil
}

// WhitespaceFormatter handles whitespace normalization
type WhitespaceFormatter struct {
	BaseFormatter
}

// CanFormat returns true for all node types (whitespace affects everything)
func (f *WhitespaceFormatter) CanFormat(_ parser.NodeType) bool {
	return true // Whitespace formatter can format any node
}

// Format applies whitespace normalization rules.
func (f *WhitespaceFormatter) Format(node parser.Node, cfg *config.Config) error {
	// Apply whitespace rules based on node type
	switch n := node.(type) {
	case *parser.Document:
		// Document-level whitespace normalization
		f.normalizeDocumentWhitespace(n, cfg)
	case *parser.Paragraph:
		// Normalize paragraph whitespace
		if cfg.Whitespace.TrimTrailingSpaces {
			n.Text = f.trimTrailingSpaces(n.Text)
		}
	case *parser.Heading:
		// Normalize heading whitespace
		if cfg.Whitespace.TrimTrailingSpaces {
			n.Text = strings.TrimSpace(n.Text)
		}
	case *parser.CodeBlock:
		// For code blocks, be more careful with whitespace
		if cfg.Whitespace.TrimTrailingSpaces {
			// Only trim trailing spaces at the end of lines, preserve indentation
			lines := strings.Split(n.Content, "\n")
			for i, line := range lines {
				lines[i] = strings.TrimRight(line, " \t")
			}
			n.Content = strings.Join(lines, "\n")
		}
	case *parser.Text:
		// Normalize text node whitespace
		if cfg.Whitespace.TrimTrailingSpaces {
			n.Content = f.trimTrailingSpaces(n.Content)
		}
	}

	return nil
}

// normalizeDocumentWhitespace handles document-level whitespace rules
func (f *WhitespaceFormatter) normalizeDocumentWhitespace(_ *parser.Document, _ *config.Config) {
	// This would be used for limiting excessive blank lines between elements
	// For now, we'll handle this in the renderer
}

// trimTrailingSpaces removes trailing spaces from each line
func (f *WhitespaceFormatter) trimTrailingSpaces(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}

// InlineFormatter handles inline elements like links, emphasis, and inline code
type InlineFormatter struct {
	BaseFormatter
}

// NewInlineFormatter creates a new inline formatter
func NewInlineFormatter() *InlineFormatter {
	return &InlineFormatter{
		BaseFormatter: BaseFormatter{
			name:     "inline",
			priority: InlineFormatterPriority,
		},
	}
}

// CanFormat returns true if this formatter can handle text nodes (where inline elements are)
func (f *InlineFormatter) CanFormat(nodeType parser.NodeType) bool {
	return nodeType == parser.NodeText || nodeType == parser.NodeParagraph
}

// Format applies inline formatting rules
func (f *InlineFormatter) Format(node parser.Node, _ *config.Config) error {
	var text string

	switch n := node.(type) {
	case *parser.Text:
		text = n.Content
		text = f.normalizeInlineElements(text)
		n.Content = text
	case *parser.Paragraph:
		text = n.Text
		text = f.normalizeInlineElements(text)
		n.Text = text
	default:
		return nil
	}

	return nil
}

// normalizeInlineElements cleans up inline markdown formatting
func (f *InlineFormatter) normalizeInlineElements(text string) string {
	// Normalize inline code backticks (ensure single backticks for simple inline code)
	text = f.normalizeInlineCode(text)

	// Normalize emphasis and strong formatting
	text = f.normalizeEmphasis(text)

	// Clean up link formatting
	text = f.normalizeLinks(text)

	return text
}

// normalizeInlineCode ensures consistent inline code formatting
func (f *InlineFormatter) normalizeInlineCode(text string) string {
	// Replace multiple backticks with single backticks where appropriate
	// This is a simplified implementation

	// Remove spaces around inline code
	re := regexp.MustCompile("`\\s+([^`]+)\\s+`")
	text = re.ReplaceAllString(text, "`$1`")

	return text
}

// normalizeEmphasis ensures consistent emphasis formatting
func (f *InlineFormatter) normalizeEmphasis(text string) string {
	// Normalize emphasis to use asterisks consistently
	// Convert _text_ to *text*
	re := regexp.MustCompile(`\b_([^_]+)_\b`)
	text = re.ReplaceAllString(text, "*$1*")

	// Normalize strong emphasis **text** (keep as is, it's already correct)
	return text
}

// normalizeLinks cleans up link formatting
func (f *InlineFormatter) normalizeLinks(text string) string {
	// Ensure proper spacing around links
	// This is a basic implementation - could be enhanced

	// Remove extra spaces in link text
	re := regexp.MustCompile(`\[\s+([^\]]+)\s+\]`)
	text = re.ReplaceAllString(text, "[$1]")

	return text
}
