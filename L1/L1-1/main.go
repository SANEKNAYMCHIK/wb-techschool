package main

import "fmt"

type Human struct {
	Name   string
	Age    int
	Height int
	Weight int
}

func NewHuman(name string, age int, height int, weight int) *Human {
	return &Human{
		Name:   name,
		Age:    age,
		Height: height,
		Weight: weight,
	}
}

func (h *Human) meeting() {
	fmt.Printf("Hi! My name is %s and I'm %d years old.\n", h.Name, h.Age)
}

func (h *Human) info() {
	fmt.Printf("Height is %d.\nWeight is %d.\n", h.Height, h.Weight)
}

type Action struct {
	ActionName string
	Human
}

func NewAction(actionName string, human Human) *Action {
	return &Action{
		ActionName: actionName,
		Human:      human,
	}
}

func (a *Action) activity() {
	a.meeting()
	fmt.Printf("I like to %s!\n", a.ActionName)
}

func main() {
	running := NewAction("run", *NewHuman("Vasya", 17, 176, 54))
	running.meeting()
	running.info()
	running.activity()
}
