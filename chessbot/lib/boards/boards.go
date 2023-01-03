package boards

import (
	"fmt"
	"sort"
	"strings"

	"chessbot/lib/pieces"
)

const (
	hr = "-------------------------------------\n"

	infinity = 100000
	deadKing = 1000
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

type Board struct {
	board         [8][8]int8
	depth         int8
	myTurn        bool
	targetCell    [2]int
	capturedPiece *pieces.Piece
	isPromotion   bool

	History            string
	Evaluation         int
	WeightedEvaluation int
	FirstMove          *Board
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
	sb.WriteString(fmt.Sprintf("# History: %v\n", b.History))

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
	verb := "to"
	if *new != pieces.EmptyCell {
		verb = "x"
		b.capturedPiece = pieces.Decode(*new)
	}
	*new = *old
	*old = pieces.EmptyCell

	b.targetCell = [2]int{newrank, newfile}

	p := pieces.Decode(*new)
	b.History += fmt.Sprintf("%v %v %v, ", p.Print(false), verb, printLocation(newrank, newfile))
}

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

/*
func (b *Board) spawnDecision(ctx *Context, oldRank, oldFile, newRank, newFile int,
	piece *pieces.Piece, allowTake, allowNoTake bool, alpha, beta int8, leafBoard **Board) bool {
	if !b.moveAvailable(newRank, newFile, piece, allowTake, allowNoTake) {
		return false
	}
	newBoard := b.ChildNode()
	newBoard.MovePiece(ctx, oldRank, oldFile, newRank, newFile)
	newBoard.EvaluateMaterial(ctx)
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
*/

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

func (b *Board) FindMoves(ctx *Context) *Board {
	return b.minimax(ctx, -infinity, infinity)
}

func (b *Board) minimax(ctx *Context, alpha, beta int) *Board {
	if b.depth >= ctx.MaxDepth+2 || b.WeightedEvaluation == deadKing || b.WeightedEvaluation == -deadKing {
		//fmt.Printf("Eval %v %v %v\n", b.depth, b.Evaluation, b.History)
		ctx.LeafNodes++
		return b
	}

	minimaxBoard := &Board{History: b.History}
	if b.myTurn {
		minimaxBoard.WeightedEvaluation = -infinity
	} else {
		minimaxBoard.WeightedEvaluation = infinity
	}

	allowNoCaptures := true
	if b.depth >= ctx.MaxDepth {
		// We have reached our depth limit.
		// Just search capture moves from this node down.
		allowNoCaptures = false
		// 	fmt.Printf("captureOnly Depth = %v History = %v\n", b.depth, b.History)
	}
	// fmt.Printf("captureOnly = %v Depth = %v\n", captureOnly, b.depth)

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
		// fmt.Printf("alpha = %v, beta = %v\n", alpha, beta)

		if beta <= alpha {
			break
			// fmt.Printf("Prune alpha %v\n", minimaxBoard.DebugString())
		}

	}

	if minimaxBoard.WeightedEvaluation == infinity || minimaxBoard.WeightedEvaluation == -infinity {
		ctx.LeafNodes++
		minimaxBoard = b
	}
	// fmt.Printf("minimaxBoard.FirstMove == nil %v Depth %d\n", minimaxBoard.FirstMove == nil, minimaxBoard.depth)
	return minimaxBoard
}

/*
func (b *Board) minimax1(ctx *Context, alpha, beta int8) *Board {
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
			if (piece.IsSameColor(ctx.Color)) != b.myTurn {
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
			// if ctx.pruneTheRest {
			// 	fmt.Printf("Pruning %v %v\n", b.DebugString(), leafBoard.DebugString())
			// 	ctx.pruneTheRest = false
			// 	return leafBoard
			// }
		}
	}
	if leafBoard.Evaluation == infinity-1 || leafBoard.Evaluation == -infinity+1 {
		// fmt.Printf("No moves found for %v %v %v\n", leafBoard.Evaluation, b.Evaluation, b.History)
		leafBoard = b
	}
	return leafBoard
}
*/

func (b *Board) getMoves(ctx *Context, allowNoCaptures bool) []*Board {
	// if !allowNoCaptures {
	// 	return []*Board{}
	// }
	type scoredMoves struct {
		score int
		board *Board
	}
	sortedMoves := []scoredMoves{}
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			piece := pieces.Decode(b.board[rank][file])
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
				// Prioritize moves where we capture most valuable pieces with a our least pieces.
				if move.capturedPiece != nil {
					score += 10*move.capturedPiece.Value() - piece.Value()
				}
				// Promoting a pawn should have the value of an extra queen.
				if move.isPromotion {
					score += pieces.Value[pieces.Queen]
				}
				// Deprioritize moving a piece that would be captured by a pawn.
				if b.getCellsOppoentPawnsCanCapture(ctx)[[2]int{rank, file}] {
					score -= piece.Value()
				}

				sortedMoves = append(sortedMoves, scoredMoves{score, move})
			}

			// switch piece.Class {
			// case pieces.Pawn:
			// 	moves[pieces.Pawn] = append(moves[pieces.Pawn], b.generatePawnMoves(ctx, rank, file, piece)...)
			// case pieces.Knight:
			// 	moves[pieces.Knight] = append(moves[pieces.Knight], b.generateKnightMoves(ctx, rank, file, piece)...)
			// case pieces.Bishop:
			// 	moves[pieces.Bishop] = append(moves[pieces.Bishop], b.generateSlidingMoves(ctx, rank, file, piece)...)
			// case pieces.Rook:
			// 	moves[pieces.Rook] = append(moves[pieces.Rook], b.generateSlidingMoves(ctx, rank, file, piece)...)
			// case pieces.Queen:
			// 	moves[pieces.Queen] = append(moves[pieces.Queen], b.generateSlidingMoves(ctx, rank, file, piece)...)
			// case pieces.King:
			// 	moves[pieces.King] = append(moves[pieces.King], b.generateSlidingMoves(ctx, rank, file, piece)...)
			// default:
			// 	continue
			// }
		}
	}

	// Prioritize moves from least valuable piece to most valuable.
	// allMoves := []*Board{}
	// for _, class := range []pieces.Class{
	// 	pieces.Pawn,
	// 	pieces.Knight,
	// 	pieces.Bishop,
	// 	pieces.Rook,
	// 	pieces.Queen,
	// 	pieces.King} {
	// 	allMoves = append(allMoves, moves[class]...)
	// }
	// return allMoves

	// Sort the list of moves by score
	sort.SliceStable(sortedMoves, func(i, j int) bool {
		// return false
		return sortedMoves[i].score > sortedMoves[j].score
	})
	moves := []*Board{}
	for _, move := range sortedMoves {
		moves = append(moves, move.board)
	}

	return moves
}

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

/*
func (b *Board) movePawn(ctx *Context, rank, file int, piece *pieces.Piece, alpha, beta int8, leafBoard **Board) {
	var direction int = 1
	if piece.Color == pieces.Black {
		direction = -1
	}
	// Pawn forward one space - no capture
	b.spawnDecision(ctx, rank, file, rank+direction, file, piece, false, true, alpha, beta, leafBoard)
	if ctx.pruneTheRest {
		return
	}

	// Pawn diagonal - capture only
	b.spawnDecision(ctx, rank, file, rank+direction, file+1, piece, true, false, alpha, beta, leafBoard)
	if ctx.pruneTheRest {
		return
	}
	b.spawnDecision(ctx, rank, file, rank+direction, file-1, piece, true, false, alpha, beta, leafBoard)
	if ctx.pruneTheRest {
		return
	}

	// Pawn forward 2 - Init only - no capture
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
*/
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

/*
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
*/

func (b *Board) getCellsOppoentPawnsCanCapture(ctx *Context) map[[2]int]bool {
	cells := map[[2]int]bool{} // rank and file
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			piece := pieces.Decode(b.board[rank][file])
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

func ParseBoard(ctx *Context, data []byte, myTurn bool) *Board {
	b := &Board{
		depth:  0,
		myTurn: myTurn,
	}

	s := string(data)
	rows := strings.Split(s, "\n")

	for _, row := range rows {
		row = strings.TrimSpace(row)

		// Skip comments
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
