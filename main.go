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
	"strings"
)

var (
	debugOut     *tview.TextView
	boxShown     bool
	helpShown    bool
	currentPage  string
	lastPage     string
	addItemCount int
)

func main() {
	boxShown = false
	helpShown = false
	currentPage = "main"
	lastPage = "main"
	addItemCount = 0
	root := NewNode("", "", []string{})

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

	pages := tview.NewPages()
	pages.AddPage("main", flex, true, true)
	pages.AddPage("help", help, true, false)

	mainGrid.SetRows(-1, 3)
	mainGrid.SetColumns(-1)
	mainGrid.AddItem(pages, 0, 0, 1, 1, 1, 1, true)
	mainGrid.AddItem(debugOut, 1, 0, 1, 1, 1, 1, false)

	app := tview.NewApplication()
	handlers := InputHandlerStack{}

	helpWidget := Widget{}
	helpWidget.Primitive = help
	helpWidget.InputHandler = createEscHandler(func() {
		currentPage = lastPage
		lastPage = "help"
		helpShown = false
		debugOut.SetText("Exiting Help")
		pages.SwitchToPage(currentPage)
		handlers.Pop()
		app.Draw()
	})

	box := tview.NewBox()
	box.SetBorder(true)
	box.SetTitle("A [red]c[yellow]o[green]l[darkcyan]o[blue]r[darkmagenta]f[red]u[yellow]l[white] [black:red]c[:yellow]o[:green]l[:darkcyan]o[:blue]r[:darkmagenta]f[:red]u[:yellow]l[white:] [::bu]title")

	boxWidget := Widget{}
	boxWidget.Primitive = box
	boxWidget.InputHandler = createEscHandler(func() {
		debugOut.SetText("Exiting Box")
		flex.RemoveItem(box)
		boxShown = false
		debugOut.SetText("Box Exited")
		handlers.Pop()
		app.Draw()
	})

	flexWidget := Widget{}
	flexWidget.Primitive = flex
	flexWidget.InputHandler = func(event *tcell.EventKey) (result *tcell.EventKey) {
		result = event

		if event.Key() == tcell.KeyRune && event.Rune() == 't' {
			if !boxShown {
				handlers.Push(boxWidget.InputHandler)
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

		if event.Rune() == '?' && !helpShown {
			handlers.Push(helpWidget.InputHandler)
			lastPage = currentPage
			currentPage = "help"
			pages.SwitchToPage(currentPage)
			helpShown = true
			result = nil
			debugOut.SetText("Showing Help")
		} else if event.Rune() == '+' {
			addItemCount += 1
			// TODO: Return a new node agenda node.
			name := addItemPage(app, &handlers, pages, addItemCount, root)
			lastPage = currentPage
			currentPage = name
			pages.SwitchToPage(currentPage)
			result = nil
			debugOut.SetText(fmt.Sprintf("Showing %s", name))
		}

		app.Draw()

		return
	}

	handlers.Push(pagesWidget.InputHandler)
	handlers.Push(flexWidget.InputHandler)

	app.SetFocus(pages)
	app.SetInputCapture(createAppInputHandler(&handlers))
	if err := app.SetRoot(mainGrid, true).Run(); err != nil {
		panic(err)
	}

	print := func(node *AgendaNode, depth int) {
		fmt.Printf("%*s%v\n", depth*5, " ", node.Title)
		fmt.Printf("%*s%v\n", depth*5, " ", node.Text)
	}

	root.Walk(print)
}

func addItemPage(app *tview.Application, handlers *InputHandlerStack, pages *tview.Pages, dialogNum int, tree *AgendaNode) string {
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

	grid.SetRows(3, -1, 3)
	grid.SetColumns(-1)

	handleEsc := createEscHandler(func() {
		tree.AddChild(node)
		currentPage = lastPage
		lastPage = name
		debugOut.SetText(fmt.Sprintf("Exiting %v", name))
		pages.SwitchToPage(currentPage)
		handlers.Pop()
		app.Draw()
	})

	title.SetBorder(true)
	title.SetTitle("List String")
	title.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			node.Title = title.GetText()
			app.SetFocus(body)
		case tcell.KeyTab:
			app.SetFocus(body)
		case tcell.KeyEsc:
			node.Title = title.GetText()
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
		switch key {
		case tcell.KeyEnter:
			node.Text = body.GetText()
			handleEsc(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
		case tcell.KeyTab:
			debugOut.SetText("<tab>")
		case tcell.KeyEsc:
			node.Text = body.GetText()
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

	handlers.Push(handleEsc)
	pages.AddPage(name, grid, true, true)

	return name
}
