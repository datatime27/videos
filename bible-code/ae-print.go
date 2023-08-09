package main
import (
    "flag"
    "fmt"
	"math/rand"
    "regexp"
    "os"
    "strings"
    "strconv"
)

var (
    inPath = flag.String("input", "", "file to load")
    outPath = flag.String("output", "", "file to load")
)

const (
template = `
function addColorText(layer, index, color) {
    var animator = layer.property('Text').property('Animators').addProperty('ADBE Text Animator');
    var colorProp = animator.property('Properties').addProperty('ADBE Text Fill Color');
    colorProp.setValueAtTime(0,[0,0,0]);
    colorProp.setValueAtTime(1,color);
    var selector = animator.property('Selectors').addProperty('ADBE Text Selector');
    var advanced = selector.property('Advanced');
    advanced.property('Based On').setValue(2); // Characters
    advanced.property('Units').setValue(2); // Index
    var start = selector.property('ADBE Text Index Start').setValue(index)
    var end = selector.property('ADBE Text Index End').setValue(index+1)
}


var font_size = 30;
var tracking = 700;

var positionX = 960;
var positionY = 540;

var width = 1800;
var height = 900;

var letterWidth = 22.2
var letterHeight = 36.0
var wrapLength = 22

var originX = positionX - width/2
var originY = positionY - height/2

var text = "%s";

app.beginUndoGroup('bible-code');
var comp = app.project.activeItem;

label = comp.layers.addBoxText([width,height], text);
label.threeDLayer = true

var textProp = label.property("Source Text");
label.position.setValue([positionX, positionY]);
var textDocument = textProp.value;
textDocument.fontSize = font_size;
textDocument.fillColor = [0, 0, 0];
textDocument.tracking = tracking;
textProp.setValue(textDocument);

var animator = label.property('Text').property('Animators').addProperty('ADBE Text Animator');
var colorProp = animator.property('Properties').addProperty('ADBE Text Fill Color');
colorProp.setValueAtTime(0,[0,0,0]);
colorProp.setValueAtTime(1,[0.8, 0.8, 0.8]);

%s
app.endUndoGroup();
`
)

type Discovery struct {
	name   string
	start  int
	end    int
	skip   int
}
func (d Discovery) String() string {
	return fmt.Sprintf("\"%s\": %d-%d skip:%d ", d.name, d.start, d.end, d.skip)
}

func parseIndexes(s string) (int) {
    re := regexp.MustCompile(`INTERSECTION: (\d+)\-(\d+)`)
    m := re.FindStringSubmatch(s)
    startIndex, err := strconv.Atoi(m[1])
    if err != nil {
        panic(err)
    }

    return startIndex
}

func parseDiscoveries(s string) []Discovery  {
    discoveries := []Discovery{}
    re := regexp.MustCompile(`"(.+?)": (\d+)\-(\d+) skip:(\d+)`)
    segs := re.FindAllStringSubmatch(s,-1)
    for _, seg := range segs {
        start, err := strconv.Atoi(seg[2])
        if err != nil {
            panic(err)
        }

        end, err := strconv.Atoi(seg[3])
        if err != nil {
            panic(err)
        }

        skip, err := strconv.Atoi(seg[4])
        if err != nil {
            panic(err)
        }

        d := Discovery {
            name: seg[1],
            start: start,
            end: end,
            skip: skip,
        }
        discoveries = append(discoveries, d)
    }
    return discoveries
}

func randColor() string {
		return fmt.Sprintf("[%g, %g, %g]",rand.Float32(), rand.Float32()/2, rand.Float32())
}
func main() {
    flag.Parse()

    dat, err := os.ReadFile(*inPath)
    if err != nil {
        panic(err)
    }
    lines := strings.Split(string(dat), "\n")
    text := string(lines[1])
    
    startIndex := parseIndexes(lines[0])
    discoveries := parseDiscoveries(lines[0])
    
    var colorText strings.Builder
    for _, d := range discoveries {
        colorText.WriteString(fmt.Sprintf("\n// %v\n", d))
        color := randColor()
        for i:=d.start; i <= d.end; i+= d.skip {
            offset := i - startIndex
            colorText.WriteString(fmt.Sprintf("addColorText(label, %d, %s)  //%s\n", offset, color, text[offset:offset+1]))
        }
    }
    
    err = os.WriteFile(*outPath, []byte(fmt.Sprintf(template, text, colorText.String())), 0777)
    if err != nil {
        panic(err)
    }
}

