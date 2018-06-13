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
	"io"
)

type AgendaNode struct {
	Parent      *AgendaNode
	Title       string
	Text        string
	NextSibling *AgendaNode
	PrevSibling *AgendaNode
	Children    []*AgendaNode
	Tags        []string
}

// func main() {
// r := NewNode("r", "r r r", []string{})

// rc1 := NewNode("rc1", "rc1 rc1 rc1", []string{})
// r.AddChild(rc1)

// rc1s1 := NewNode("rc1s1", "rc1s1 rc1s1 rc1s1", []string{})
// rc1.AddSibling(rc1s1)

// rc1s1c1 := NewNode("rc1s1c1", "rc1s1c1 rc1s1c1 rc1s1c1", []string{})
// rc1s1.AddChild(rc1s1c1)

// rc1s2 := NewNode("rc1s2", "rc1s2 rc1s2 rc1s2", []string{})
// rc1s1.AddSibling(rc1s2)

// rc2 := NewNode("rc2", "rc2 rc2 rc2", []string{})
// r.AddChild(rc2)

// rc1s3 := NewNode("rc1s3", "rc1s3 rc1s3 rc1s3", []string{})
// rc1s1.AddSibling(rc1s3)

// rc1s3c1 := NewNode("rc1s3c1", "rc1s3c1 rc1s3c1 rc1s3c1", []string{})
// rc1s3c1.AddSibling(rc1s3c1)

// r.PrintTree(os.Stdout, 5)
// }

func (tree *AgendaNode) PrintTree(w io.Writer, indent int) {
	tree.Walk(func(node *AgendaNode, depth int) {
		node.Write(w, depth, indent)
	})
}

func (node *AgendaNode) Write(w io.Writer, indentLevel int, indentScale int) {
	io.WriteString(w, fmt.Sprintf("%*s%v\n", indentLevel*indentScale, " ", node.Title))
	io.WriteString(w, fmt.Sprintf("%*s%v\n", indentLevel*indentScale, " ", node.Text))
}

func NewNode(title, text string, tags []string) *AgendaNode {
	new := &AgendaNode{Title: title, Text: text, Tags: tags}
	return new
}

func (parent *AgendaNode) AddChild(node *AgendaNode) {
	parent.Children = append(parent.Children, node)
	node.Parent = parent
}

func (node *AgendaNode) AddSibling(new *AgendaNode) {
	for ; node.NextSibling != nil; node = node.NextSibling {
	}
	node.NextSibling = new
	new.PrevSibling = node
	new.Parent = node.Parent
}

func (parent *AgendaNode) IndexChild(child *AgendaNode) int {
	for i := range parent.Children {
		if parent.Children[i] == child {
			return i
		}
	}
	return -1
}

func (root *AgendaNode) Prev(sigil *AgendaNode) (result *AgendaNode) {
	result = nil
	var lastNode *AgendaNode = nil

	root.Walk(func(node *AgendaNode, _ int) {
		if node == sigil {
			result = lastNode
		}
		lastNode = node
	})

	return
}

func (root *AgendaNode) Next(sigil *AgendaNode) (result *AgendaNode) {
	result = nil
	var lastNode *AgendaNode = nil

	root.Walk(func(node *AgendaNode, _ int) {
		if lastNode == sigil {
			result = node
		}
		lastNode = node
	})

	return
}

// Invoke callback on each node in the tree.
// Walk performs a depth-first, sibling-second traversal of the tree.
// @param callback
//        Takes the node being visited and how many levels deep in the tree the node is.
//
func (node *AgendaNode) Walk(callback func(*AgendaNode, int)) {
	var walk func(*AgendaNode, int, bool)

	walk = func(node *AgendaNode, depth int, processSiblings bool) {
		callback(node, depth)

		for i := range node.Children {
			walk(node.Children[i], depth+1, true)
		}

		if processSiblings {
			for node = node.NextSibling; node != nil; node = node.NextSibling {
				walk(node, depth, false)
			}
		}
	}

	walk(node, 0, true)
}

// Move a node "down".
//   If there are other children:
//     Shuffle the organization so this child gets a higher index .
//   Otherwise:
//     If the parent has a parent and that parent has a higher-index child:
//       Move to parent-parent's next highest index child's children's list at the start.
//     Otherwise:
//       Stop. Do nothing.
//
func (node *AgendaNode) ShuffleDown() {
	parent := node.Parent
	index := parent.IndexChild(node)
	if index == -1 {
		panic(fmt.Errorf("Couldn't find node in children!"))
	}

	if index == len(parent.Children)-1 {
		// Shuffle to next parent.
		if parent.Parent == nil {
			// There is no "super" parent, the parent is the root node.
			return
		}
		parentIndex := parent.Parent.IndexChild(parent)
		if len(parent.Parent.Children) <= parentIndex+1 {
			// There are no more children to pass ownership to!
			return
		}
		// Change ownership to next parent and remove ownership from previous parent.
		newParent := parent.Parent.Children[parentIndex+1]
		newParent.Children = append(newParent.Children, node)
		parent.Children = append(parent.Children[:index], parent.Children[index+1:]...)
	} else {
		// Shuffle index up within children.
		parent.Children[index], parent.Children[index+1] = parent.Children[index+1], parent.Children[index]
	}
}

// Move a node "up".
//   If there are other children:
//     Shuffle the organization so this child gets a lower index.
//   Otherwise:
//     If the parent has a parent and that parent has a lower-index child:
//       Move to parent-parent's next lowest index child's children's list at the end.
//     Otherwise:
//       Stop. Do nothing.
//
func (node *AgendaNode) ShuffleUp() {
	parent := node.Parent
	index := parent.IndexChild(node)
	if index == -1 {
		panic(fmt.Errorf("Couldn't find node in children!"))
	}

	if len(parent.Children) == 1 {
		// Shuffle to previous parent.
		if parent.Parent == nil {
			// There is no "super" parent, the parent is the root node.
			return
		}
		parentIndex := parent.Parent.IndexChild(parent)
		if len(parent.Parent.Children) == 1 {
			// There are no more children to pass ownership to!
			return
		}
		// Change ownership to next parent and remove ownership from previous parent.
		newParent := parent.Parent.Children[parentIndex-1]
		newParent.Children = append([]*AgendaNode{node}, newParent.Children...)
		parent.Children = append(parent.Children[:index], parent.Children[index+1:]...)
	} else {
		// Shuffle index up within children.
		parent.Children[index], parent.Children[index-1] = parent.Children[index-1], parent.Children[index]
	}
}