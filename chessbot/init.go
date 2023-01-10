package main

import (
	"fmt"
	"os"

	"chessbot/lib/boards"
	"flag"
)

var (
	path        = flag.String("path", "", "location of file with board layout")
	displayPath = flag.String("display", "html/display.html", "location output display html")
)

// Standard board layout
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

// Smaller board for testing.
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
	flag.Parse()
	if *path == "" {
		panic("Must provide --path")
	}
	board := Init()
	fmt.Printf("%v\n", board.String())
	if err := os.WriteFile(*path, []byte(board.Format(false)), 0666); err != nil {
		panic(err)
	}
	if err := board.WriteHTML(*displayPath); err != nil {
		panic(err)
	}
	fmt.Printf("Wrote: %v\n", *path)
}
