package boards

import (
	"os"

	// "fmt"
	"html/template"
	// "strings"
	"chessbot/lib/pieces"
)

var (
	htmlTemplate = template.Must(template.New("html").Parse(`
<html>
<head>
<style>
td {
    text-align: center;
    width: 26px;
    font-size: 24px;
    font-weight: bold;
}
.border {
    border-style:solid;
    border-color: red;
    border-width: 5px;
    width: 70px;
    height: 70px;
}

</style>
</head>
<table cellpadding="0" cellspacing="0" style="background:white;border:1px #c8ccd1 solid;padding:0;margin:auto">
<tbody>
<tr><td></td><td>a</td><td>b</td><td>c</td><td>d</td><td>e</td><td>f</td><td>g</td><td>h</td></tr>
<tr>
<td>8</td>
<td colspan="8" rowspan="8">
<div style="position:relative">
<img src="icons/Chessboard480.svg.png" />

{{range $item := .Items}}
<div style="position:absolute;top:{{$item.Top}}px;left:{{$item.Left}}px;"><img class="{{$item.Border}}" src="{{$item.Src}}" /></div>{{end}}
<td>8</td>
</tr>

<tr><td>7</td><td>7</td></tr>
<tr><td>6</td><td>6</td></tr>
<tr><td>5</td><td>5</td></tr>
<tr><td>4</td><td>4</td></tr>
<tr><td>3</td><td>3</td></tr>
<tr><td>2</td><td>2</td></tr>
<tr><td>1</td><td>1</td></tr>
<tr>
<td></td>
<td>a</td>
<td>b</td>
<td>c</td>
<td>d</td>
<td>e</td>
<td>f</td>
<td>g</td>
<td>h</td>
</tr>
</tbody>
</table>
</html>
`))

	srcLookup = map[pieces.Color]map[pieces.Class]string{
		pieces.White: map[pieces.Class]string{
			pieces.Rook:   "icons/Chess_rlt45.svg.png",
			pieces.Knight: "icons/Chess_nlt45.svg.png",
			pieces.Bishop: "icons/Chess_blt45.svg.png",
			pieces.Queen:  "icons/Chess_qlt45.svg.png",
			pieces.King:   "icons/Chess_klt45.svg.png",
			pieces.Pawn:   "icons/Chess_plt45.svg.png",
		},
		pieces.Black: map[pieces.Class]string{
			pieces.Rook:   "icons/Chess_rdt45.svg.png",
			pieces.Knight: "icons/Chess_ndt45.svg.png",
			pieces.Bishop: "icons/Chess_bdt45.svg.png",
			pieces.Queen:  "icons/Chess_qdt45.svg.png",
			pieces.King:   "icons/Chess_kdt45.svg.png",
			pieces.Pawn:   "icons/Chess_pdt45.svg.png",
		},
	}
)

type Config struct {
	Items []Item
}
type Item struct {
	Top    int
	Left   int
	Src    string
	Border string
}

func (b *Board) WriteHTML(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	config := Config{}
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			p := pieces.Decode(b.board[rank][file])
			if p.Class == pieces.ClassNone {
				continue
			}

			item := Item{
				Top:  640 - (rank+1)*80,
				Left: file * 80,
				Src:  srcLookup[p.Color][p.Class],
			}
			if b.targetCell == [2]int{rank, file} {
				item.Border = "border"
			}

			config.Items = append(config.Items, item)
		}
	}
	return htmlTemplate.Execute(f, config)
}
