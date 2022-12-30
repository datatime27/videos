package boards

import (
	"fmt"
	"strings"

	"pieces"
)

const (
	hr = "-------------------------------------\n"

	infinity = 127
	deadKing = 100
)

var (
	filePrint = map[int]string{
		0: "a",
		1: "b",
		2: "c",
		3: "d",
		4: "e",
		5: "f",
		6: "g",
		7: "h",
	}
	fileParse = map[string]int{
		"a": 0,
		"b": 1,
		"c": 2,
		"d": 3,
		"e": 4,
		"f": 5,
		"g": 6,
		"h": 7,
	}

	rankPrint = map[int]string{
		0: "1",
		1: "2",
		2: "3",
		3: "4",
		4: "5",
		5: "6",
		6: "7",
		7: "8",
	}
	rankParse = map[string]int{
		"1": 0,
		"2": 1,
		"3": 2,
		"4": 3,
		"5": 4,
		"6": 5,
		"7": 6,
		"8": 7,
	}

	rookDirs = [][]int{
		{-1, 0},
		{1, 0},
		{0, 1},
		{0, -1},
	}
	diagonalDirs = [][]int{
		{-1, -1},
		{1, 1},
		{-1, 1},
		{1, -1},
	}
	allDirs = append(rookDirs, diagonalDirs...)

	sliderPieceDirs = map[int8][][]int{
		pieces.Rook:   rookDirs,
		pieces.Bishop: diagonalDirs,
		pieces.Queen:  allDirs,
		pieces.King:   allDirs,
	}
	sliderPieceDistance = map[int8]int{
		pieces.Rook:   7,
		pieces.Bishop: 7,
		pieces.Queen:  7,
		pieces.King:   1,
	}
)

type Context struct {
	MaxDepth int8
	Color    int8

	AllNodes  int64
	LeafNodes int64

	pruneTheRest bool
}

func NewContext(maxDepth, color int8) *Context {
	return &Context{
		MaxDepth: maxDepth,
		Color:    color,
	}
}

type Board struct {
	board      [8][8]int8
	depth      int8
	myTurn     bool
	History    string
	Evaluation int8
	FirstMove  *Board
}

type vector struct {
	rank int
	file int
}

func printLocation(rank, file int) string {
	if rank < 0 || rank >= 8 || file < 0 || file >= 8 {
		return fmt.Sprintf("Offboard[rank:%v,file:%v]", rank, file)
	}

	return filePrint[file] + rankPrint[rank]
}

func fileHeader() string {
	var sb strings.Builder
	sb.WriteString("  | ")
	for file := 0; file < 8; file++ {
		sb.WriteString(filePrint[file])
		sb.WriteString(" | ")
	}
	sb.WriteString("\n")
	return sb.String()
}
func (b *Board) String() string {
	var sb strings.Builder

	// sb.WriteString(fmt.Sprintf("color: %v, depth: %v, myTurn: %v\n", b.color, b.depth, b.myTurn))
	sb.WriteString(fileHeader())
	sb.WriteString(hr)
	for rank := 7; rank >= 0; rank-- {
		sb.WriteString(fmt.Sprintf("%v | ", rankPrint[rank]))
		for file := 0; file < 8; file++ {
			p := pieces.Decode(b.board[rank][file])
			sb.WriteString(fmt.Sprintf("%v | ", p.Print()))
		}
		sb.WriteString(fmt.Sprintf("%v\n", rankPrint[rank]))
	}
	sb.WriteString(hr)
	sb.WriteString(fileHeader())

	return sb.String()
}
func (b *Board) DebugString() string {
	return fmt.Sprintf("Depth: %v Evalution: %v History: %v", b.depth, b.Evaluation, b.History)
}
func (b *Board) ChildNode() *Board {
	newBoard := &Board{
		board:   b.board,
		depth:   b.depth + 1,
		myTurn:  !b.myTurn,
		History: b.History,
	}
	if b.FirstMove == nil {
		newBoard.FirstMove = newBoard
	} else {
		newBoard.FirstMove = b.FirstMove
	}
	return newBoard
}

func getLocation(location string) (int, int, bool) {
	if len(location) != 2 {
		panic("Unknown location: " + location)
	}
	file := fileParse[string(location[0])]
	rank := rankParse[string(location[1])]

	isValid := rank >= 0 && rank < 8 && file >= 0 && file < 8
	return rank, file, isValid
}

func (b *Board) GetPiece(location string) *pieces.Piece {
	rank, file, isVaild := getLocation(location)
	if !isVaild {
		return nil
	}
	return pieces.Decode(b.board[rank][file])
}

func (b *Board) SetPiece(location, piece string) {
	rank, file, isVaild := getLocation(location)
	if !isVaild {
		return
	}
	p := pieces.Parse(piece)
	b.board[rank][file] = p.Encode()
}

func (b *Board) MovePiece(ctx *Context, oldrank, oldfile, newrank, newfile int) {
	old := &b.board[oldrank][oldfile]

	new := &b.board[newrank][newfile]
	*new = *old
	*old = pieces.None

	p := pieces.Decode(*new)
	b.History += fmt.Sprintf("%v to %v, ", p.Print(), printLocation(newrank, newfile))
	b.EvaluateMaterial(ctx)
}

func (b *Board) EvaluateMaterial(ctx *Context) {
	var sum int = 0
	kings := map[int8]bool{}

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			piece := pieces.Decode(b.board[rank][file])
			if piece.Piece == pieces.King {
				kings[piece.Color] = true
			} else {
				v := piece.Value()

				if piece.IsSameColor(ctx.Color) {
					sum += v
				} else {
					sum -= v
				}
			}
		}
	}

	// Both Kings are still present - return material value
	if len(kings) == 2 {
		b.Evaluation = int8(sum)
		return
	}

	// My king is dead - this is infinitely bad
	if !kings[ctx.Color] {
		b.Evaluation = -deadKing
		return
	}

	// Check or check mate - this is infintely good
	b.Evaluation = deadKing
}
func (b *Board) spawnDecision(ctx *Context, oldRank, oldFile, newRank, newFile int,
	piece *pieces.Piece, allowTake, allowNoTake bool, alpha, beta int8, leafBoard **Board) bool {
	if !b.moveAvailable(newRank, newFile, piece, allowTake, allowNoTake) {
		return false
	}
	newBoard := b.ChildNode()
	newBoard.MovePiece(ctx, oldRank, oldFile, newRank, newFile)
	ctx.AllNodes++
	var minimaxBoard *Board
	// Terminal node - evaulate
	if newBoard.depth == ctx.MaxDepth || newBoard.Evaluation == deadKing || newBoard.Evaluation == -deadKing {
		// fmt.Printf("Eval %v\n", printLocation(newRank, newFile))
		ctx.LeafNodes++
		minimaxBoard = newBoard
	} else {
		// fmt.Printf("%v Child FindMoves %v\n", printLocation(oldRank, oldFile), printLocation(newRank, newFile))
		minimaxBoard = newBoard.minimax(ctx, alpha, beta) // Recursive call down
	}

	// fmt.Printf("%v %v\n", b.DebugString(), minimaxBoard.DebugString())

	if b.myTurn {
		if minimaxBoard.Evaluation > (*leafBoard).Evaluation ||
			(minimaxBoard.Evaluation == (*leafBoard).Evaluation && b.depth < minimaxBoard.depth) {
			*leafBoard = minimaxBoard
		}
		if minimaxBoard.Evaluation > alpha {
			alpha = minimaxBoard.Evaluation
			// fmt.Printf("alpha = %v\n", alpha)
		}
	} else if !b.myTurn {
		if minimaxBoard.Evaluation < (*leafBoard).Evaluation ||
			(minimaxBoard.Evaluation == (*leafBoard).Evaluation && b.depth > minimaxBoard.depth) {
			*leafBoard = minimaxBoard
		}
		if minimaxBoard.Evaluation < beta {
			beta = minimaxBoard.Evaluation
			// fmt.Printf("beta = %v\n", beta)
		}
	}
	// fmt.Printf("alpha = %v, beta = %v\n", alpha, beta)

	if beta <= alpha {
		ctx.pruneTheRest = true
		fmt.Printf("Prune alpha %v\n", minimaxBoard.DebugString())
	}
	return true
}

func (b *Board) moveAvailable(rank, file int, piece *pieces.Piece, allowTake, allowNoTake bool) bool {
	isValid := rank >= 0 && rank < 8 && file >= 0 && file < 8
	if !isValid {
		return false
	}

	newPiece := pieces.Decode(b.board[rank][file])
	// Cell is empty

	if allowNoTake && newPiece.Piece == pieces.None {
		return true
	}
	// Take opponent piece
	if allowTake && newPiece.IsOppositeColor(piece.Color) {
		return true
	}

	// Move  not allowed
	return false
}

func (b *Board) FindMoves(ctx *Context) *Board {
	return b.minimax(ctx, -infinity, infinity)
}
func (b *Board) minimax(ctx *Context, alpha, beta int8) *Board {
	ctx.pruneTheRest = false
	leafBoard := &Board{History: b.History}
	if b.myTurn {
		leafBoard.Evaluation = -infinity
	} else {
		leafBoard.Evaluation = infinity
	}

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			piece := pieces.Decode(b.board[rank][file])
			// Only look at the correct pieces for this turn
			if (piece.IsOppositeColor(ctx.Color)) == b.myTurn {
				continue
			}

			switch piece.Piece {
			case pieces.Pawn:
				b.movePawn(ctx, rank, file, piece, alpha, beta, &leafBoard)
			case pieces.Rook, pieces.Bishop, pieces.Queen, pieces.King:
				b.moveSlidingPiece(ctx, rank, file, piece, alpha, beta, &leafBoard)
			default:
				continue
			}
			if ctx.pruneTheRest {
				fmt.Printf("Pruning %v %v\n", b.DebugString(), leafBoard.DebugString())
				ctx.pruneTheRest = false
				return leafBoard
			}
		}
	}
	if leafBoard.Evaluation == infinity-1 || leafBoard.Evaluation == -infinity+1 {
		// fmt.Printf("No moves found for %v %v %v\n", leafBoard.Evaluation, b.Evaluation, b.History)
		leafBoard = b
	}
	return leafBoard
}

func (b *Board) movePawn(ctx *Context, rank, file int, piece *pieces.Piece, alpha, beta int8, leafBoard **Board) {
	var direction int = 1
	if piece.Color == pieces.Black {
		direction = -1
	}
	// Pawn forward one space - no take
	b.spawnDecision(ctx, rank, file, rank+direction, file, piece, false, true, alpha, beta, leafBoard)
	if ctx.pruneTheRest {
		return
	}

	// Pawn diagonal - take only
	b.spawnDecision(ctx, rank, file, rank+direction, file+1, piece, true, false, alpha, beta, leafBoard)
	if ctx.pruneTheRest {
		return
	}
	b.spawnDecision(ctx, rank, file, rank+direction, file-1, piece, true, false, alpha, beta, leafBoard)
	if ctx.pruneTheRest {
		return
	}

	// Pawn forward 2 - Init only - no take
	// This doesn't work if the cell in front of the pawn is occupied
	// fmt.Printf("movePawn %v %v %v\n", rank+direction, file, b.board[rank+direction][file])
	if ((rank == 1 && piece.Color == pieces.White) || (rank == 6 && piece.Color == pieces.Black)) &&
		b.board[rank+direction][file] == pieces.None {
		b.spawnDecision(ctx, rank, file, rank+direction*2, file, piece, false, true, alpha, beta, leafBoard)
		if ctx.pruneTheRest {
			return
		}
	}

}
func (b *Board) moveSlidingPiece(ctx *Context, rank, file int, piece *pieces.Piece, alpha, beta int8, leafBoard **Board) {
	dirs := sliderPieceDirs[piece.Piece]
	distance := sliderPieceDistance[piece.Piece]
	for _, dir := range dirs {
		rankDir := dir[0]
		fileDir := dir[1]
		// fmt.Printf("%v:  Dir: %v %v\n", piece.Print(), rankDir, fileDir)

		for i := 1; i <= distance; i++ {
			newRank := rank + rankDir*i
			newFile := file + fileDir*i

			// fmt.Printf("dir[%v,%v] q to %v\n", rankDir, fileDir, printLocation(newRank, newFile))

			// Piece is unable to move in this direction - No need to continue further
			if !b.spawnDecision(ctx, rank, file, newRank, newFile, piece, true, true, alpha, beta, leafBoard) {
				// fmt.Printf("Unable to move in this direction i: %v %v:  %v\n", i, piece.Print(), printLocation(newRank, newFile))
				break
			}
			if ctx.pruneTheRest {
				return
			}

			newPiece := pieces.Decode(b.board[newRank][newFile])

			// piece has hit a piece along this direction (opponent piece) - No need to continue further
			if newPiece.IsOppositeColor(piece.Color) {
				// fmt.Printf("Opponent hit i: %v %v:  %v %v\n", i, piece.Print(), printLocation(newRank, newFile), b.board[newRank][newFile])
				break
			}
			// fmt.Printf("Continue forward i: %v %v:  %v\n", i, piece.Print(), printLocation(newRank, newFile))

		}
	}
}

func ParseBoard(data []byte, myTurn bool) *Board {
	b := &Board{
		depth:  0,
		myTurn: myTurn,
	}

	s := string(data)
	rows := strings.Split(s, "\n")

	for _, row := range rows {
		cells := strings.Split(row, "|")
		if len(cells) < 10 {
			continue
		}
		rankString := strings.TrimSpace(cells[0])
		rank, found := rankParse[rankString]
		if !found {
			continue
		}
		for file, cell := range cells[1:9] {
			piece := pieces.Parse(strings.TrimSpace(cell))
			b.board[rank][file] = piece.Encode()
		}
	}

	return b
}
