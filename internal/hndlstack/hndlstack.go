package hndlstack

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

func (s *InputHandlerStack) Top() (InputHandler) {
	if len(s.InputHandler) < 1 {
		return nil
	}

	handler := s.InputHandler[len(s.InputHandler)-1]

	return handler
}
