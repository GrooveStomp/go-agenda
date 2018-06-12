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
	Parent   *AgendaNode
	Title    string
	Text     string
	Sibling  *AgendaNode
	Children []*AgendaNode
	Tags     []string
}

// func main() {
//   tree := NewNode("tree", "text0 text0 text0 text0", []string{})

//   treec1 := NewNode("treec1", "text1 text1 text1 text1", []string{})
//   tree.AddChild(treec1)

//   treec1s1 := NewNode("treec1s1", "text1b text1b text1b text1b", []string{})
//   treec1.AddSibling(treec1s1)

//   treec1s1c1 := NewNode("treec1s1c1", "text text text text", []string{})
//   treec1s1.AddChild(treec1s1c1)

//   treec1s2 := NewNode("treec1s2", "text1c text1c text1c text1c", []string{})
//   treec1s1.AddSibling(treec1s2)

//   treec2 := NewNode("treec2", "text2 text2 text2 text2", []string{})
//   tree.AddChild(treec2)

//   tree.Walk(func(node *AgendaNode, indentLevel int) {
//     node.Print(os.Stdout, indentLevel, 5)
//   })
// }

func (node *AgendaNode) Print(w io.Writer, indentLevel int, indentScale int) {
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
	for ; node.Sibling != nil; node = node.Sibling {
	}
	node.Sibling = new
	new.Parent = node.Parent
}

func (parent *AgendaNode) Index(child *AgendaNode) int {
	for i := range parent.Children {
		if parent.Children[i] == child {
			return i
		}
	}
	return -1
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
	index := parent.Index(node)
	if index == -1 {
		panic(fmt.Errorf("Couldn't find node in children!"))
	}

	if index == len(parent.Children)-1 {
		// Shuffle to next parent.
		if parent.Parent == nil {
			// There is no "super" parent, the parent is the root node.
			return
		}
		parentIndex := parent.Parent.Index(parent)
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
	index := parent.Index(node)
	if index == -1 {
		panic(fmt.Errorf("Couldn't find node in children!"))
	}

	if len(parent.Children) == 1 {
		// Shuffle to previous parent.
		if parent.Parent == nil {
			// There is no "super" parent, the parent is the root node.
			return
		}
		parentIndex := parent.Parent.Index(parent)
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
