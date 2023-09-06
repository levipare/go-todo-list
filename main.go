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
var textInput = ""

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

	// Blocking input loop
loop:
	for {
		s.Clear()
		s.HideCursor()

		drawItems(s)
		drawItemInput(s)
		drawClearKeyMap(s)

		s.Show()

		switch ev := s.PollEvent().(type) {
		case *tcell.EventKey:
			if selectedItem == len(todoList.Items) {
				if ev.Key() == 256 {
					textInput = textInput + string(ev.Rune())
				}

			} else if ev.Key() == 256 {
				switch ev.Rune() {
				case ' ':
					if selectedItem < len(todoList.Items) {
						todoList.GetItem(selectedItem).Completed = !todoList.GetItem(selectedItem).Completed
					}
				}
			}

			switch ev.Key() {
			case tcell.KeyEscape:
				break loop
			case tcell.KeyDown, tcell.KeyTab:
				if selectedItem < len(todoList.Items) {
					selectedItem++
				}
			case tcell.KeyUp, tcell.KeyBacktab:
				if selectedItem > 0 {
					selectedItem--
				}
			case tcell.KeyEnter, ' ':
				if selectedItem < len(todoList.Items) {
					todoList.GetItem(selectedItem).Completed = !todoList.GetItem(selectedItem).Completed
				} else if selectedItem == len(todoList.Items) && textInput != "" {
					todoList.AddItem(todo.NewItem(textInput, false))
					textInput = ""
					selectedItem = len(todoList.Items)
				}
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if selectedItem == len(todoList.Items) && len(textInput) > 0 {
					textInput = textInput[:len(textInput)-1]
				}
			case tcell.KeyCtrlX:
				todoList = todo.NewList()
			}

		}

	}
}

func drawClearKeyMap(s tcell.Screen) {
	for i, c := range "to clear list -> ctrl+x" {
		_, h := s.Size()
		s.SetContent(i, h-1, c, nil, tcell.StyleDefault.Dim(true))
	}
}

func drawItemInput(s tcell.Screen) {
	oX, oY, w, h := 0, 0, 0, 0
	oY = len(todoList.Items) + 1
	w = 32
	h = 3

	style := tcell.StyleDefault
	if selectedItem == len(todoList.Items) {
		style = style.Foreground(tcell.ColorIndianRed)
		s.ShowCursor(oX+1+len(textInput), oY+1)
	}

	s.SetContent(oX, oY, '╭', nil, style)
	s.SetContent(oX+w-1, oY, '╮', nil, style)
	s.SetContent(oX, oY+h-1, '╰', nil, style)
	s.SetContent(oX+w-1, oY+h-1, '╯', nil, style)

	for x := oX + 1; x < oX+w-1; x++ {
		s.SetContent(x, oY, '─', nil, style)
		s.SetContent(x, oY+h-1, '─', nil, style)

		if x-oX-1 < len(textInput) {
			s.SetContent(x, oY+1, rune(textInput[x-oX-1]), nil, tcell.StyleDefault)
		}
	}

	for _, x := range []int{oX, oX + w - 1} {
		for y := oY + 1; y < oY+h-1; y++ {
			s.SetContent(x, y, '│', nil, style)
		}
	}

	for i, c := range " new todo " {
		s.SetContent(oX+i+1, oY, c, nil, style)
	}

}

func drawItems(s tcell.Screen) {
	for i, item := range todoList.Items {
		style := tcell.StyleDefault
		if item.Completed {
			style = style.Dim(true)
		}

		if i == selectedItem {
			style = style.Foreground(tcell.ColorIndianRed)
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
