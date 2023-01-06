package main

import (
	"fmt"
	"os"

	"chessbot/lib/boards"
	"flag"
)

var (
	path = flag.String("path", "", "location of file with board layout")
)

func main() {
	flag.Parse()
	if *path == "" {
		panic("Must provide --path")
	}
	move := flag.Arg(0)

	ctx := boards.NewContext(0, 0)
	data, err := os.ReadFile(*path)
	if err != nil {
		panic(err)
	}
	move = flag.Arg(0)
	if len(move) != 4 {
		panic("Must provide move like a2b3")
	}
	src := move[0:2]
	dst := move[2:4]
	board := boards.ParseBoard(ctx, data, true)
	piece := board.GetPiece(src)
	ctx.Color = piece.Color
	board.EvaluateMaterial(ctx)
	fmt.Printf("%v\n", board)

	fmt.Println("Moving piece...\n")

	board.SetMove(ctx, src, dst)
	board.EvaluateMaterial(ctx)
	fmt.Printf("%v\n", board)

	if err := os.WriteFile(*path, []byte(board.Format(false)), 0666); err != nil {
		panic(err)
	}
	fmt.Println("Wrote file")

}
