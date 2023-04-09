package boards

import (
	"fmt"
	"html/template"
	"net/url"
	"os"
	"strings"

	"chessbot/lib/pieces"
)

var (
	htmlTemplate = template.Must(template.New("html").Parse(`
<html>
<head>
<style>
.board {
	background:white;
	border:1px #c8ccd1 solid;
	padding:0;
	margin:auto;
	margin-right: 20px;
}
.board tr td{
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
.history td {
    font-size: 24px;
}
div{
    font-size: 24px;
}
input {
    font-size: 24px;
}

</style>
<script>
var cellwidth=80;
const fileLookup = new Map();
fileLookup.set(0, "a");
fileLookup.set(1, "b");
fileLookup.set(2, "c");
fileLookup.set(3, "d");
fileLookup.set(4, "e");
fileLookup.set(5, "f");
fileLookup.set(6, "g");
fileLookup.set(7, "h");

function pieceClick(piece,location){
    moveoutput = document.getElementById("moveoutput")
    if (moveoutput.value.length == 0){ 
    		moveoutput.value = piece+location;
    } else {
    		moveoutput.value += location;
    }
}

function emptyClick(event){
    rank = 8 - Math.floor((event.y-event.target.getBoundingClientRect().top)/cellwidth)
    file = Math.floor((event.x-event.target.getBoundingClientRect().left)/cellwidth)
    
    moveoutput = document.getElementById("moveoutput")
    moveoutput.value += fileLookup.get(file)+rank
}
</script>

</head>
<body>
<table class="board" align="left" cellpadding="0" cellspacing="0" >
<tbody>
<tr><td></td><td>a</td><td>b</td><td>c</td><td>d</td><td>e</td><td>f</td><td>g</td><td>h</td></tr>
<tr>
<td>8</td>
<td colspan="8" rowspan="8">
<div style="position:relative">
<img src="icons/Chessboard480.svg.png" onclick="emptyClick(event)"/>

{{range $item := .Items}}
<div style="position:absolute;top:{{$item.Top}}px;left:{{$item.Left}}px;"><img onclick="pieceClick('{{$item.Piece}}','{{$item.Location}}')" class="{{$item.Border}}" src="{{$item.Src}}" /></div>{{end}}
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
FEN: {{.Fen}}<br>
FEN Link: <a target="fen" href="{{.FenLink}}">{{.FenLink}}</a><br>
<div>
    <h3>move output:</h3>
    <input type=text id="moveoutput" value="" size=8/>
    <input type=button value="copy to clipboard" onclick="navigator.clipboard.writeText(moveoutput.value);"/>
    <input type=button value="clear" onclick="moveoutput.value=''"/>
</div>
<div>Evaluation: {{.Evaluation}}</div>
<div>WeightedEvaluation: {{.WeightedEvaluation}}</div>
<h3>History:</h3>
<table border=1 class="history">
{{range $row := .History}}
<tr>
{{range $cell := $row}}
<td>{{$cell}}</td>{{end}}
</tr>{{end}}
</table>
</body>
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
	Items              []Item
	History            [][2]string
	Evaluation         int
	WeightedEvaluation int
	Fen                string
	FenLink            string
}
type Item struct {
	Location string
	Piece    string
	Top      int
	Left     int
	Src      string
	Border   string
}

func (b *Board) WriteHTML(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	history := [][2]string{}
	for index, move := range b.History {
		if index%2 == 0 {
			history = append(history, [2]string{move, ""})
		} else {
			history[len(history)-1][1] = move
		}
	}

	fen := b.formatFen()

	config := Config{
		Evaluation:         b.Evaluation,
		WeightedEvaluation: b.WeightedEvaluation,
		History:            history,
		Fen:                fen,
		FenLink:            b.formatFenLink(fen),
	}
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			p := pieces.Decode(b.board[rank][file])
			if p.Class == pieces.ClassNone {
				continue
			}

			item := Item{
				Location: printLocation(rank, file),
				Piece:    p.Print(false),
				Top:      640 - (rank+1)*80,
				Left:     file * 80,
				Src:      srcLookup[p.Color][p.Class],
			}
			if b.targetCell == [2]int{rank, file} {
				item.Border = "border"
			}

			config.Items = append(config.Items, item)
		}
	}
	return htmlTemplate.Execute(f, config)
}

func (b *Board) formatFenLink(fen string) string {
	baseURL := "https://lichess1.org/export/fen.gif"
	v := url.Values{}
	v.Set("fen", fen)
	v.Set("variant", "standard")
	v.Set("theme", "brown")
	v.Set("piece", "cburnett")
	if b.srcCell != b.targetCell {
		v.Set("lastMove", printCell(b.srcCell)+printCell(b.targetCell))
	}
	return baseURL + "?" + v.Encode()
}

func (b *Board) formatFen() string {
	var sb strings.Builder
	numBlankSquares := 0
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			p := pieces.Decode(b.board[rank][file])
			if p.Class == pieces.ClassNone {
				numBlankSquares++
				continue
			}
			if numBlankSquares > 0 {
				sb.WriteString(fmt.Sprintf("%d", numBlankSquares))
				numBlankSquares = 0
			}
			sb.WriteString(p.Print(false))

		}
		if numBlankSquares > 0 {
			sb.WriteString(fmt.Sprintf("%d", numBlankSquares))
			numBlankSquares = 0
		}
		sb.WriteString("/")
	}
	return strings.Trim(sb.String(), "/")
}
