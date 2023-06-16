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
	acceptMove  = flag.Bool("y", false, "Automatically accept move and file write.")
)

func main() {
	flag.Parse()
	if *path == "" {
		panic("Must provide --path")
	}

	data, err := os.ReadFile(*path)
	if err != nil {
		panic(err)
	}

	board := boards.ParseBoard(&boards.Context{}, data, true)
	fmt.Printf("%v\n", board)

	if err := board.WriteHTML(*displayPath); err != nil {
		panic(err)
	}
	fmt.Println("Wrote file")
}
