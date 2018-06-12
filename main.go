package main

/*

TODO:
[ ] Allow editing an agenda item
  [✓] <enter> on the item: brings up an edit dialog prepopulated with text.
  [✓] <enter> in the dialog: returns to main view, list item is updated.
[ ] Properly draw nested lists
[✓] Allow moving an item up and down in the list.
[✓] Allow adding a child item.
[ ] Allow indenting/de-indenting an item.
[ ] Allow nesting agendas and fluidly adding siblings and whatnot.
  [ ] Ah crap.  That means the text from the pop-up can't just be taken "as-is"
[ ] Modularize everything!  Ha ha.
  [ ] Separation of concerns, that kind of things.  UI and BE are intermixed pretty liberally.
[ ] Prototype and experiment with fluidity and mechanics of building nestable, hierarchical lists in the UI.
[ ] Implement a textarea widget, perhaps built upon gemacs or micro or gomacs?

*/

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os"
	"strings"
)

var (
	log      DebugLog
	boxShown bool
)

func main() {
	var editNode func(*AgendaNode, *AgendaNode, func(*AgendaNode))
	var editNodeCallback func(*AgendaNode)

	newNodeStack := Stack{}
	boxShown = false
	rootAgendaNode := NewNode("", "", []string{})

	mainGrid := tview.NewGrid()

	log.Primitive = tview.NewTextView()
	log.Primitive.SetBorder(true)
	log.Log("Program loaded")

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Agenda")
	list.SetSelectedFunc(func(index int, title string, body string, _ rune) {
		node := rootAgendaNode.Children[index]
		editNode(nil, node, editNodeCallback)
	})

	flex := tview.NewFlex()
	flex.SetFullScreen(false)
	flex.SetDirection(tview.FlexColumn)
	flex.AddItem(list, 0, 1, true)

	help := tview.NewTextView()
	help.SetBorder(true)
	help.SetTitle("help")
	help.SetText(strings.TrimSpace(helpText))

	app := tview.NewApplication()
	inputStack := InputHandlerStack{}
	pageStack := PageStack{}

	pages := tview.NewPages()
	pages.AddPage("main", flex, true, true)
	pages.AddPage("help", help, true, false)
	pageStack.Push(&Page{Name: "main", Primitive: flex})

	mainGrid.SetRows(-1, 3)
	mainGrid.SetColumns(-1)
	mainGrid.AddItem(pages, 0, 0, 1, 1, 1, 1, true)
	mainGrid.AddItem(log.Primitive, 1, 0, 1, 1, 1, 1, false)

	helpWidget := Widget{}
	helpWidget.Primitive = help
	helpWidget.InputHandler = createEscHandler(func() {
		inputStack.Pop()
		pageStack.Pop()
		pages.SwitchToPage(pageStack.Top().Name)
		log.Log("Exiting help, switching to %v", pageStack.Top().Name)
		app.Draw()
	})

	box := tview.NewBox()
	box.SetBorder(true)
	box.SetTitle("A [red]c[yellow]o[green]l[darkcyan]o[blue]r[darkmagenta]f[red]u[yellow]l[white] [black:red]c[:yellow]o[:green]l[:darkcyan]o[:blue]r[:darkmagenta]f[:red]u[:yellow]l[white:] [::bu]title")

	boxWidget := Widget{}
	boxWidget.Primitive = box
	boxWidget.InputHandler = createEscHandler(func() {
		inputStack.Pop()
		boxShown = false
		log.Log("Exiting box, switching to main")
		flex.RemoveItem(box)
		app.Draw()
	})

	flexWidget := Widget{}
	flexWidget.Primitive = flex
	flexWidget.InputHandler = func(event *tcell.EventKey) (result *tcell.EventKey) {
		result = event

		if event.Key() == tcell.KeyRune && event.Rune() == 't' {
			if !boxShown {
				boxWidget.InputHandlerIndex = inputStack.Push(boxWidget.InputHandler)
				flex.AddItem(box, 0, 1, true)
				app.SetFocus(box)
				boxShown = true
				log.Log("Showing box")
				result = nil
			}

			app.Draw()
		}

		return
	}

	addNodeCallback := func(node *AgendaNode) {
		list.AddItem(node.Title, node.Text, 0, nil)
	}

	editNodeCallback = func(node *AgendaNode) {
		index := list.GetCurrentItem()
		list.SetItemText(index, node.Title, node.Text)
	}

	editNode = func(scratch *AgendaNode, node *AgendaNode, callback func(node *AgendaNode)) {
		editNodeWidget := NewEditAgendaNodeWidget(app, node, scratch)
		editNodeWidget.InputHandler = createEscHandler(func() {
			if newNodeStack.Count() <= 1 {
				inputStack.Enable(flexWidget.InputHandlerIndex)
			}
			if scratch != nil {
				scratch.AddChild(node)
			} else {
				rootAgendaNode.AddChild(node)
			}
			inputStack.Pop()
			pageStack.Pop()
			newNodeStack.Pop()
			pages.SwitchToPage(pageStack.Top().Name)
			log.Log("Exiting %v, switching to %v", editNodeWidget.Name, pageStack.Top().Name)
			callback(node)
			app.Draw()
		})
		editNodeWidget.InputHandlerIndex = inputStack.Push(editNodeWidget.InputHandler)
		inputStack.Disable(flexWidget.InputHandlerIndex)
		pageStack.Push(&Page{Name: editNodeWidget.Name, Primitive: editNodeWidget.Primitive})
		pages.AddPage(editNodeWidget.Name, editNodeWidget.Primitive, true, true)
		pages.SwitchToPage(editNodeWidget.Name)
		log.Log("Showing %s", editNodeWidget.Name)
	}

	pagesWidget := Widget{}
	pagesWidget.Primitive = pages
	pagesWidget.InputHandler = func(event *tcell.EventKey) (result *tcell.EventKey) {
		result = event

		switch event.Key() {
		case tcell.KeyPgUp:
			log.Log("<pgup>")

			// get index of currently selected item.
			index := list.GetCurrentItem()
			if index == 0 {
				return nil
			}
			log.Log("index: %v, num items: %v", index, list.GetItemCount())

			node := rootAgendaNode.Children[index]
			prevNode := rootAgendaNode.Children[index-1]

			list.SetItemText(index-1, node.Title, node.Text)
			list.SetItemText(index, prevNode.Title, prevNode.Text)
			list.SetCurrentItem(index - 1)
			node.ShuffleUp()

			result = nil

		case tcell.KeyPgDn:
			log.Log("<pgdn>")

			// get index of currently selected item.
			index := list.GetCurrentItem()
			if list.GetItemCount() == index+1 {
				return nil
			}
			log.Log("index: %v, num items: %v", index, list.GetItemCount())

			node := rootAgendaNode.Children[index]
			nextNode := rootAgendaNode.Children[index+1]

			list.SetItemText(index+1, node.Title, node.Text)
			list.SetItemText(index, nextNode.Title, nextNode.Text)
			list.SetCurrentItem(index + 1)
			node.ShuffleDown()

			result = nil

		case tcell.KeyRune:
			switch event.Rune() {
			case '?':
				if pageStack.Top().Name == "help" {
					break
				}

				helpWidget.InputHandlerIndex = inputStack.Push(helpWidget.InputHandler)
				pageStack.Push(&Page{Name: "help", Primitive: helpWidget.Primitive})
				pages.SwitchToPage("help")
				log.Log("Showing help")
				result = nil

			case '+':
				var scratchNode *AgendaNode = nil
				if newNodeStack.Top() != nil {
					scratchNode = newNodeStack.Top().(*AgendaNode)
				}
				node := &AgendaNode{}
				newNodeStack.Push(node)
				editNode(scratchNode, node, addNodeCallback)
				result = nil
			}
		}
		app.Draw()

		return
	}

	pagesWidget.InputHandlerIndex = inputStack.Push(pagesWidget.InputHandler)
	flexWidget.InputHandlerIndex = inputStack.Push(flexWidget.InputHandler)

	app.SetFocus(pages)
	app.SetInputCapture(createAppInputHandler(&inputStack))
	if err := app.SetRoot(mainGrid, true).Run(); err != nil {
		panic(err)
	}

	rootAgendaNode.Walk(func(node *AgendaNode, indentLevel int) {
		node.Print(os.Stdout, indentLevel, 5)
	})
}

var helpText = `
?           Show this help text.
+           Add a new item.
<tab>       Expand an item.
<ctrl+c>    Quit
<esc>       Quit any popups, dialogs or modals.
<enter>     Edit selected item.
<ctrl+up>   Move an item up in the list. (Preserves nesting level.)
<ctrl+down> Move an item down in the list. (Preserves nesting level.)
`
