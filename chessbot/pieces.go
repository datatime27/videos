package pieces

import (
	"strings"
)

const (
	None   = 0
	Pawn   = 1
	Knight = 2
	Bishop = 3
	Rook   = 4
	Queen  = 5
	King   = 6

	White = 1
	Black = 2
)

var (
	parse = map[string]int8{
		" ": None,
		"K": King,
		"Q": Queen,
		"R": Rook,
		"B": Bishop,
		"N": Knight,
		"P": Pawn,
	}
	print = map[int8]string{
		None:   " ",
		King:   "K",
		Queen:  "Q",
		Rook:   "R",
		Bishop: "B",
		Knight: "N",
		Pawn:   "P",
	}

	value = map[int8]int{
		None:   0,
		Queen:  9,
		Rook:   5,
		Bishop: 3,
		Knight: 3,
		Pawn:   1,
	}
)

type Piece struct {
	Piece int8
	Color int8
}

func (p *Piece) Print() string {
	s := print[p.Piece]
	if p.Color == Black {
		s = strings.ToLower(s)
	}
	return s
}

func (p *Piece) Value() int {
	return value[p.Piece]
}

func (p *Piece) Encode() int8 {
	return int8(p.Color<<4 + p.Piece)
}

func Decode(code int8) *Piece {
	return &Piece{
		Color: code >> 4,
		Piece: code & 15,
	}
}

func (p *Piece) IsSameColor(color int8) bool {
	if p.Color == White && color == White {
		return true
	}
	if p.Color == Black && color == Black {
		return true
	}
	return false
}

func (p *Piece) IsOppositeColor(color int8) bool {
	if p.Color == White && color == Black {
		return true
	}
	if p.Color == Black && color == White {
		return true
	}
	return false
}

func Parse(notation string) *Piece {
	if notation == "" {
		return &Piece{
			Piece: None,
			Color: None,
		}
	}
	char := strings.ToUpper(notation)
	side := White
	if char != notation {
		side = Black
	}
	return &Piece{
		Piece: parse[char],
		Color: int8(side),
	}
}
