package main

import (
	"flag"
	"fmt"
	"os"

	"chessbot/lib/boards"
	"chessbot/lib/pieces"
)

var (
	path  = flag.String("path", "", "location of file with board layout")
	color = flag.String("color", "w", "color of the  player w|b")
)

func main() {
	flag.Parse()
	if *path == "" {
		panic("Must provide --path")
	}

	ctx := boards.NewContext(0, *color)
	data, err := os.ReadFile(*path)
	if err != nil {
		panic(err)
	}

	move := flag.Arg(0)
	if len(move) == 0 {
		panic("Must provide move")
	}
	board := boards.ParseBoard(ctx, data, true)
	board.EvaluateMaterial(ctx)
	fmt.Printf("%v\n", board)

	fmt.Println("Moving piece...\n")

	board.SetMove(move)
	board.EvaluateMaterial(ctx)
	fmt.Printf("%v\n", board)

	if err := os.WriteFile(*path, []byte(board.Format(false)), 0666); err != nil {
		panic(err)
	}
	fmt.Println("Wrote file")

}
