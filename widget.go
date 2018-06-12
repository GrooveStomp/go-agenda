package main

import (
	"github.com/rivo/tview"
)

type Widget struct {
	//Id WidgetId
	Name              string
	Primitive         tview.Primitive
	InputHandler      InputHandler
	InputHandlerIndex int
}
