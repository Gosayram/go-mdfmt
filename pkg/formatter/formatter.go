package formatter

import (
	"github.com/Gosayram/go-mdfmt/pkg/config"
	"github.com/Gosayram/go-mdfmt/pkg/parser"
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
			priority: 100,
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

	// Normalize heading style
	if cfg.Heading.Style == "atx" {
		heading.Style = "atx"
	} else {
		heading.Style = "setext"
	}

	// TODO: Implement heading level normalization
	if cfg.Heading.NormalizeLevels {
		// Normalize heading levels to prevent jumps
	}

	return nil
}

// ParagraphFormatter formats paragraph nodes
type ParagraphFormatter struct {
	BaseFormatter
}

func init() {
	// Initialize formatters with proper values
	defaultFormatters := []*BaseFormatter{
		{name: "heading", priority: 100},
		{name: "paragraph", priority: 90},
		{name: "list", priority: 80},
		{name: "code", priority: 70},
		{name: "whitespace", priority: 10},
	}
	_ = defaultFormatters // Use the formatters as needed
}

// CanFormat returns true if this formatter can handle paragraphs
func (f *ParagraphFormatter) CanFormat(nodeType parser.NodeType) bool {
	return nodeType == parser.NodeParagraph
}

// Format formats paragraph nodes
func (f *ParagraphFormatter) Format(node parser.Node, cfg *config.Config) error {
	// TODO: Implement paragraph text reflow based on line width
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
	codeBlock, ok := node.(*parser.CodeBlock)
	if !ok {
		return nil
	}

	// Set consistent fence style
	if codeBlock.Fenced {
		codeBlock.Fence = cfg.Code.FenceStyle
	}

	// TODO: Implement language detection
	if cfg.Code.LanguageDetection {
		// Auto-detect language based on content
	}

	return nil
}

// WhitespaceFormatter handles whitespace normalization
type WhitespaceFormatter struct {
	BaseFormatter
}

// CanFormat returns true for all node types (whitespace affects everything)
func (f *WhitespaceFormatter) CanFormat(nodeType parser.NodeType) bool {
	return true
}

// Format normalizes whitespace
func (f *WhitespaceFormatter) Format(node parser.Node, cfg *config.Config) error {
	// TODO: Implement whitespace normalization
	// - Remove excessive blank lines
	// - Trim trailing spaces
	// - Ensure final newline
	return nil
}
