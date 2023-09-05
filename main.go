package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

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

	return todoList;
}

func writeListToFile(filename string, lst todo.TodoList)  {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for _, item := range lst.Items {
		err := writer.Write([]string{item.Description, strconv.FormatBool(item.Completed)})
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	todoList := readListFromFile("items.csv")
	todoList.AddItem(todo.NewItem("Make text editor", false))
	fmt.Println(todoList)
	writeListToFile("items.csv", todoList)
}
