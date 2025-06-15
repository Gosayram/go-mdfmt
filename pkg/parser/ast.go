package parser

import (
	"fmt"
	"strings"
)

// NodeType represents the type of a node in the AST
type NodeType int

const (
	// Document node types
	NodeDocument NodeType = iota
	NodeHeading
	NodeParagraph
	NodeList
	NodeListItem
	NodeCodeBlock
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

func (n *Document) Type() NodeType { return NodeDocument }
func (n *Document) String() string { return "Document" }

// Heading represents a heading node
type Heading struct {
	Level int
	Text  string
	Style string // "atx" or "setext"
}

func (n *Heading) Type() NodeType { return NodeHeading }
func (n *Heading) String() string {
	return fmt.Sprintf("Heading(level=%d, text=%q)", n.Level, n.Text)
}

// Paragraph represents a paragraph node
type Paragraph struct {
	Text string
}

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

func (n *List) Type() NodeType { return NodeList }
func (n *List) String() string {
	return fmt.Sprintf("List(ordered=%t, items=%d)", n.Ordered, len(n.Items))
}

// ListItem represents a list item node
type ListItem struct {
	Text   string
	Marker string
}

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

func (n *CodeBlock) Type() NodeType { return NodeCodeBlock }
func (n *CodeBlock) String() string {
	return fmt.Sprintf("CodeBlock(lang=%q, fenced=%t)", n.Language, n.Fenced)
}

// Text represents a text node
type Text struct {
	Content string
}

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
	var nodes []Node
	nodes = append(nodes, doc)
	for _, child := range doc.Children {
		nodes = append(nodes, child)
	}
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
