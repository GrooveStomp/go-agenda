package main

/*

TODO:
- Implement a textarea widget, perhaps built upon gemacs or micro or gomacs?
- Prototype and experiment with fluidity and mechanics of building nestable, hierarchical lists in the UI.

Org-Mode is a hierarchical outline.
It consists of headlines and collapsing portions of the document.

When I use org-mode, there are a few specific features I make use of:
- title
- tags
- description

Part of what makes org-mode good to use is that you can quickly and arbitrarily
define nested sections.  Adding a new nested section should be quick and
painless, and moving sections around (indenting, outdenting, removing as a
headline (collapsing to parent's level) should also feel natural.

*/

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os"
	"strings"
)

var (
	debugOut     *tview.TextView
	boxShown     bool
	addItemCount int
	flexWidget   Widget
)

func main() {
	boxShown = false
	addItemCount = 0
	rootAgendaNode := NewNode("", "", []string{})

	mainGrid := tview.NewGrid()

	debugOut = tview.NewTextView()
	debugOut.SetTitle("Debug Output")
	debugOut.SetBorder(true)
	debugOut.SetText("Hello Daniela")

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Agenda")
	list.AddItem("Example 1", "Description 1", 0, nil)
	list.AddItem("Example 2", "Description 2", 0, nil)
	list.AddItem("Example 3", "Description 3", 0, nil)
	list.AddItem("Example 4", "Description 4", 0, nil)
	list.AddItem("Example 5", "Description 5", 0, nil)

	flex := tview.NewFlex()
	flex.SetFullScreen(false)
	flex.SetDirection(tview.FlexColumn)
	flex.AddItem(list, 0, 1, true)

	help := tview.NewTextView()
	help.SetBorder(true)
	help.SetTitle("help")
	help.SetText(strings.TrimSpace(`
?        Show this help text.
+        Add a new item.
x        Mark an item as complete.
<tab>    Expand an item.
<ctrl+c> Quit
<esc>    Quit any popups, dialogs or modals.
`))

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
	mainGrid.AddItem(debugOut, 1, 0, 1, 1, 1, 1, false)

	helpWidget := Widget{}
	helpWidget.Primitive = help
	helpWidget.InputHandler = createEscHandler(func() {
		inputStack.Pop()
		pageStack.Pop()
		pages.SwitchToPage(pageStack.Top().Name)
		debugOut.SetText(fmt.Sprintf("Exiting help, switching to %v", pageStack.Top().Name))
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
		debugOut.SetText("Exiting Box, switching to main")
		flex.RemoveItem(box)
		app.Draw()
	})

	flexWidget = Widget{}
	flexWidget.Primitive = flex
	flexWidget.InputHandler = func(event *tcell.EventKey) (result *tcell.EventKey) {
		result = event

		if event.Key() == tcell.KeyRune && event.Rune() == 't' {
			if !boxShown {
				boxWidget.InputHandlerIndex = inputStack.Push(boxWidget.InputHandler)
				flex.AddItem(box, 0, 1, true)
				app.SetFocus(box)
				boxShown = true
				debugOut.SetText("Showing Box")
				result = nil
			}

			app.Draw()
		}

		return
	}

	pagesWidget := Widget{}
	pagesWidget.Primitive = pages
	pagesWidget.InputHandler = func(event *tcell.EventKey) (result *tcell.EventKey) {
		result = event
		if event.Key() != tcell.KeyRune {
			return
		}

		switch event.Rune() {
		case '?':
			if pageStack.Top().Name == "help" {
				break
			}

			helpWidget.InputHandlerIndex = inputStack.Push(helpWidget.InputHandler)
			pageStack.Push(&Page{Name: "help", Primitive: helpWidget.Primitive})
			pages.SwitchToPage("help")
			debugOut.SetText("Showing Help")
			result = nil

		case '+':
			addItemCount += 1
			name, primitive, inputHandler := addItemPage(app, &inputStack, &pageStack, pages, addItemCount, rootAgendaNode)
			inputStack.Push(inputHandler)
			pageStack.Push(&Page{Name: name, Primitive: primitive})
			pages.SwitchToPage(name)
			debugOut.SetText(fmt.Sprintf("Showing %s", name))
			result = nil
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

	rootAgendaNode.Walk(func(node *AgendaNode, indentLevel int) {
		node.Print(os.Stdout, indentLevel, 5)
	})
}

func addItemPage(app *tview.Application, inputStack *InputHandlerStack, pageStack *PageStack, pages *tview.Pages, dialogNum int, tree *AgendaNode) (string, tview.Primitive, InputHandler) {
	/*
		This should have its own input handler to:
		- tab between elements.
		- record data in a common datastructure.
	*/
	node := NewNode("", "", []string{})
	name := fmt.Sprintf("addAgendaItem%v", dialogNum)

	grid := tview.NewGrid()
	title := tview.NewInputField()
	body := tview.NewInputField()

	inputStack.Disable(flexWidget.InputHandlerIndex)

	grid.SetRows(3, -1, 3)
	grid.SetColumns(-1)

	handleEsc := createEscHandler(func() {
		inputStack.Enable(flexWidget.InputHandlerIndex)
		tree.AddChild(node)
		inputStack.Pop()
		pageStack.Pop()
		pages.SwitchToPage(pageStack.Top().Name)
		debugOut.SetText(fmt.Sprintf("Exiting %v, switching to %v", name, pageStack.Top().Name))
		app.Draw()
	})

	title.SetBorder(true)
	title.SetTitle("List String")
	title.SetDoneFunc(func(key tcell.Key) {
		node.Title = title.GetText()
		switch key {
		case tcell.KeyEnter:
			app.SetFocus(body)
		case tcell.KeyTab:
			app.SetFocus(body)
		case tcell.KeyEsc:
			handleEsc(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
		case tcell.KeyBacktab:
			debugOut.SetText("<backtab>")
		default:
		}
	})
	title.SetChangedFunc(func(text string) {
		debugOut.SetText(text)
		app.Draw()
	})

	body.SetBorder(true)
	body.SetTitle("Full Description")
	body.SetDoneFunc(func(key tcell.Key) {
		node.Text = body.GetText()
		switch key {
		case tcell.KeyEnter:
			handleEsc(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
		case tcell.KeyTab:
			debugOut.SetText("<tab>")
		case tcell.KeyEsc:
			handleEsc(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
		case tcell.KeyBacktab:
			app.SetFocus(title)
		default:
		}
	})
	body.SetChangedFunc(func(text string) {
		debugOut.SetText(text)
		app.Draw()
	})

	tags := tview.NewInputField()
	tags.SetBorder(true)
	tags.SetTitle("Tags")
	tags.SetPlaceholder("tag1 tag-2 spaces_separate_tags")
	tags.SetDoneFunc(func(tcell.Key) {
		node.Tags = []string{tags.GetText()} // TODO(AARONO): Split text on spaces!
	})

	grid.AddItem(title, 0, 0, 1, 1, 1, 1, true)
	grid.AddItem(body, 1, 0, 1, 1, 1, 1, false)
	grid.AddItem(tags, 2, 0, 1, 1, 1, 1, false)

	pages.AddPage(name, grid, true, true)

	return name, grid, handleEsc
}
