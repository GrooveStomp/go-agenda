package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
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

	flex := tview.NewFlex()
	flex.SetFullScreen(true)
	flex.SetDirection(tview.FlexColumn)
	flex.AddItem(list, 0, 1, true)

	boxShown := false

	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyRune {
			return event
		}

		if event.Rune() == 't' {
			if boxShown {
				flex.RemoveItem(box)
				boxShown = false
			} else {
				flex.AddItem(box, 0, 1, true)
				boxShown = true
			}
			app.Draw()
			return nil
		}
		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
