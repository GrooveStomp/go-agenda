package main

import (
	"github.com/rivo/tview"
)

type Page struct {
	Name      string
	Primitive tview.Primitive
}

type PageStack struct {
	Pages []*Page
}

func (stack *PageStack) Push(p *Page) {
	stack.Pages = append(stack.Pages, p)
}

func (stack *PageStack) Pop() (p *Page) {
	p = stack.Top()
	if p == nil {
		return
	}

	stack.Pages = stack.Pages[:len(stack.Pages)-1]
	return
}

func (stack *PageStack) Top() (p *Page) {
	if len(stack.Pages) < 1 {
		return nil
	}

	p = stack.Pages[len(stack.Pages)-1]
	return
}

func (stack *PageStack) IndexName(name string) int {
	for i := range stack.Pages {
		p := stack.Pages[i]
		if p.Name == name {
			return i
		}
	}
	return -1
}

func (stack *PageStack) IndexPrimitive(primitive tview.Primitive) int {
	for i := range stack.Pages {
		p := stack.Pages[i]
		if p.Primitive == primitive {
			return 1
		}
	}
	return -1
}
