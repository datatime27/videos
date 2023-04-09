package boards

import (
	"fmt"
	"sort"
	"strings"

	"chessbot/lib/pieces"
)

const (
	hr = "-------------------------------------\n"

	infinity      = 100000
	deadKing      = 1000
	historyPrefix = "# History: "
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

	knightDirs = [][]int{
		{2, 1},
		{2, -1},
		{1, 2},
		{1, -2},
		{-1, 2},
		{-1, -2},
		{-2, 1},
		{-2, -1},
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
	allSlidingDirs = append(rookDirs, diagonalDirs...)

	sliderPieceDirs = map[pieces.Class][][]int{
		pieces.Rook:   rookDirs,
		pieces.Bishop: diagonalDirs,
		pieces.Queen:  allSlidingDirs,
		pieces.King:   allSlidingDirs,
	}
	sliderPieceDistance = map[pieces.Class]int{
		pieces.Rook:   7,
		pieces.Bishop: 7,
		pieces.Queen:  7,
		pieces.King:   1,
	}
)

// Struct to pass along all global fields
type Context struct {
	MaxDepth int8
	Color    pieces.Color

	AllNodes  int64
	LeafNodes int64

	pruneTheRest bool
}

func NewContext(maxDepth int8, color pieces.Color) *Context {
	return &Context{
		MaxDepth: maxDepth,
		Color:    color,
	}
}

// Main struct to contain the representation of the board after a given move.
type Board struct {
	board         [8][8]int8
	depth         int8
	myTurn        bool
	srcCell       [2]int
	targetCell    [2]int
	capturedPiece *pieces.Piece
	isPromotion   bool

	History            []string
	Evaluation         int
	WeightedEvaluation int
	FirstMove          *Board
}

func printCell(cell [2]int) string {
	return printLocation(cell[0], cell[1])
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
	return b.Format(true)
}
func (b *Board) Format(bold bool) string {
	var sb strings.Builder

	// sb.WriteString(fmt.Sprintf("color: %v, depth: %v, myTurn: %v\n", b.color, b.depth, b.myTurn))
	sb.WriteString(fileHeader())
	sb.WriteString(hr)
	for rank := 7; rank >= 0; rank-- {
		sb.WriteString(fmt.Sprintf("%v | ", rankPrint[rank]))
		for file := 0; file < 8; file++ {
			p := pieces.Decode(b.board[rank][file])
			sb.WriteString(fmt.Sprintf("%v | ", p.Print(bold && b.targetCell == [2]int{rank, file})))
		}
		sb.WriteString(fmt.Sprintf("%v\n", rankPrint[rank]))
	}
	sb.WriteString(hr)
	sb.WriteString(fileHeader())
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("# Evaluation: %d\n", b.Evaluation))
	sb.WriteString(fmt.Sprintf("# WeightedEvaluation: %d\n", b.WeightedEvaluation))
	sb.WriteString(b.formatHistory())

	return sb.String()
}
func (b *Board) formatHistory() string {
	var sb strings.Builder
	for index, move := range b.History {
		if index%2 == 0 {
			sb.WriteString(historyPrefix + move + ", ")
		} else {
			sb.WriteString(move + "\n")
		}
	}
	return sb.String()
}

func (b *Board) DebugString() string {
	return fmt.Sprintf("Depth: %v Evalution: %v History: %v", b.depth, b.WeightedEvaluation, b.History)
}
func (b *Board) ChildNode() *Board {
	newBoard := &Board{
		board:   b.board,
		depth:   b.depth + 1,
		myTurn:  !b.myTurn,
		History: append([]string{}, b.History...),
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

func (b *Board) SetMove(ctx *Context, move string) {
	var piece *pieces.Piece
	var src, dst string

	if move == "O-O" || move == "O-O-O" {
		b.moveCastle(ctx, move)
	} else {

		if len(move) == 4 {
			src = move[0:2]
			dst = move[2:4]
			piece = b.GetPiece(src)

		} else if len(move) == 5 {
			piece = pieces.Parse(string(move[0]))
			src = move[1:3]
			dst = move[3:5]
			if piece.Class != b.GetPiece(src).Class {
				panic("Wrong piece at source location")
			}
		} else {
			panic("Unable to read move " + move)
		}

		if piece.Class == pieces.ClassNone {
			panic("Cannot move empty cell")
		}
		if piece.IsSameColor(ctx.Color) {
			panic("Cannot move your own piece manually")
		}

		oldfile := fileParse[string(src[0])]
		oldrank := rankParse[string(src[1])]
		newfile := fileParse[string(dst[0])]
		newrank := rankParse[string(dst[1])]
		b.MovePiece(ctx, oldrank, oldfile, newrank, newfile)
	}

	myTurn := b.myTurn
	b.myTurn = false
	if b.isCheck(ctx) {
		b.History[len(b.History)-1] += "+"
	}
	b.myTurn = myTurn

}

func (b *Board) moveCastle(ctx *Context, move string) {
	const srcKing = 4
	const dstNearKing = 6
	const dstFarKing = 2
	const srcNearRook = 7
	const dstNearRook = 5
	const srcFarRook = 0
	const dstFarRook = 3

	var srcRook, dstRook, dstKing int
	if move == "O-O" {
		srcRook = srcNearRook
		dstRook = dstNearRook
		dstKing = dstNearKing
	} else if move == "O-O-O" {
		srcRook = srcFarRook
		dstRook = dstFarRook
		dstKing = dstFarKing
	}

	// We are moving opponent's piece
	enemyRank := 7
	if ctx.Color == pieces.Black {
		enemyRank = 0
	}

	king := pieces.Decode(b.board[enemyRank][srcKing])
	if king.Class != pieces.King {
		panic("King is not found on " + printLocation(enemyRank, srcKing))
	}
	rook := pieces.Decode(b.board[enemyRank][srcRook])
	if rook.Class != pieces.Rook {
		panic("Rook is not found on " + printLocation(enemyRank, srcRook))
	}

	b.board[enemyRank][srcKing] = pieces.EmptyCell // remove king
	b.board[enemyRank][dstKing] = king.Encode()    // add king
	b.board[enemyRank][srcRook] = pieces.EmptyCell // remove rook
	b.board[enemyRank][dstRook] = rook.Encode()    // add rook
	b.History = append(b.History, move)
	b.srcCell = [2]int{enemyRank, dstRook}
	b.targetCell = [2]int{enemyRank, dstKing}
	return

}

func (b *Board) MovePiece(ctx *Context, oldrank, oldfile, newrank, newfile int) {
	old := &b.board[oldrank][oldfile]

	new := &b.board[newrank][newfile]
	verb := "to"
	if *new != pieces.EmptyCell {
		verb = "x"
		b.capturedPiece = pieces.Decode(*new)
	}
	*new = *old
	*old = pieces.EmptyCell

	b.srcCell = [2]int{oldrank, oldfile}
	b.targetCell = [2]int{newrank, newfile}

	p := pieces.Decode(*new)
	b.History = append(b.History, fmt.Sprintf("%v %v %v", p.Print(false), verb, printLocation(newrank, newfile)))
	// fmt.Printf("%v\n", b.String())
	// fmt.Printf("%v\n", printLocation(b.targetCell[0], b.targetCell[1]))
}

func (b *Board) isCheck(ctx *Context) bool {
	for _, move := range b.getMoves(ctx, false) {
		// fmt.Printf("captured piece %v\n", move.capturedPiece.Print(true))
		if move.capturedPiece != nil && move.capturedPiece.Class == pieces.King {
			return true
		}
	}
	return false
}

// Evaluate the final score of the board based on the pieces on the board.
func (b *Board) EvaluateMaterial(ctx *Context) {
	var sum int = 0
	var weightedSum int = 0
	kings := map[pieces.Color]bool{}

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			piece := pieces.Decode(b.board[rank][file])
			if piece.Class == pieces.King {
				kings[piece.Color] = true
			} else {
				value := piece.Value()
				weightValue := piece.ValueWeightedByLocation(rank, file)

				if piece.IsSameColor(ctx.Color) {
					sum += value
					weightedSum += weightValue
				} else {
					sum -= value
					weightedSum -= weightValue
				}
			}
		}
	}

	// Both Kings are still present - return material value
	if len(kings) == 2 {
		b.Evaluation = sum
		b.WeightedEvaluation = weightedSum
		return
	}

	// My king is dead - this is infinitely bad
	if !kings[ctx.Color] {
		b.Evaluation = -deadKing
		b.WeightedEvaluation = -deadKing
		return
	}

	// Their king is dead - this is infintely good
	b.Evaluation = deadKing
	b.WeightedEvaluation = deadKing
}

// Check if a given move is allowed.
// Some pawn moves are only allowed for capture (diagonal).
// Some pawn moves are NOT allowed for capture (forward).
func (b *Board) moveAvailable(
	rank, file int, piece *pieces.Piece, allowCapture, allowNoCapture bool) bool {
	isValid := rank >= 0 && rank < 8 && file >= 0 && file < 8
	if !isValid {
		return false
	}

	newPiece := pieces.Decode(b.board[rank][file])

	// Cell is empty
	if allowNoCapture && newPiece.Class == pieces.ClassNone {
		return true
	}

	// Capture opponent piece
	if allowCapture && newPiece.IsOppositeColor(piece.Color) {
		return true
	}

	// Move  not allowed
	return false
}

// Entry point for chessbot to begin recursion.
func (b *Board) FindMoves(ctx *Context) *Board {
	leafBoard := b.minimax(ctx, -infinity, infinity)
	myTurn := leafBoard.FirstMove.myTurn
	leafBoard.FirstMove.myTurn = true
	if leafBoard.FirstMove.isCheck(ctx) {
		leafBoard.FirstMove.History[len(leafBoard.FirstMove.History)-1] += "+"
	}
	leafBoard.FirstMove.myTurn = myTurn
	return leafBoard

}

// Recursive depth first function to step through each board in the decision tree.
// Alpha and beta pruning is used to trim unneeded nodes.
func (b *Board) minimax(ctx *Context, alpha, beta int) *Board {
	// We actually search depths 2 more than MaxDepth so that we can evalation
	// capture moves after MaxDepth
	if b.depth >= ctx.MaxDepth*2 || b.WeightedEvaluation == deadKing || b.WeightedEvaluation == -deadKing {
		ctx.LeafNodes++
		return b
	}

	// Hacky way to initialize minmax value to infinity
	minimaxBoard := &Board{}
	if b.myTurn {
		minimaxBoard.WeightedEvaluation = -infinity
	} else {
		minimaxBoard.WeightedEvaluation = infinity
	}

	// We have reached our depth limit.
	// Just search capture moves from this node down.
	allowNoCaptures := true
	if b.depth >= ctx.MaxDepth {
		allowNoCaptures = false
	}
	// fmt.Printf("Depth: %d allowNoCaptures: %v LastMove: %v\n", b.depth, allowNoCaptures, b.History[len(b.History)-1])

	for _, newBoard := range b.getMoves(ctx, allowNoCaptures) {
		// Recursively traverse downwards
		eval := newBoard.minimax(ctx, alpha, beta)

		if b.myTurn {
			// if eval.WeightedEvaluation > minimaxBoard.WeightedEvaluation ||
			// 	(eval.WeightedEvaluation == minimaxBoard.WeightedEvaluation && eval.depth < minimaxBoard.depth) {
			if eval.WeightedEvaluation > minimaxBoard.WeightedEvaluation {
				minimaxBoard = eval
			}
			if eval.WeightedEvaluation > alpha {
				alpha = eval.WeightedEvaluation
			}
		} else if !b.myTurn {
			// if eval.WeightedEvaluation < minimaxBoard.WeightedEvaluation ||
			// 	(eval.WeightedEvaluation == minimaxBoard.WeightedEvaluation && eval.depth > minimaxBoard.depth) {
			if eval.WeightedEvaluation < minimaxBoard.WeightedEvaluation {
				minimaxBoard = eval
			}
			if eval.WeightedEvaluation < beta {
				beta = eval.WeightedEvaluation
			}
		}

		if beta <= alpha {
			break
		}

	}

	// There were no eligible moves here, so just return the existing board.
	if minimaxBoard.WeightedEvaluation == infinity || minimaxBoard.WeightedEvaluation == -infinity {
		ctx.LeafNodes++
		minimaxBoard = b
	}
	return minimaxBoard
}

// Get all eligible moves on the board for the current player.
// allowNoCaptures determines if we just want capture moves or all moves.
func (b *Board) getMoves(ctx *Context, allowNoCaptures bool) []*Board {
	type scoredMoves struct {
		score int
		board *Board
	}
	sortedMoves := []scoredMoves{}
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			piece := pieces.Decode(b.board[rank][file])
			// Skip empty cells
			if piece.Class == pieces.ClassNone {
				continue
			}
			// Only look at the correct pieces for this turn
			if (piece.IsSameColor(ctx.Color)) != b.myTurn {
				continue
			}
			moves := []*Board{}

			switch piece.Class {
			case pieces.Pawn:
				moves = append(moves, b.generatePawnMoves(ctx, rank, file, piece, allowNoCaptures)...)
			case pieces.Knight:
				moves = append(moves, b.generateKnightMoves(ctx, rank, file, piece, allowNoCaptures)...)
			case pieces.Bishop, pieces.Rook, pieces.Queen, pieces.King:
				moves = append(moves, b.generateSlidingMoves(ctx, rank, file, piece, allowNoCaptures)...)
			default:
				continue
			}

			// Assign score to each move
			for _, move := range moves {
				score := 0
				// Prioritize moves where we capture a valuable piece
				// with a less valuable pieces.
				if move.capturedPiece != nil {
					score += 10*move.capturedPiece.Value() - piece.Value()
				}
				// Promoting a pawn should have the value of an extra queen.
				if move.isPromotion {
					score += pieces.Value[pieces.Queen]
				}
				// Deprioritize moving a piece that would be captured by a pawn.
				if b.getCellsOpponentPawnsCanCapture(ctx)[[2]int{rank, file}] {
					score -= piece.Value()
				}

				sortedMoves = append(sortedMoves, scoredMoves{score, move})
			}
		}
	}

	// Sort the list of moves by score and return them
	sort.SliceStable(sortedMoves, func(i, j int) bool {
		return sortedMoves[i].score > sortedMoves[j].score
	})
	moves := []*Board{}
	for _, move := range sortedMoves {
		moves = append(moves, move.board)
	}

	return moves
}

// Generate a new board given the old and new rank+file coordinates.
// Some pawn moves are only allowed for capture (diagonal).
// Some pawn moves are NOT allowed for capture (forward).
// Some pawn moves are used for promotion.
func (b *Board) generateMove(ctx *Context, oldRank, oldFile, newRank, newFile int,
	piece *pieces.Piece, allowCapture, allowNoCapture bool, promote bool) []*Board {
	if !b.moveAvailable(newRank, newFile, piece, allowCapture, allowNoCapture) {
		return []*Board{}
	}
	newBoard := b.ChildNode()
	newBoard.MovePiece(ctx, oldRank, oldFile, newRank, newFile)

	// Pawn Promotion
	if promote {
		newQueen := pieces.Piece{
			Class: pieces.Queen,
			Color: piece.Color,
		}
		newBoard.board[newRank][newFile] = newQueen.Encode()
		newBoard.isPromotion = true
	}
	newBoard.EvaluateMaterial(ctx)
	ctx.AllNodes++
	return []*Board{newBoard}
}

func (b *Board) generatePawnMoves(ctx *Context, rank, file int, piece *pieces.Piece, allowNoCaptures bool) []*Board {
	moves := []*Board{}
	var direction int = 1
	if piece.Color == pieces.Black {
		direction = -1
	}

	// Pawn diagonal - capture only
	moves = append(moves, b.generateMove(ctx, rank, file, rank+direction, file+1, piece, true, false, false)...)
	moves = append(moves, b.generateMove(ctx, rank, file, rank+direction, file-1, piece, true, false, false)...)

	// If we only are searching for capture moves, then we are done here.
	if !allowNoCaptures {
		return moves
	}

	// Pawn promotion. Always pick Queen.
	if ((rank == 1 && piece.Color == pieces.Black) || (rank == 6 && piece.Color == pieces.White)) &&
		b.board[rank+direction][file] == pieces.EmptyCell {
		moves = append(moves, b.generateMove(ctx, rank, file, rank+direction, file, piece, true, true, true)...)
	}

	// Pawn forward 2 - Init only - no capture
	// This doesn't work if the cell in front of the pawn is occupied
	// fmt.Printf("movePawn %v %v %v\n", rank+direction, file, b.board[rank+direction][file])
	if ((rank == 1 && piece.Color == pieces.White) || (rank == 6 && piece.Color == pieces.Black)) &&
		b.board[rank+direction][file] == pieces.EmptyCell {
		moves = append(moves, b.generateMove(ctx, rank, file, rank+direction*2, file, piece, false, true, false)...)
	}

	// Pawn forward one space - no capture
	// fmt.Printf("Pawn forward one space %v\n", printLocation(rank, file))
	moves = append(moves, b.generateMove(ctx, rank, file, rank+direction, file, piece, false, true, false)...)

	return moves
}

func (b *Board) generateKnightMoves(ctx *Context, rank, file int, piece *pieces.Piece, allowNoCaptures bool) []*Board {
	moves := []*Board{}
	for _, dir := range knightDirs {
		rankDir := dir[0]
		fileDir := dir[1]

		newRank := rank + rankDir
		newFile := file + fileDir

		newBoard := b.generateMove(ctx, rank, file, newRank, newFile, piece, true, allowNoCaptures, false)
		moves = append(moves, newBoard...)
	}
	return moves
}

// Bishops, Rooks, Queen, and King are really all just variations of a sliding
// piece. So we use the same logic here.
func (b *Board) generateSlidingMoves(ctx *Context, rank, file int, piece *pieces.Piece, allowNoCaptures bool) []*Board {
	moves := []*Board{}
	dirs := sliderPieceDirs[piece.Class]
	distance := sliderPieceDistance[piece.Class]

	for _, dir := range dirs {
		rankDir := dir[0]
		fileDir := dir[1]

		for i := 1; i <= distance; i++ {
			newRank := rank + rankDir*i
			newFile := file + fileDir*i

			newBoard := b.generateMove(ctx, rank, file, newRank, newFile, piece, true, true, false)
			// Piece is unable to move in this direction - No need to continue further
			if len(newBoard) == 0 {
				break
			}
			// If we only want to return capture moves and this is not a
			// capture move, then we need to skip it.
			if !allowNoCaptures && newBoard[0].capturedPiece == nil {
				continue
			}
			moves = append(moves, newBoard...)

			// Piece has captured an opponent piece along this direction.
			// No need to continue further.
			if newBoard[0].capturedPiece != nil {
				break
			}
		}
	}
	return moves
}

// Generates a list of cells that the opponents' pawns can capture.
// This is used to help order move priorities.
func (b *Board) getCellsOpponentPawnsCanCapture(ctx *Context) map[[2]int]bool {
	cells := map[[2]int]bool{} // rank and file
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			piece := pieces.Decode(b.board[rank][file])

			// Skip empty cells
			if piece.Class == pieces.ClassNone {
				continue
			}
			// Only look at opponent pawns that can capture my pieces.
			if (piece.IsSameColor(ctx.Color)) == b.myTurn {
				continue
			}

			if piece.Class != pieces.Pawn {
				continue
			}

			var direction int = 1
			if piece.Color == pieces.Black {
				direction = -1
			}

			// Pawn diagonal - capture only
			if b.moveAvailable(rank+direction, file+1, piece, true, false) {
				cells[[2]int{rank + direction, file + 1}] = true
			}
			if b.moveAvailable(rank+direction, file-1, piece, true, false) {
				cells[[2]int{rank + direction, file - 1}] = true
			}
		}
	}
	return cells
}

// Read board text file.
func ParseBoard(ctx *Context, data []byte, myTurn bool) *Board {
	b := &Board{
		depth:  0,
		myTurn: myTurn,
	}

	s := string(data)
	rows := strings.Split(s, "\n")

	for _, row := range rows {
		row = strings.TrimSpace(row)

		if strings.HasPrefix(row, historyPrefix) {
			segs := strings.Split(strings.TrimPrefix(row, historyPrefix), ",")
			for _, seg := range segs {
				seg = strings.TrimSpace(seg)
				if len(seg) > 0 {
					b.History = append(b.History, seg)
				}
			}
			continue
		}

		// Skip other comment lines
		if len(row) == 0 || row[0] == '#' {
			continue
		}
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
	b.EvaluateMaterial(ctx)
	return b
}
