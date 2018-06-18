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

text 1-1, text 1-2 and text 1-3 are all Continuations of each other.
text 1-2 and text 1-3 should not have headings of their own.

text 1a-1 and text 1b-1 are children of text 1-1 (but *not* 1-2 or 1-3), and we
call them siblings of each other.

+-------------------------------------------------------------------------------
| Operations:
| Next Parent
|   Given a node in the tree, find the next suitable parent which would occur
|   "down" in the list as rendered visually via PrintTree.
| Prev Parent
|   Given a node in the tree, find the previous suitable parent which would
|   occur "up" in the list as rendered visually via PrintTree.
| Outdent
|   Move the given node up one level in the tree so it becomes a child of its
|   current parent's parent.
| Indent
|   Move the given node down one level in the tree so it becomes a child of the
|   current node.
+-------------------------------------------------------------------------------
*/

import (
	"fmt"
	"io"
)

type AgendaNode struct {
	Parent           *AgendaNode
	Title            string
	Text             string
	NextContinuation *AgendaNode
	PrevContinuation *AgendaNode
	Children         []*AgendaNode
	Tags             []string
}

// func main() {
// r := NewNode("r", "r r r", []string{})

// rc1 := NewNode("rc1", "rc1 rc1 rc1", []string{})
// r.AddChild(rc1)

// rc1s1 := NewNode("rc1s1", "rc1s1 rc1s1 rc1s1", []string{})
// rc1.AddContinuation(rc1s1)

// rc1s1c1 := NewNode("rc1s1c1", "rc1s1c1 rc1s1c1 rc1s1c1", []string{})
// rc1s1.AddChild(rc1s1c1)

// rc1s2 := NewNode("rc1s2", "rc1s2 rc1s2 rc1s2", []string{})
// rc1s1.AddContinuation(rc1s2)

// rc2 := NewNode("rc2", "rc2 rc2 rc2", []string{})
// r.AddChild(rc2)

// rc1s3 := NewNode("rc1s3", "rc1s3 rc1s3 rc1s3", []string{})
// rc1s1.AddContinuation(rc1s3)

// rc1s3c1 := NewNode("rc1s3c1", "rc1s3c1 rc1s3c1 rc1s3c1", []string{})
// rc1s3c1.AddContinuation(rc1s3c1)

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

func NewNode(title, text string, tags ...string) *AgendaNode {
	new := &AgendaNode{Title: title, Text: text, Tags: tags}
	return new
}

func (parent *AgendaNode) AddChild(child *AgendaNode) {
	parent.Children = append(parent.Children, child)
	child.Parent = parent
}

func (parent *AgendaNode) InsertChild(child *AgendaNode, index int) error {
	if index < 0 || index > len(parent.Children) {
		return fmt.Errorf("index is out of range")
	}

	switch index {
	case 0:
		parent.Children = append([]*AgendaNode{child}, parent.Children...)
	case len(parent.Children):
		parent.Children = append(parent.Children, child)
	default:
		rest := parent.Children[index:]

		str := ""
		for i := range parent.Children[:index] {
			str = fmt.Sprintf("%s %s", str, parent.Children[i].Title)
		}
		str = fmt.Sprintf("%s %s", str, child.Title)
		for i := range rest {
			str = fmt.Sprintf("%s %s", str, rest[i].Title)
		}

		log.Log("%s", str)

		new := append([]*AgendaNode{child}, parent.Children[index:]...)
		parent.Children = append(parent.Children[:index], new...)
	}

	child.Parent = parent

	return nil
}

func (parent *AgendaNode) RemoveChild(child *AgendaNode) {
	index := parent.IndexChild(child)

	if index == len(parent.Children)-1 {
		parent.Children = parent.Children[:index]
	} else {
		parent.Children = append(parent.Children[:index], parent.Children[index+1:]...)
	}

	child.Parent = nil
}

func (node *AgendaNode) AddContinuation(new *AgendaNode) {
	for ; node.NextContinuation != nil; node = node.NextContinuation {
	}
	node.NextContinuation = new
	new.PrevContinuation = node
}

func (parent *AgendaNode) IndexChild(child *AgendaNode) int {
	for i := range parent.Children {
		if parent.Children[i] == child {
			return i
		}
	}
	return -1
}

// Returns parent, or a continuation of parent where that new node has a valid parent.
func (subject *AgendaNode) ParentContinuationWithParent() (*AgendaNode, error) {
	node := subject.Parent
	if node == nil {
		return nil, fmt.Errorf("Parent is nil")
	}

	if node.Parent != nil {
		return node, nil
	}

	// Traverse to the first node in the chain of continuations.
	for ; node.PrevContinuation != nil; node = node.PrevContinuation {
	}

	if node.Parent == nil {
		return nil, fmt.Errorf("Couldn't find grandparent.")
	}

	return node, nil
}

func (node *AgendaNode) IsContinuation() bool {
	return node.Parent == nil
}

func (root *AgendaNode) Prev(subject *AgendaNode) (wanted *AgendaNode) {
	wanted = nil
	var last *AgendaNode = nil

	root.Walk(func(visitee *AgendaNode, _ int) {
		if visitee == subject {
			wanted = last
		}
		if !visitee.IsContinuation() {
			last = visitee
		}
	})

	return
}

func (root *AgendaNode) Next(subject *AgendaNode) (wanted *AgendaNode) {
	wanted = nil
	var last *AgendaNode = nil

	root.Walk(func(visitee *AgendaNode, _ int) {
		if last == subject && !visitee.IsContinuation() {
			wanted = visitee
		}
		if !visitee.IsContinuation() {
			last = visitee
		}
	})

	return
}

// Invoke callback on each node in the tree.
// Walk performs a depth-first, continuation-second traversal of the tree.
// @param callback
//        Takes the node being visited and how many levels deep in the tree the node is.
//
func (node *AgendaNode) Walk(callback func(*AgendaNode, int)) {
	var walk func(*AgendaNode, int)

	walk = func(node *AgendaNode, depth int) {
		callback(node, depth)

		for i := range node.Children {
			walk(node.Children[i], depth+1)
		}

		if node.NextContinuation != nil {
			walk(node.NextContinuation, depth)
		}
	}

	// NOTE(AARONO): Skip the root node.
	for i := range node.Children {
		walk(node.Children[i], 0)
	}
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
func (subject *AgendaNode) MakeNextSibling() {
	if subject.IsContinuation() {
		log.Log("No parent, only works for non-continuation nodes.")
		return
	}

	parent := subject.Parent
	index := parent.IndexChild(subject)
	count := len(parent.Children)

	if index == count-1 {
		log.Log("node already at last index.")
		return
	}

	rest := parent.Children[index+1:]
	parent.Children = parent.Children[:index]
	parent.Children = append(parent.Children, rest[0])
	parent.Children = append(parent.Children, subject)
	parent.Children = append(parent.Children, rest[1:]...)
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
func (subject *AgendaNode) MakePrevSibling() {
	// This only works for non-continuation nodes. ie., must have a parent.
	if subject.IsContinuation() {
		log.Log("No parent, only works for non-continuation nodes.")
		return
	}

	parent := subject.Parent
	index := parent.IndexChild(subject)

	if index == 0 {
		log.Log("node already at first index.")
		return
	}

	a, b := parent.Children[index-1], parent.Children[index]
	parent.Children[index-1], parent.Children[index] = b, a
}

// Makes subject the next sibling of its current parent.
// aka Outdent.
//
func (subject *AgendaNode) MoveUpTree() {
	parentSib, err := subject.ParentContinuationWithParent()
	if err != nil {
		log.Log("Couldn't find parent's parent")
		return
	}

	subject.Parent.RemoveChild(subject)
	newParent := parentSib.Parent
	index := newParent.IndexChild(parentSib)
	if index == -1 {
		panic("Oh noes!")
	}

	err = newParent.InsertChild(subject, index+1)
	if err != nil {
		panic(err)
	}
}

// Makes subject a child of its previous sibling.
// aka Indent.
//
func (subject *AgendaNode) MoveDownTree() {
	if subject.IsContinuation() {
		log.Log("Node must have a parent.")
		return
	}
	parent := subject.Parent
	index := parent.IndexChild(subject)
	if len(parent.Children) == 1 {
		log.Log("Must have at least one sibling.")
		return
	}
	if index == 0 {
		log.Log("Can't indent from here. Try moving up first.")
		return
	}

	parent.Children = append(parent.Children[:index], parent.Children[index+1:]...)
	newParent := parent.Children[index-1] // Previous child before subject.
	for ; newParent.NextContinuation != nil; newParent = newParent.NextContinuation {
	}
	newParent.AddChild(subject)
}

// Swap two nodes in thre tree so that left appears where right was, and right
// appears where left was.
//
func (tree *AgendaNode) swap(left *AgendaNode, right *AgendaNode) {
	tmp := &AgendaNode{}
	left.replaceWith(tmp)
	right.replaceWith(left)
	tmp.replaceWith(right)
}

// Changes tree pointers so dst appears in the place where src was.
//
func (src *AgendaNode) replaceWith(dst *AgendaNode) {
	dst.Parent = src.Parent
	dst.NextContinuation = src.NextContinuation
	dst.PrevContinuation = src.PrevContinuation

	if dst.NextContinuation != nil {
		dst.NextContinuation.PrevContinuation = dst
	}

	if dst.PrevContinuation != nil {
		dst.PrevContinuation.NextContinuation = dst
	}

	if dst.Parent != nil {
		index := dst.Parent.IndexChild(src)
		if index == -1 {
			panic("Move! Parent!")
		} else {
			dst.Parent.Children[index] = dst
		}
	}
}
