// Package parser provides markdown parsing functionality and Abstract Syntax Tree definitions.
package parser

import (
	"fmt"
	"strings"
)

// NodeType represents the type of a node in the AST
type NodeType int

// NodeType constants
const (
	// NodeDocument represents a document node containing the entire markdown structure
	NodeDocument NodeType = iota
	// NodeHeading represents a heading node (# Title)
	NodeHeading
	// NodeParagraph represents a paragraph of text
	NodeParagraph
	// NodeList represents an ordered or unordered list
	NodeList
	// NodeListItem represents a single item within a list
	NodeListItem
	// NodeCodeBlock represents a code block (fenced or indented)
	NodeCodeBlock
	// NodeText represents plain text content
	NodeText
)

// Node represents a basic node in the markdown AST
type Node interface {
	Type() NodeType
	String() string
}

// Document represents the root document node
type Document struct {
	Children []Node
}

// Type returns the node type for Document nodes.
func (n *Document) Type() NodeType { return NodeDocument }
func (n *Document) String() string { return "Document" }

// Heading represents a heading node
type Heading struct {
	Level int
	Text  string
	Style string // "atx" or "setext"
}

// Type returns the node type for Heading nodes.
func (n *Heading) Type() NodeType { return NodeHeading }
func (n *Heading) String() string {
	return fmt.Sprintf("Heading(level=%d, text=%q)", n.Level, n.Text)
}

// Paragraph represents a paragraph node
type Paragraph struct {
	Text string
}

// Type returns the node type for Paragraph nodes.
func (n *Paragraph) Type() NodeType { return NodeParagraph }
func (n *Paragraph) String() string {
	return fmt.Sprintf("Paragraph(text=%q)", n.Text)
}

// List represents a list node
type List struct {
	Ordered bool
	Items   []*ListItem
	Marker  string
}

// Type returns the node type for List nodes.
func (n *List) Type() NodeType { return NodeList }
func (n *List) String() string {
	return fmt.Sprintf("List(ordered=%t, items=%d)", n.Ordered, len(n.Items))
}

// ListItem represents a list item node
type ListItem struct {
	Text     string
	Marker   string
	Children []Node // Support for nested lists and other elements
}

// Type returns the node type for ListItem nodes.
func (n *ListItem) Type() NodeType { return NodeListItem }
func (n *ListItem) String() string {
	return fmt.Sprintf("ListItem(text=%q)", n.Text)
}

// CodeBlock represents a code block node
type CodeBlock struct {
	Language string
	Content  string
	Fenced   bool
	Fence    string
}

// Type returns the node type for CodeBlock nodes.
func (n *CodeBlock) Type() NodeType { return NodeCodeBlock }
func (n *CodeBlock) String() string {
	return fmt.Sprintf("CodeBlock(lang=%q, fenced=%t)", n.Language, n.Fenced)
}

// Text represents a text node
type Text struct {
	Content string
}

// Type returns the node type for Text nodes.
func (n *Text) Type() NodeType { return NodeText }
func (n *Text) String() string {
	return fmt.Sprintf("Text(content=%q)", n.Content)
}

// Walker provides a simple way to iterate over nodes
type Walker struct {
	nodes []Node
	index int
}

// NewWalker creates a new walker for the given document
func NewWalker(doc *Document) *Walker {
	nodes := append([]Node{doc}, doc.Children...)
	return &Walker{nodes: nodes, index: -1}
}

// Next returns the next node in the walk
func (w *Walker) Next() (Node, bool) {
	w.index++
	if w.index >= len(w.nodes) {
		return nil, false
	}
	return w.nodes[w.index], true
}

// NodeTypeString returns a string representation of the node type
func NodeTypeString(t NodeType) string {
	switch t {
	case NodeDocument:
		return "Document"
	case NodeHeading:
		return "Heading"
	case NodeParagraph:
		return "Paragraph"
	case NodeList:
		return "List"
	case NodeListItem:
		return "ListItem"
	case NodeCodeBlock:
		return "CodeBlock"
	case NodeText:
		return "Text"
	default:
		return "Unknown"
	}
}

// DebugString returns a debug representation of a document
func DebugString(doc *Document) string {
	var sb strings.Builder
	sb.WriteString("Document\n")
	for _, child := range doc.Children {
		sb.WriteString("  ")
		sb.WriteString(child.String())
		sb.WriteString("\n")
	}
	return sb.String()
}

// GetAllNodes returns all nodes in the document as a flat slice.
func (n *Document) GetAllNodes() []Node {
	return append([]Node{}, n.Children...)
}
