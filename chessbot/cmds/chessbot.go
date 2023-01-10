package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"chessbot/lib/boards"
	"chessbot/lib/pieces"
	"flag"
)

var (
	path  = flag.String("path", "..\\boards\\test.txt", "location of file with board layout")
	color = flag.String("color", "w", "color of the  player w|b")
	depth = flag.Int("depth", 1, "How many ply forward to look ahead")
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
	fmt.Printf("Value: %d\n", board.Evaluation)
	fmt.Printf("%v\n", board)

	fmt.Println("Calculating....\n")

	start := time.Now()
	leafBoard := board.FindMoves(ctx)
	duration := time.Duration(time.Now().Sub(start))
	nodesPerSec := float64(ctx.AllNodes) / duration.Seconds()
	leafNodesPerSec := float64(ctx.LeafNodes) / duration.Seconds()
	fmt.Printf("Time: %v\n", duration)
	fmt.Printf("All Nodes: %d (%.2f per sec)\n", ctx.AllNodes, nodesPerSec)
	fmt.Printf("Leaf Nodes: %d (%.2f per sec)\n\n", ctx.LeafNodes, leafNodesPerSec)

	fmt.Printf("Leaf Move:\n%v\n\n", leafBoard.String())
	fmt.Printf("Next Move:\n%v\n\n", leafBoard.FirstMove.String())

	// answer := StringPrompt("Write file? y/n")
	// if answer != "n" {
	// 	if err := os.WriteFile(*path, []byte(leafBoard.FirstMove.Format(false)), 0666); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("Wrote file")
	// }

}
