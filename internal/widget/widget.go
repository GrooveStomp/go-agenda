package widget

import (
)

// type WidgetId int

// var currentWidgetId WidgetId = 0

// func init() {
// }

// func nextWidgetId() (result WidgetId) {
// 	result = currentWidgetId
// 	currentWidgetId += 1
// }
	
type Widget struct {
	//Id WidgetId
	Primitive tview.Primitive
  InputHandler InputHandler	
}

// func NewWidget() *Widget {
// 	w := Widget{}
// 	w.Id = nextWidgetId()
// 	return &w
// }
