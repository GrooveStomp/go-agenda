package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strings"
)

func main() {
	box := tview.NewBox()
	box.SetBorder(true)
	box.SetTitle("A [red]c[yellow]o[green]l[darkcyan]o[blue]r[darkmagenta]f[red]u[yellow]l[white] [black:red]c[:yellow]o[:green]l[:darkcyan]o[:blue]r[:darkmagenta]f[:red]u[:yellow]l[white:] [::bu]title")

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
	flex.SetFullScreen(true)
	flex.SetDirection(tview.FlexColumn)
	flex.AddItem(list, 0, 1, true)

	pages := tview.NewPages()
	pages.AddPage("main", flex, true, true)
	pages.AddPage("help", help, true, false)

	boxShown := false
	helpShown := false
	currentPage := "main"
	lastPage := "main"

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
			
			if boxShown {
				flex.RemoveItem(box)
				boxShown = false
			} else {
				flex.AddItem(box, 0, 1, true)
				boxShown = true
			}
			result = nil
		} else if event.Rune() == '?' {
			if helpShown {
				tmp := currentPage
				currentPage = lastPage
				lastPage = tmp
				pages.SwitchToPage(currentPage)
				helpShown = false
			} else {
				lastPage = currentPage
				currentPage = "help"
				pages.SwitchToPage(currentPage)
				helpShown = true
			}
		}
		
		app.Draw()

		return
	})

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
