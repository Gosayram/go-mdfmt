package parser

// Parser interface defines methods for parsing Markdown content and validating the parser
type Parser interface {
	Parse(content []byte) (*Document, error)
	Validate() error
}

// New creates a new parser instance using the default parser implementation
func New() Parser {
	return NewGoldmarkParser()
}

// DefaultParser returns the default parser implementation using the Goldmark parser
func DefaultParser() Parser {
	return NewGoldmarkParser()
}

// BasicParser represents a simple placeholder parser (for testing)
type BasicParser struct{}

// NewBasicParser creates a new basic parser
func NewBasicParser() *BasicParser {
	return &BasicParser{}
}

// Parse implements a basic placeholder parser
func (p *BasicParser) Parse(content []byte) (*Document, error) {
	// Create a simple document with one paragraph
	doc := &Document{
		Children: []Node{
			&Paragraph{
				Text: string(content),
			},
		},
	}
	return doc, nil
}

// Validate validates the parser configuration
func (p *BasicParser) Validate() error {
	return nil
}

// Helper functions for node manipulation

// FindNodes finds all nodes of a specific type in the tree
func FindNodes(doc *Document, nodeType NodeType) []Node {
	var found []Node
	walker := NewWalker(doc)

	for node, ok := walker.Next(); ok; node, ok = walker.Next() {
		if node.Type() == nodeType {
			found = append(found, node)
		}
	}

	return found
}

// FindFirstNode finds the first node of a specific type
func FindFirstNode(doc *Document, nodeType NodeType) Node {
	walker := NewWalker(doc)

	for node, ok := walker.Next(); ok; node, ok = walker.Next() {
		if node.Type() == nodeType {
			return node
		}
	}

	return nil
}
