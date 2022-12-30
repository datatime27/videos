package main

import (
	"fmt"
	"os"

	"flag"
	"boards"
)

var (
	path = flag.String("path", "", "location of file with board layout")
)

func Init() *boards.Board {
	b := &boards.Board{}

	// White
	b.SetPiece("a1", "R")
	b.SetPiece("b1", "N")
	b.SetPiece("c1", "B")
	b.SetPiece("d1", "Q")
	b.SetPiece("e1", "K")
	b.SetPiece("f1", "B")
	b.SetPiece("g1", "N")
	b.SetPiece("h1", "R")

	b.SetPiece("a2", "P")
	b.SetPiece("b2", "P")
	b.SetPiece("c2", "P")
	b.SetPiece("d2", "P")
	b.SetPiece("e2", "P")
	b.SetPiece("f2", "P")
	b.SetPiece("g2", "P")
	b.SetPiece("h2", "P")

	// Black
	b.SetPiece("a8", "r")
	b.SetPiece("b8", "n")
	b.SetPiece("c8", "b")
	b.SetPiece("d8", "q")
	b.SetPiece("e8", "k")
	b.SetPiece("f8", "b")
	b.SetPiece("g8", "n")
	b.SetPiece("h8", "r")

	b.SetPiece("a7", "p")
	b.SetPiece("b7", "p")
	b.SetPiece("c7", "p")
	b.SetPiece("d7", "p")
	b.SetPiece("e7", "p")
	b.SetPiece("f7", "p")
	b.SetPiece("g7", "p")
	b.SetPiece("h7", "p")

	return b
}
func InitTest() *boards.Board {
	b := &boards.Board{}

	// White
	// b.SetPiece("a1", "R")
	//b.SetPiece("b1", "N")
	b.SetPiece("c1", "B")
	b.SetPiece("d1", "Q")
	b.SetPiece("e1", "K")
	b.SetPiece("f1", "B")
	//b.SetPiece("g1", "N")
	// b.SetPiece("h1", "R")

	// b.SetPiece("a2", "P")
	// b.SetPiece("b2", "P")
	// b.SetPiece("c2", "P")
	b.SetPiece("d2", "P")
	b.SetPiece("e2", "P")
	// b.SetPiece("f2", "P")
	// b.SetPiece("g2", "P")
	// b.SetPiece("h2", "P")

	// Black
	// b.SetPiece("a8", "r")
	//b.SetPiece("b8", "n")
	b.SetPiece("c8", "b")
	b.SetPiece("d8", "q")
	b.SetPiece("e8", "k")
	b.SetPiece("f8", "b")
	//b.SetPiece("g8", "n")
	// b.SetPiece("h8", "r")

	// b.SetPiece("a7", "p")
	// b.SetPiece("b7", "p")
	// b.SetPiece("c7", "p")
	b.SetPiece("d7", "p")
	b.SetPiece("e7", "p")
	// b.SetPiece("f7", "p")
	// b.SetPiece("g7", "p")
	// b.SetPiece("h7", "p")

	return b
}

func main() {
	board := InitTest()
	s := board.String()
	fmt.Printf("%v\n", s)
	if err := os.WriteFile(*path, []byte(s), 0666); err != nil {
		panic(err)
	}

}
