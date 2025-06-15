// Package formatter provides formatting functionality for markdown nodes.
package formatter

import (
	"github.com/Gosayram/go-mdfmt/pkg/config"
	"github.com/Gosayram/go-mdfmt/pkg/parser"
)

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

// Format formats heading nodes
func (f *HeadingFormatter) Format(node parser.Node, cfg *config.Config) error {
	heading, ok := node.(*parser.Heading)
	if !ok {
		return nil
	}

	// Apply heading style preferences
	if cfg.Heading.Style == "atx" {
		// Ensure ATX-style headers (#, ##, ###, etc.)
		// Set the style on the heading
		_ = heading // Use heading to avoid unused variable error
	}

	// Add heading level normalization logic if needed
	if cfg.Heading.NormalizeLevels {
		// Normalize heading levels to prevent jumps
		// Implementation would go here
		_ = heading // Use heading to avoid unused variable error
	}

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

// Format formats paragraph nodes
func (f *ParagraphFormatter) Format(_ parser.Node, _ *config.Config) error {
	// Implementation would go here
	return nil
}

// ListFormatter formats list nodes
type ListFormatter struct {
	BaseFormatter
}

// CanFormat returns true if this formatter can handle lists
func (f *ListFormatter) CanFormat(nodeType parser.NodeType) bool {
	return nodeType == parser.NodeList || nodeType == parser.NodeListItem
}

// Format formats list nodes
func (f *ListFormatter) Format(node parser.Node, cfg *config.Config) error {
	switch n := node.(type) {
	case *parser.List:
		// Set consistent bullet style
		if !n.Ordered {
			n.Marker = cfg.List.BulletStyle
		}
		// TODO: Implement consistent indentation
	case *parser.ListItem:
		// Format list item marker
		// TODO: Implement parent-child relationship if needed
		n.Marker = cfg.List.BulletStyle
	}
	return nil
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
	if cfg.Code.FenceStyle == "```" {
		code.Fence = "```"
	} else if cfg.Code.FenceStyle == "~~~" {
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

// Format normalizes whitespace
func (f *WhitespaceFormatter) Format(_ parser.Node, _ *config.Config) error {
	// Implementation would go here
	return nil
}
