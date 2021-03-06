package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	EditAgendaNodeDialogNum = 1
)

func NewEditAgendaNodeWidget(app *tview.Application, node *AgendaNode, scratch *AgendaNode) (widget *Widget) {
	widget = &Widget{}

	title := tview.NewInputField()
	body := tview.NewInputField()

	titleText := "Title"
	if scratch != nil {
		titleText = fmt.Sprintf("Title (child of %v):", scratch.Title)
	}

	title.SetBorder(true)
	title.SetTitle(titleText)
	title.SetText(node.Title)
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
		node.Title = title.GetText()
		log.Log(text)
		app.Draw()
	})

	body.SetBorder(true)
	body.SetTitle("Body")
	body.SetText(node.Text)
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
		node.Text = body.GetText()
		log.Log(text)
		app.Draw()
	})

	grid := tview.NewGrid()
	grid.SetRows(3, -1, 3)
	grid.SetColumns(-1)

	grid.AddItem(title, 0, 0, 1, 1, 1, 1, true)
	grid.AddItem(body, 1, 0, 1, 1, 1, 1, false)

	widget.Primitive = grid
	widget.Name = fmt.Sprintf("EditAgenda%v", EditAgendaNodeDialogNum)
	EditAgendaNodeDialogNum++
	return
}
