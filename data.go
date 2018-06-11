package main

/*
Use this sample for reference:
+-------------------------------------------------------------------------------
|* Heading 1
|  text 1-1
|
|  * Sub-Heading 1a
|    text 1a-1
|
|    * Sub-Sub Heading 1aa
|      text 1aa-1
|
|    text 1a-2
|
|  text 1-2
|
|  * Sub-Heading 1b
|    text 1b-1
|
|  text 1-3
+-------------------------------------------------------------------------------

text 1-1, text 1-2 and text 1-3 are all Siblings of each other.
text 1-2 and text 1-3 should not have headings of their own.

text 1a-1 and text 1b-1 are children of text 1-1 (but *not* 1-2 or 1-3)

*/

import (
	"fmt"
)

type AgendaNode struct {
	Title    string
	Text     string
	Sibling  *AgendaNode
	Children []*AgendaNode
	Tags     []string
}

func main() {
	level0 := NewNode("level0", "text0 text0 text0 text0", []string{})

	level1 := NewNode("level1", "text1 text1 text1 text1", []string{})
	level0.AddChild(level1)

	level1b := NewNode("level1b", "text1b text1b text1b text1b", []string{})
	level1.AddSibling(level1b)

	level1c := NewNode("level1c", "text1c text1c text1c text1c", []string{})
	level1.AddSibling(level1c)

	level2 := NewNode("level2", "text2 text2 text2 text2", []string{})
	level1.AddChild(level2)

	print := func(node *AgendaNode, depth int) {
		fmt.Printf("%*s%v\n", depth*5, " ", node.Title)
		fmt.Printf("%*s%v\n", depth*5, " ", node.Text)
	}

	level0.Walk(print)
}

func (node *AgendaNode) Display(level int) {
	fmt.Println(node.Title)
	fmt.Println("> text")
}

func NewNode(title, text string, tags []string) *AgendaNode {
	new := &AgendaNode{Title: title, Text: text, Tags: tags}
	return new
}

func (parent *AgendaNode) AddChild(node *AgendaNode) {
	parent.Children = append(node.Children, node)
}

func (sibling *AgendaNode) AddSibling(node *AgendaNode) {
	for ; sibling.Sibling != nil; sibling = sibling.Sibling {
	}
	sibling.Sibling = node
}

// Traverses sibling nodes until nil is encountered. No children are touched.
//
func (node *AgendaNode) VisitSiblings(siblingNum int, callback func(*AgendaNode, int)) {
	for ; node.Sibling != nil; node = node.Sibling {
		callback(node, siblingNum)
		siblingNum++
	}
}

// Traverses child nodes until nil is encountered. No siblings are touched.
//
func (node *AgendaNode) VisitChildren(childNum int, callback func(*AgendaNode, int)) {
	for i := range node.Children {
		child := node.Children[i]
		callback(child, i)
	}
}

// Invoke callback on each node in the tree.
// Walk performs a depth-first, sibling-second traversal of the tree.
// @param callback
//        Takes the node being visited and how many levels deep in the tree the node is.
//
func (node *AgendaNode) Walk(callback func(*AgendaNode, int)) {
	type cbfn func(*AgendaNode, int)
	var walk func(*AgendaNode, int, cbfn)

	walk = func(node *AgendaNode, depth int, callback cbfn) {
		callback(node, depth)

		for i := range node.Children {
			child := node.Children[i]
			walk(child, depth+1, callback)
		}

		for ; node.Sibling != nil; node = node.Sibling {
			walk(node.Sibling, depth, callback)
		}
	}

	walk(node, 0, callback)
}
