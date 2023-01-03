package main

import (
	"fmt"
	"os"

	"chessbot/lib/boards"
	"chessbot/lib/pieces"
	"flag"
)

var (
	path  = flag.String("path", "", "location of file with board layout")
	color = flag.String("color", "w", "color of the  player w|b")
	depth = flag.Int("depth", 3, "How many ply forward to look ahead")
)

func main() {
	flag.Parse()
	if *path == "" {
		panic("Must provide --path")
	}

	var myColor pieces.Color
	if *color == "w" {
		myColor = pieces.White
	} else if *color == "b" {
		myColor = pieces.Black
	} else {
		panic("--color must be 'w' or 'b'")
	}
	ctx := boards.NewContext(int8(*depth), myColor)
	data, err := os.ReadFile(*path)
	if err != nil {
		panic(err)
	}

	board := boards.ParseBoard(ctx, data, true)
	// board.MovePiece(6, 0, 2, 0) // Pawn
	// board.MovePiece(0, 3, 3, 1) // Queen
	fmt.Printf("Value: %d\n", board.Evaluation)
	fmt.Printf("%v\n", board)

	fmt.Println("Calculating....\n")

	m := board.GetCellsOppoentPawnsCanCapture(ctx)
	fmt.Printf("cells: %v\n", m)
}
