package main

import (
	"github.com/gdamore/tcell"
)

type InputHandler func(*tcell.EventKey) *tcell.EventKey

type InputHandlerStack struct {
	InputHandler []InputHandler
}

func (s *InputHandlerStack) Push(f InputHandler) {
	s.InputHandler = append(s.InputHandler, f)
}

func (s *InputHandlerStack) Pop() (handler InputHandler) {
	handler = s.Top()
	if handler == nil {
		return
	}

	s.InputHandler = s.InputHandler[:len(s.InputHandler)-1]

	return handler
}

func (s *InputHandlerStack) Top() InputHandler {
	if len(s.InputHandler) < 1 {
		return nil
	}

	handler := s.InputHandler[len(s.InputHandler)-1]

	return handler
}

func createEscHandler(callback func()) InputHandler {
	return func(eventKey *tcell.EventKey) *tcell.EventKey {
		if eventKey.Key() == tcell.KeyEsc {
			callback()
			// NOTE(AARONO): Apparently returning nil here causes laggy behavior where
			// it seems like Esc needs to be hit twice.
			return nil
		}
		return eventKey
	}
}

func createAppInputHandler(stack *InputHandlerStack) InputHandler {
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
