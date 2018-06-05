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

	"code.groovestomp.com/agenda/internal/hndlstack"
)

type Widget struct {
	Primitive tview.Primitive
	InputHandler hndlstack.InputHandler
}

type AgendaItem struct {
	Title  string
	Bodies []string // We can have N bodies separated by N-1 Links to other agenda items.
	Links  []AgendaItem
	Tags   []string
}

func createEscHandler(callback func()) hndlstack.InputHandler {
	return func(eventKey *tcell.EventKey) *tcell.EventKey {
		if eventKey.Key() == tcell.KeyEsc {
			callback()
			return nil
		}
		return eventKey
	}
}

var debugOut *tview.TextView

func main() {
	box := tview.NewBox()
	box.SetBorder(true)
	box.SetTitle("A [red]c[yellow]o[green]l[darkcyan]o[blue]r[darkmagenta]f[red]u[yellow]l[white] [black:red]c[:yellow]o[:green]l[:darkcyan]o[:blue]r[:darkmagenta]f[:red]u[:yellow]l[white:] [::bu]title")

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

	flex := tview.NewFlex()
	flex.SetFullScreen(false)
	flex.SetDirection(tview.FlexColumn)
	flex.AddItem(list, 0, 1, true)

	pages := tview.NewPages()
	pages.AddPage("main", flex, true, true)
	pages.AddPage("help", help, true, false)

	mainGrid.SetRows(-4, 3)
	mainGrid.SetColumns(-1)
	mainGrid.AddItem(pages, 0, 0, 1, 1, 1, 1, true)
	mainGrid.AddItem(debugOut, 1, 0, 1, 1, 1, 1, false)

	boxShown := false
	helpShown := false
	currentPage := "main"
	lastPage := "main"
	addItemCount := 0

	app := tview.NewApplication()
	handlers := hndlstack.InputHandlerStack{}

	helpWidget := Widget{}
	helpWidget.Primitive = help
	helpWidget.InputHandler = createEscHandler(func() {
		currentPage = lastPage
		lastPage = "help"
		helpShown = false
		debugOut.SetText("Exiting Help")
		pages.SwitchToPage(currentPage)
		handlers.Pop()
	})

	boxWidget := Widget{}
	boxWidget.Primitive = box
	boxWidget.InputHandler = createEscHandler(func() {
		debugOut.SetText("Exiting Box")
		flex.RemoveItem(box)
		boxShown = false
		debugOut.SetText("Box Exited")
		handlers.Pop()
	})

	flexWidget := Widget{}
	flexWidget.Primitive = flex
	flexWidget.InputHandler = func(event *tcell.EventKey) (result *tcell.EventKey) {
		result = event
		if event.Key() != tcell.KeyRune {
			return
		}

		if event.Rune() == 't' {
			if currentPage != "main" {
				return
			}

			if !boxShown {
				handlers.Push(boxWidget.InputHandler)
				flex.AddItem(box, 0, 1, true)
				app.SetFocus(box)
				boxShown = true
				result = nil
				debugOut.SetText("Showing Box")
			}
		}

		app.Draw()

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
			name := addItemPage(pages, addItemCount)
			lastPage = currentPage
			currentPage = name
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
}

func addItemPage(pages *tview.Pages, dialogNum int) string {
	/*
		This should have its own input handler to:
		- tab between elements.
		- record data in a common datastructure.
	*/
	name := fmt.Sprintf("addAgendaItem%v", dialogNum)

	grid := tview.NewGrid()
	grid.SetRows(3, -1, 3)
	grid.SetColumns(-1)

	title := tview.NewInputField()
	title.SetBorder(true)
	title.SetTitle("List String")

	body := tview.NewInputField()
	body.SetBorder(true)
	body.SetTitle("Full Description")

	tags := tview.NewInputField()
	tags.SetBorder(true)
	tags.SetTitle("Tags")
	tags.SetPlaceholder("tag1 tag-2 spaces_separate_tags")

	grid.AddItem(title, 0, 0, 1, 1, 1, 1, true)
	grid.AddItem(body, 1, 0, 1, 1, 1, 1, false)
	grid.AddItem(tags, 2, 0, 1, 1, 1, 1, false)

	// NOTE(AARONO): This input handler doesn't work when attached at this level.
	grid.SetInputCapture(createEscHandler(func() {
		debugOut.SetText(fmt.Sprintf("Exiting %v", name))
		pages.RemovePage(name)
	}))

	pages.AddPage(name, grid, true, false)
	pages.SwitchToPage(name)

	return name
}

func createAppInputHandler(stack *hndlstack.InputHandlerStack) hndlstack.InputHandler {
	return func(event *tcell.EventKey) *tcell.EventKey {
		result := event
		for i := len(stack.InputHandler) - 1; i >= 0; i-- {
			handler := stack.InputHandler[i]
			res := handler(event)
			if res == nil {
				return nil
			}
		}
		return result
	}
}
