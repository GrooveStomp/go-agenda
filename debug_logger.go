package main

import (
	"fmt"
	"github.com/rivo/tview"
)

type DebugLog struct {
	Primitive *tview.TextView
}

func (log DebugLog) Log(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	log.Primitive.SetText(str)
}
