package main

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"

	"github.com/levipare/go-todo-list/todo"
)

func readListFromFile(filename string) todo.TodoList {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	todoList := todo.NewList()

	for _, item := range lines {
		isCompleted, _ := strconv.ParseBool(item[1])
		todoList.AddItem(todo.NewItem(item[0], isCompleted))
	}

	return todoList
}

func writeListToFile(filename string, lst todo.TodoList) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for _, item := range lst.Items {
		record := []string{item.Description, strconv.FormatBool(item.Completed)}
		err := writer.Write(record)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
}

var todoList todo.TodoList
var selectedItem = 0

func main() {
	todoList = readListFromFile("items.csv")
	render()
	writeListToFile("items.csv", todoList)
}

func render() {
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err = s.Init(); err != nil {
		panic(err)
	}
	defer s.Fini()

	s.SetStyle(tcell.StyleDefault)
	s.Clear()

	// Blocking input loop
loop:
	for {
		drawItems(s, todoList)
		s.Show()

		switch ev := s.PollEvent().(type) {
		case *tcell.EventKey:
			if ev.Key() == 256 {
				switch ev.Rune() {
				case ' ':
					if selectedItem < len(todoList.Items) {
						todoList.GetItem(selectedItem).Completed = !todoList.GetItem(selectedItem).Completed
					}
				}
			} else {
				switch ev.Key() {
				case tcell.KeyEscape:
					break loop
				case tcell.KeyDown, tcell.KeyTab:
					if selectedItem < len(todoList.Items)-1 {
						selectedItem++
					}
				case tcell.KeyUp, tcell.KeyBacktab:
					if selectedItem > 0 {
						selectedItem--
					}
				case tcell.KeyEnter, ' ':
					if selectedItem < len(todoList.Items) {
						todoList.GetItem(selectedItem).Completed = !todoList.GetItem(selectedItem).Completed
					}
				}
			}

		}

	}
}

func drawItems(s tcell.Screen, lst todo.TodoList) {
	for i, item := range lst.Items {
		style := tcell.StyleDefault
		if item.Completed {
			style = style.Dim(true)
		}

		if i == selectedItem {
			style = style.Foreground(tcell.ColorDeepPink)
		}

		symbol := ' '
		if item.Completed {
			symbol = 'X'
		}

		s.SetContent(0, i, '[', nil, style)
		s.SetContent(1, i, symbol, nil, style)
		s.SetContent(2, i, ']', nil, style)
		for j, c := range item.Description {
			s.SetContent(j+4, i, c, nil, style)
		}
	}
}
