package main

import (
	"github.com/rivo/tview"
)

type Widget struct {
	//Id WidgetId
	Primitive    tview.Primitive
	InputHandler InputHandler
}
