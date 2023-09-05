package todo

type TodoItem struct {
	Description string
	Completed bool
}

type TodoList struct {
	Items []TodoItem
}

// Item Functions
func NewItem(desc string, completed bool) TodoItem {
	return TodoItem{desc, completed}
}


// List Functions
func NewList(items ...TodoItem)  TodoList{
	lst := TodoList{}
	lst.Items = append(lst.Items, items...) // Append optional items
	return lst
}

func (lst *TodoList) GetItem(index int) *TodoItem{
	return &lst.Items[index]
}

func (lst* TodoList) AddItem(item TodoItem) {
	lst.Items = append(lst.Items, item)
}

func (lst* TodoList) RemoveItem(index int){
	lst.Items = append(lst.Items[0:index], lst.Items[index+1:]...)
}

