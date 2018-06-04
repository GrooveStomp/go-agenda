package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strings"
)

type inputHandler func(*tcell.EventKey) *tcell.EventKey

func createEscHandler(callback func()) inputHandler {
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

	help.SetInputCapture(createEscHandler(func() {
		currentPage = lastPage
		lastPage = "help"
		helpShown = false
		debugOut.SetText("Exiting Help")
		pages.SwitchToPage(currentPage)
	}))

	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) (result *tcell.EventKey) {
		result = event
		if event.Key() != tcell.KeyRune {
			return
		}

		if event.Rune() == 't' {
			if currentPage != "main" {
				return
			}

			if !boxShown {
				flex.AddItem(box, 0, 1, true)
				app.SetFocus(box)
				boxShown = true
				result = nil
				debugOut.SetText("Showing Box")
			}
		} else if event.Rune() == '?' {
			lastPage = currentPage
			currentPage = "help"
			pages.SwitchToPage(currentPage)
			helpShown = true
			result = nil
			debugOut.SetText("Showing Help")
		}

		app.Draw()

		return
	})

	// box.SetInputCapture(func(eventKey *tcell.EventKey) *tcell.EventKey {
	// 	if eventKey.Key() == tcell.KeyRune && eventKey.Rune() == 't' {
	// 		flex.RemoveItem(box)
	// 		boxShown = false
	// 		debugOut.SetText("Exiting Box")
	// 		return nil
	// 	}
	// 	return eventKey
	// })
	box.SetInputCapture(createEscHandler(func() {
		flex.RemoveItem(box)
		boxShown = false
		debugOut.SetText("Exiting Box")
	}))

	if err := app.SetRoot(mainGrid, true).Run(); err != nil {
		panic(err)
	}
}
