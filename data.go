package main

import (
	 "fmt"
)

type AgendaNode struct {
	Title string
	Text string
	Sibling *AgendaNode
	Child *AgendaNode
	Tags []string
}

func main() {
	root := NewNode("Root", "Root Text", []string{})
	child1a := NewNode("Child1a", "Child1a Text", []string{})
	root.AddChild(child1a)
	child1b := NewNode("Child1b", "Child1b Text", []string{})
	root.AddChild(child1b)
	child2a := NewNode("Child2a", "Child2a Text", []string{})
	child1a.AddChild(child2a)

	print := func(node *AgendaNode)

	WalkDepth()
}

func (node *AgendaNode) Display(level int) {
	fmt.Println(node.Title)
	fmt.Println(> text)
}

func NewNode(title, text string, tags []string) *AgendaNode {
	new := &AgendaNode{Title: title, Text: text, Tags: tags}
	return new
}

func (parent *AgendaNode) AddChild(node *AgendaNode) {
	var sibling *AgendaNode
	if parent.Child != nil {
		sibling = parent.Child
		sibling.AddSibling(node)
	} else {
		parent.Child = node
	}
}

func (sibling *AgendaNode) AddSibling(node *AgendaNode) {
	for ; sibling.Sibling != nil; sibling = sibling.Sibling {}
	sibling.Sibling = node
}

// Traverses child nodes until nil is encountered. No siblings are touched.
//
func (node *AgendaNode) VisitDescendants(callback func(node *AgendaNode)) {
	func walk := func(level int, node *AgendaNode, callback func(node *AgendaNode)) {
		callback(node)
		
		if node.Child == nil {
			return
		}

		walk(level + 1, node.Child, callback)
		return
	}

	walk(0, node, callback)
}

// Traverses sibling nodes until nil is encountered. No children are touched.
//
func (tree *AgendaNode) VisitSiblings(callback func(node *AgendaNode)) {
	func walk := func(node *AgendaNode, callback func(node *AgendaNode)) {
		for ; node.Sibling != nil; node = node.Sibling {
			callback(node)
		}
	}

	walk(node, callback)
}

func (tree *AgendaNode) WalkDepth(callback func(level, *AgendaNode)) {
	func walk := func(level int, tree *AgendaNode, callback func(level, *AgendaNode)) {
		node := tree

		depthLevel := level
		for ; depthLevel++, node.Child != nil; node = node.Child {
			callback(depthLevel, node)
		}

		for ; node.Sibling != nil; node = node.Sibling {
			node.WalkDepth(level, callback)
		}
	}

	walk(0, tree, callback)
}
