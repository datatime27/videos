package pieces

import (
	"strings"
)

type Class int8
type Color int8

const (
	EmptyCell = 0

	ClassNone Class = 0
	Pawn            = 1
	Knight          = 2
	Bishop          = 3
	Rook            = 4
	Queen           = 5
	King            = 6

	ColorNone Color = 0
	White           = 1
	Black           = 2

	BoldText    = "\033[1m\033[4m"
	DefaultText = "\033[0m"
)

var (
	parse = map[string]Class{
		" ": ClassNone,
		"K": King,
		"Q": Queen,
		"R": Rook,
		"B": Bishop,
		"N": Knight,
		"P": Pawn,
	}
	fullName = map[string]string{
		"K": "King",
		"Q": "Queen",
		"R": "Rook",
		"B": "Bishop",
		"N": "Knight",
		"P": "Pawn",
	}

	print = map[Class]string{
		ClassNone: " ",
		King:      "K",
		Queen:     "Q",
		Rook:      "R",
		Bishop:    "B",
		Knight:    "N",
		Pawn:      "P",
	}

	Value = map[Class]int{
		ClassNone: 0,
		Pawn:      10,
		Knight:    30,
		Bishop:    30,
		Rook:      50,
		Queen:     90,
	}
)

func FullName(letter string) string {
	return fullName[letter]
}

type Piece struct {
	Class Class
	Color Color
}

func (p *Piece) Print(bold bool) string {
	s := print[p.Class]
	if p.Color == Black {
		s = strings.ToLower(s)
	}
	if bold {
		s = BoldText + s + DefaultText
	}
	return s
}

func (p *Piece) Value() int {
	return Value[p.Class]
}

// This provides evaluation of the board by including weighted values for
// pieces by class.
func (p *Piece) ValueWeightedByLocation(rank, file int) int {
	const pawnMultiplier = 1
	const kingMultiplier = 5
	const mainMultiplier = 10

	value := Value[p.Class]
	rankWeight := 0
	fileWeight := 0

	switch p.Class {
	// Encourage Pawns to move forward
	case Pawn:
		if p.Color == White {
			rankWeight = rank - 1
		} else {
			rankWeight = 6 - rank
		}
		rankWeight *= pawnMultiplier

	// Discourage Kings from moving forward
	case King:
		if p.Color == White {
			rankWeight = 7 - rank
		} else {
			rankWeight = rank
		}
		rankWeight *= kingMultiplier

	// Encourage other pieces towards the center of the board
	case Knight, Bishop, Rook, Queen:
		if rank <= 3 {
			rankWeight = rank
		} else {
			rankWeight = 7 - rank
		}
		rankWeight *= mainMultiplier

		if file <= 3 {
			fileWeight = file
		} else {
			fileWeight = 7 - file
		}
		fileWeight *= mainMultiplier
	}
	return value + rankWeight + fileWeight
}

func (p *Piece) Encode() int8 {
	return int8(p.Color)<<4 + int8(p.Class)
}

func Decode(code int8) *Piece {
	return &Piece{
		Color: Color(code >> 4),
		Class: Class(code & 15),
	}
}

func (p *Piece) IsSameColor(color Color) bool {
	if p.Color == White && color == White {
		return true
	}
	if p.Color == Black && color == Black {
		return true
	}
	return false
}

func (p *Piece) IsOppositeColor(color Color) bool {
	if p.Color == White && color == Black {
		return true
	}
	if p.Color == Black && color == White {
		return true
	}
	return false
}
func OppositeColor(color Color) Color {
	if color == Black {
		return White
	}
	if color == White {
		return Black
	}
	return ColorNone
}

func Parse(notation string) *Piece {
	if notation == "" {
		return &Piece{
			Class: ClassNone,
			Color: ColorNone,
		}
	}
	char := strings.ToUpper(notation)
	side := White
	if char != notation {
		side = Black
	}
	return &Piece{
		Class: parse[char],
		Color: Color(side),
	}
}
