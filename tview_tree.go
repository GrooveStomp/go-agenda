package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Tree struct {
	*tview.Box
	Root     *AgendaNode
	Indent   int
	Selected *AgendaNode
}

func NewTree(root *AgendaNode) *Tree {
	return &Tree{
		Box:      tview.NewBox(),
		Root:     root,
		Indent:   5,
		Selected: root,
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

		default:
		}
	})
}

// var (
// 	debug *tview.TextView
// )

// func main() {
// 	app := tview.NewApplication()

// 	debug = tview.NewTextView()
// 	debug.SetBorder(true)

// 	grid := tview.NewGrid()
// 	grid.SetRows(-1, 3)
// 	grid.SetColumns(-1)

// 	tree := NewTree(NewAgendaTree())
// 	tree.SetBorder(true)
// 	tree.SetTitle("Tree")

// 	grid.AddItem(tree, 0, 0, 1, 1, 1, 1, true)
// 	grid.AddItem(debug, 1, 0, 1, 1, 1, 1, false)

// 	if err := app.SetRoot(grid, true).Run(); err != nil {
// 		panic(err)
// 	}

// 	NewAgendaTree().PrintTree(os.Stdout, 5)
// }
