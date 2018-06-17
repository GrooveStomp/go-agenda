package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Tree struct {
	*tview.Box
	Root         *AgendaNode
	Indent       int
	Selected     *AgendaNode
	selectedFunc func(*AgendaNode)
}

func NewTree(root *AgendaNode) *Tree {
	return &Tree{
		Box:      tview.NewBox(),
		Root:     root,
		Indent:   5,
		Selected: nil,
	}
}

func (t *Tree) Draw(screen tcell.Screen) {
	t.Box.Draw(screen)
	x, y, width, _ /*height*/ := t.GetInnerRect()

	t.Root.Walk(func(node *AgendaNode, indentLevel int) {
		tview.Print(screen, node.Title, x+(indentLevel*t.Indent), y, width, tview.AlignLeft, tview.Styles.PrimaryTextColor)
		if t.Selected == node {
			textWidth := len(node.Title)
			for bx := 0; bx < textWidth; bx++ {
				m, c, style, _ := screen.GetContent(x+bx, y)
				fg, _, _ := style.Decompose()
				if fg == tview.Styles.PrimaryTextColor {
					fg = tview.Styles.PrimitiveBackgroundColor
				}
				style = style.Background(tview.Styles.PrimaryTextColor).Foreground(fg)
				screen.SetContent(x+(indentLevel*t.Indent)+bx, y, m, c, style)
			}
		}
		y++
		tview.Print(screen, node.Text, x+(indentLevel*t.Indent), y, width, tview.AlignLeft, tview.Styles.TertiaryTextColor)
		y++
	})
}

func (t *Tree) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return t.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyDown:
			next := t.Root.Next(t.Selected)
			if next != nil {
				t.Selected = next
			}

		case tcell.KeyUp:
			previous := t.Root.Prev(t.Selected)
			if previous != nil {
				t.Selected = previous
			}

		case tcell.KeyEnter:
			t.selectedFunc(t.Selected)

		case tcell.KeyCtrlU:
			t.Root.ShuffleUp(t.Selected)

		case tcell.KeyCtrlD:
			t.Root.ShuffleDown(t.Selected)

		case tcell.KeyLeft:
			t.Selected.MoveUpTree()

		case tcell.KeyRight:
			t.Selected.MoveDownTree()

		default:
		}
	})
}

func (t *Tree) SetSelectedFunc(callback func(node *AgendaNode)) {
	t.selectedFunc = callback
}

// func main() {
// 	app := tview.NewApplication()

// 	tree := NewTree(NewAgendaTree())
// 	tree.SetBorder(true)
// 	tree.SetTitle("Tree")

// 	if err := app.SetRoot(tree, true).Run(); err != nil {
// 		panic(err)
// 	}

// 	NewAgendaTree().PrintTree(os.Stdout, 5)
// }

// func NewAgendaTree() *AgendaNode {
// 	r := NewNode("r", "r r r", []string{})

// 	rc1 := NewNode("rc1", "rc1 rc1 rc1", []string{})
// 	r.AddChild(rc1)

// 	rc1s1 := NewNode("rc1s1", "rc1s1 rc1s1 rc1s1", []string{})
// 	rc1.AddSibling(rc1s1)

// 	rc1s1c1 := NewNode("rc1s1c1", "rc1s1c1 rc1s1c1 rc1s1c1", []string{})
// 	rc1s1.AddChild(rc1s1c1)

// 	rc1s2 := NewNode("rc1s2", "rc1s2 rc1s2 rc1s2", []string{})
// 	rc1s1.AddSibling(rc1s2)

// 	rc2 := NewNode("rc2", "rc2 rc2 rc2", []string{})
// 	r.AddChild(rc2)

// 	rc1s3 := NewNode("rc1s3", "rc1s3 rc1s3 rc1s3", []string{})
// 	rc1s1.AddSibling(rc1s3)

// 	rc1s3c1 := NewNode("rc1s3c1", "rc1s3c1 rc1s3c1 rc1s3c1", []string{})
// 	rc1s3.AddSibling(rc1s3c1)

// 	return r
// }
