package main

import "fmt"

type animal interface {
	move()
}

type squirrel struct {
	name string
}

func (s *squirrel) run() {
	fmt.Println("Run run run !")
}

func (s *squirrel) move() {
	fmt.Println("move")
}

func doMove(a animal) {
	fmt.Println("ggg")
}

func main() {
	s1 := squirrel{"yasser"}
	fmt.Print(&s1)
}
