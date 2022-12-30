package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"flag"
	"boards"
	"pieces"
)

var (
	path  = flag.String("path", "", "location of file with board layout")
	color = flag.String("color", "w", "color of the  player w|b")
	depth = flag.Int("depth", 3, "How many ply forward to look ahead")
)

// StringPrompt asks for a string value using the label
func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func main() {

	var myColor int8
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

	board := boards.ParseBoard(data, true)
	// board.MovePiece(6, 0, 2, 0) // Pawn
	// board.MovePiece(0, 3, 3, 1) // Queen
	fmt.Printf("Value: %d\n", board.Evaluation)
	fmt.Printf("%v\n", board)

	fmt.Println("Calculating....\n")

	start := time.Now()
	leafBoard := board.FindMoves(ctx)
	duration := time.Duration(time.Now().Sub(start))
	fmt.Printf("Time: %v\nAll Nodes: %d\nLeaf Nodes: %d\nValue: %d\n", duration, ctx.AllNodes, ctx.LeafNodes, leafBoard.Evaluation)
	fmt.Printf("History: %v\n", leafBoard.History)

	fmt.Printf("First Move:\n%v\n", leafBoard.FirstMove.String())

	fmt.Printf("Leaf Move:\nValue:%v\n%v\n", leafBoard.Evaluation, leafBoard.String())

	answer := StringPrompt("Write file? y/n")
	if answer == "y" {
		if err := os.WriteFile(*path, []byte(leafBoard.FirstMove.String()), 0666); err != nil {
			panic(err)
		}
		fmt.Println("Wrote file")
	}

}
