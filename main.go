package main

/*

TODO:
[ ] Bug: Can't add new child for existing item.
    > Steps:
      - Create item.
      - Edit item.
      - Create new item while editing.
      - Return to main.
    > Expected:
      - New node is a child of edited node.
    > Actual:
      - New node is a child of parent node.

[ ] Feature: Render children somehow in the edit dialog.
[ ] Feature: When adding children in the edit dialog, implicitly create sibling nodes around them.
[ ] Feature: Collapse siblings if children are moved. (Maybe not without undo?)
[ ] Feature: Implement a textarea widget, perhaps built upon gemacs or micro or gomacs or gemacs?
[ ] Feature: Undo.

*/

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os"
	"strings"
)

var (
	log            DebugLog
	boxShown       bool
	rootAgendaNode *AgendaNode
)

func main() {
	var editNode func(*AgendaNode, *AgendaNode)

	newNodeStack := Stack{}
	boxShown = false
	rootAgendaNode := NewAgendaTree()

	mainGrid := tview.NewGrid()

	log.Primitive = tview.NewTextView()
	log.Primitive.SetBorder(true)
	log.Log("Program loaded")

	tree := NewTree(rootAgendaNode)
	tree.SetBorder(true)
	tree.SetTitle("Agenda")
	tree.SetSelectedFunc(func(node *AgendaNode) {
		editNode(nil, node)
	})

	flex := tview.NewFlex()
	flex.SetFullScreen(false)
	flex.SetDirection(tview.FlexColumn)
	flex.AddItem(tree, 0, 1, true)

	help := tview.NewTextView()
	help.SetBorder(true)
	help.SetTitle("help")
	help.SetText(strings.TrimSpace(helpText))

	app := tview.NewApplication()
	inputStack := InputHandlerStack{}
	pageStack := PageStack{}

	pages := tview.NewPages()
	pages.AddPage("main", flex, true, true)
	pages.AddPage("help", help, true, false)
	pageStack.Push(&Page{Name: "main", Primitive: flex})

	mainGrid.SetRows(-1, 3)
	mainGrid.SetColumns(-1)
	mainGrid.AddItem(pages, 0, 0, 1, 1, 1, 1, true)
	mainGrid.AddItem(log.Primitive, 1, 0, 1, 1, 1, 1, false)

	helpWidget := Widget{}
	helpWidget.Primitive = help
	helpWidget.InputHandler = createEscHandler(func() {
		inputStack.Pop()
		pageStack.Pop()
		pages.SwitchToPage(pageStack.Top().Name)
		log.Log("Exiting help, switching to %v", pageStack.Top().Name)
		app.Draw()
	})

	box := tview.NewBox()
	box.SetBorder(true)
	box.SetTitle("A [red]c[yellow]o[green]l[darkcyan]o[blue]r[darkmagenta]f[red]u[yellow]l[white] [black:red]c[:yellow]o[:green]l[:darkcyan]o[:blue]r[:darkmagenta]f[:red]u[:yellow]l[white:] [::bu]title")

	boxWidget := Widget{}
	boxWidget.Primitive = box
	boxWidget.InputHandler = createEscHandler(func() {
		inputStack.Pop()
		boxShown = false
		log.Log("Exiting box, switching to main")
		flex.RemoveItem(box)
		app.Draw()
	})

	flexWidget := Widget{}
	flexWidget.Primitive = flex
	flexWidget.InputHandler = func(event *tcell.EventKey) (result *tcell.EventKey) {
		result = event

		if event.Key() == tcell.KeyRune && event.Rune() == 't' {
			if !boxShown {
				boxWidget.InputHandlerIndex = inputStack.Push(boxWidget.InputHandler)
				flex.AddItem(box, 0, 1, true)
				app.SetFocus(box)
				boxShown = true
				log.Log("Showing box")
				result = nil
			}

			app.Draw()
		}

		return
	}

	editNode = func(scratch *AgendaNode, node *AgendaNode) {
		editNodeWidget := NewEditAgendaNodeWidget(app, node, scratch)
		editNodeWidget.InputHandler = createEscHandler(func() {
			if newNodeStack.Count() <= 1 {
				inputStack.Enable(flexWidget.InputHandlerIndex)
			}

			if node.Parent == nil && node.NextContinuation == nil && node.PrevContinuation == nil {
				if scratch != nil {
					scratch.AddChild(node)
				} else {
					rootAgendaNode.AddChild(node)
					tree.Selected = node
				}
			}
			inputStack.Pop()
			pageStack.Pop()
			newNodeStack.Pop()
			pages.SwitchToPage(pageStack.Top().Name)
			log.Log("Exiting %v, switching to %v", editNodeWidget.Name, pageStack.Top().Name)
			app.Draw()
		})
		editNodeWidget.InputHandlerIndex = inputStack.Push(editNodeWidget.InputHandler)
		inputStack.Disable(flexWidget.InputHandlerIndex)
		pageStack.Push(&Page{Name: editNodeWidget.Name, Primitive: editNodeWidget.Primitive})
		pages.AddPage(editNodeWidget.Name, editNodeWidget.Primitive, true, true)
		pages.SwitchToPage(editNodeWidget.Name)
		log.Log("Showing %s", editNodeWidget.Name)
	}

	pagesWidget := Widget{}
	pagesWidget.Primitive = pages
	pagesWidget.InputHandler = func(event *tcell.EventKey) (result *tcell.EventKey) {
		result = event

		switch event.Key() {
		case tcell.KeyCtrlR:
			app.Draw()
			result = nil

		case tcell.KeyRune:
			switch event.Rune() {
			case '?':
				if pageStack.Top().Name == "help" {
					break
				}

				helpWidget.InputHandlerIndex = inputStack.Push(helpWidget.InputHandler)
				pageStack.Push(&Page{Name: "help", Primitive: helpWidget.Primitive})
				pages.SwitchToPage("help")
				log.Log("Showing help")
				result = nil

			case '+':
				var scratchNode *AgendaNode = nil
				if newNodeStack.Top() != nil {
					scratchNode = newNodeStack.Top().(*AgendaNode)
				}
				node := &AgendaNode{}
				newNodeStack.Push(node)
				editNode(scratchNode, node)
				result = nil
			}
		}
		app.Draw()

		return
	}

	pagesWidget.InputHandlerIndex = inputStack.Push(pagesWidget.InputHandler)
	flexWidget.InputHandlerIndex = inputStack.Push(flexWidget.InputHandler)

	app.SetFocus(pages)
	app.SetInputCapture(createAppInputHandler(&inputStack))
	if err := app.SetRoot(mainGrid, true).Run(); err != nil {
		panic(err)
	}

	rootAgendaNode.PrintTree(os.Stdout, 5)
}

var helpText = `
?           Show this help text.
+           Add a new item.
<ctrl+c>    Quit
<esc>       Quit any popups, dialogs or modals.
<enter>     Edit selected item.
k           Select previous item in list.
j           Select next item in list.
<alt>+h     Outdent the item one level.
<alt>+l     Indent the item one level.
<alt>+k     Move an item up in the list. (Preserves nesting level.)
<alt>+j     Move an item down in the list. (Preserves nesting level.)
`

func NewAgendaTree() *AgendaNode {
	var p *AgendaNode
	var c *AgendaNode
	var index int
	r := NewNode("r", "r r r")

	rc1 := NewNode("rc1", "rc1 rc1 rc1")
	r.AddChild(rc1)
	p = r
	c = rc1
	index = p.IndexChild(c)
	if index == -1 {
		panic("1")
	}

	rc1s1 := NewNode("rc1s1", "rc1s1 rc1s1 rc1s1")
	rc1.AddContinuation(rc1s1)

	rc1s1c1 := NewNode("rc1s1c1", "rc1s1c1 rc1s1c1 rc1s1c1")
	rc1s1.AddChild(rc1s1c1)
	p = rc1s1
	c = rc1s1c1
	index = p.IndexChild(c)
	if index == -1 {
		panic("3")
	}

	rc1s2 := NewNode("rc1s2", "rc1s2 rc1s2 rc1s2")
	rc1s1.AddContinuation(rc1s2)

	rc2 := NewNode("rc2", "rc2 rc2 rc2")
	r.AddChild(rc2)
	p = r
	c = rc2
	index = p.IndexChild(c)
	if index == -1 {
		panic("5")
	}

	rc1s3 := NewNode("rc1s3", "rc1s3 rc1s3 rc1s3")
	rc1s1.AddContinuation(rc1s3)

	rc1s3c1 := NewNode("rc1s3c1", "rc1s3c1 rc1s3c1 rc1s3c1")
	rc1s3.AddContinuation(rc1s3c1)

	return r
}
