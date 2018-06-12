package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	NodeNum = 0
)

func NewEditAgendaNodeWidget(app *tview.Application, node *AgendaNode) (widget *Widget) {
	widget = &Widget{}

	title := tview.NewInputField()
	body := tview.NewInputField()

	title.SetBorder(true)
	title.SetTitle("Title")
	title.SetDoneFunc(func(key tcell.Key) {
		node.Title = title.GetText()
		switch key {
		case tcell.KeyEnter:
			log.Log("<enter>")
			app.SetFocus(body)
		case tcell.KeyTab:
			log.Log("<tab>")
			app.SetFocus(body)
		case tcell.KeyEsc:
			log.Log("<esc>")
			widget.InputHandler(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
		case tcell.KeyBacktab:
			log.Log("<backtab>")
		default:
		}
	})
	title.SetChangedFunc(func(text string) {
		log.Log(text)
		app.Draw()
	})

	body.SetBorder(true)
	body.SetTitle("Body")
	body.SetDoneFunc(func(key tcell.Key) {
		node.Text = body.GetText()
		switch key {
		case tcell.KeyEnter:
			widget.InputHandler(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
			log.Log("<enter>")
		case tcell.KeyTab:
			log.Log("<tab>")
		case tcell.KeyEsc:
			log.Log("<esc>")
			widget.InputHandler(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
		case tcell.KeyBacktab:
			log.Log("<backtab>")
			app.SetFocus(title)
		default:
		}
	})
	body.SetChangedFunc(func(text string) {
		log.Log(text)
		app.Draw()
	})

	grid := tview.NewGrid()
	grid.SetRows(3, -1, 3)
	grid.SetColumns(-1)

	grid.AddItem(title, 0, 0, 1, 1, 1, 1, true)
	grid.AddItem(body, 1, 0, 1, 1, 1, 1, false)

	widget.Primitive = grid
	widget.Name = fmt.Sprintf("node%v", NodeNum)
	NodeNum++
	return
}
