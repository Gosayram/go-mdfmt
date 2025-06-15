package parser

import (
	"io"
)

// Parser represents a markdown parser
type Parser interface {
	// Parse parses markdown content from a reader
	Parse(r io.Reader) (Node, error)
	// ParseBytes parses markdown content from bytes
	ParseBytes(data []byte) (Node, error)
	// ParseString parses markdown content from string
	ParseString(content string) (Node, error)
}

// Options represents parser options
type Options struct {
	// Extensions specifies enabled extensions
	Extensions []string
	// Strict enables strict parsing mode
	Strict bool
}

// DefaultOptions returns default parser options
func DefaultOptions() *Options {
	return &Options{
		Extensions: []string{"table", "strikethrough", "autolink", "tasklist"},
		Strict:     false,
	}
}

// New creates a new parser with default options
func New() Parser {
	return NewWithOptions(DefaultOptions())
}

// NewWithOptions creates a new parser with specified options
func NewWithOptions(opts *Options) Parser {
	return &parser{
		options: opts,
	}
}

// parser implements the Parser interface
type parser struct {
	options *Options
}

// Parse parses markdown content from a reader
func (p *parser) Parse(r io.Reader) (Node, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return p.ParseBytes(data)
}

// ParseBytes parses markdown content from bytes
func (p *parser) ParseBytes(data []byte) (Node, error) {
	return p.ParseString(string(data))
}

// ParseString parses markdown content from string
func (p *parser) ParseString(content string) (Node, error) {
	// TODO: Implement actual parsing logic
	// For now, create a simple document with a text node
	doc := &Document{}
	if len(content) > 0 {
		text := &Text{Content: content}
		doc.Children = append(doc.Children, text)
	}
	return doc, nil
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
