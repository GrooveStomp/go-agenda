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

func (parent *AgendaNode) AddChild(node *AgendaNode) {
	parent.Children = append(parent.Children, node)
	node.Parent = parent
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

func (root *AgendaNode) Prev(subject *AgendaNode) (wanted *AgendaNode) {
	wanted = nil
	var last *AgendaNode = nil

	root.Walk(func(visitee *AgendaNode, _ int) {
		if visitee == subject {
			wanted = last
		}
		last = visitee
	})

	return
}

func (root *AgendaNode) Next(subject *AgendaNode) (wanted *AgendaNode) {
	wanted = nil
	var last *AgendaNode = nil

	root.Walk(func(visitee *AgendaNode, _ int) {
		if last == subject {
			wanted = visitee
		}
		last = visitee
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
	// TODO(AARONO): Debug this!  Doesn't seem to do anything.
	// This only works for non-continuation nodes. ie., must have a parent.
	if subject.Parent == nil {
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
	// TODO(AARONO): Debug this!  Doesn't seem to do anything.
	// This only works for non-continuation nodes. ie., must have a parent.
	if subject.Parent == nil {
		log.Log("No parent, only works for non-continuation nodes.")
		return
	}

	parent := subject.Parent
	index := parent.IndexChild(subject)

	if index == 0 {
		log.Log("node already at first index.")
		return
	}

	rest := append(parent.Children[:index], parent.Children[index+1:]...)
	parent.Children = append([]*AgendaNode{subject}, rest...)
}

// Moves a node up so it becomes a child of its parent's parent.
// aka Outdent.
//
func (subject *AgendaNode) MoveUpTree() {
	if subject.Parent == nil {
		log.Log("Couldn't find parent")
		return
	}
	curParent := subject.Parent
	canonicalParent := curParent

	index := curParent.IndexChild(subject)
	if index == -1 {
		log.Log("Couldn't find index of node in parent")
		return
	}

	if curParent.Parent == nil {
		tmpParent := curParent
		for ; tmpParent.PrevContinuation != nil; tmpParent = tmpParent.PrevContinuation {
		}
		for ; tmpParent.NextContinuation != nil; tmpParent = tmpParent.NextContinuation {
			if tmpParent.Parent != nil {
				canonicalParent = tmpParent
				break
			}
		}
	}
	if canonicalParent.Parent == nil {
		log.Log("Couldn't find parent's parent")
		return
	}

	curParent.Children = append(curParent.Children[:index], curParent.Children[index+1:]...)
	parent := canonicalParent.Parent
	parent.Children = append([]*AgendaNode{subject}, parent.Children...)
	subject.Parent = parent
}

// Makes subject a child of its closest sibling.
// aka Indent.
//
func (subject *AgendaNode) MoveDownTree() {
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
