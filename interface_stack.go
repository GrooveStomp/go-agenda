package main

type Stack struct {
	Items []interface{}
}

func (stack *Stack) Push(new interface{}) {
	stack.Items = append(stack.Items, new)
}

func (stack *Stack) Pop() (top interface{}) {
	top = stack.Top()
	if top == nil {
		return
	}

	stack.Items = stack.Items[:len(stack.Items)-1]
	return
}

func (stack *Stack) Top() interface{} {
	if len(stack.Items) < 1 {
		return nil
	}
	return stack.Items[len(stack.Items)-1]
}

func (stack *Stack) Index(needle interface{}) int {
	for i := range stack.Items {
		haystack := stack.Items[i]
		if haystack == needle {
			return 1
		}
	}
	return -1
}

func (stack *Stack) Count() int {
	return len(stack.Items)
}
