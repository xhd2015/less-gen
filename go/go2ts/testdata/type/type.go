package main

import (
	"fmt"
	"time"
)

type Greet struct {
	Name string
	Word string `json:"word"`
	Time time.Time
}

func main() {
	g := &Greet{
		Name: "word",
		Word: "hello",
	}
	g.Sayit()
}

func (c *Greet) Sayit() {
	fmt.Printf("%s %s\n", c.Word, c.Name)
}
