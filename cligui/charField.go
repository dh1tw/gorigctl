package cligui

import "strings"

import ui "github.com/gizak/termui"

type CharField struct {
	ui.Block
	Text       string
	TextColor  ui.Attribute
	WrapLength int
	Alignment  ui.Align
	PaddingTop int
}

func NewCharField(s string) *CharField {
	return &CharField{
		Block:      *ui.NewBlock(),
		Text:       s,
		TextColor:  ui.ThemeAttr("par.text.fg"),
		WrapLength: 0,
		PaddingTop: 1,
	}
}

func (d *CharField) Buffer() ui.Buffer {
	buf := d.Block.Buffer()

	// fg, bg := d.TextFgColor, d.TextBgColor
	// cs := DefaultDecimalBuilder.Build(d.Text, fg, bg)

	cField := make([][][]rune, 0, len(d.Text))

	for _, r := range d.Text {
		fontChar, ok := smallFont[r]
		if ok {
			cField = append(cField, fontChar)
		}
	}

	// truncate if charField is larger than canvas
	for getCharFieldWidth(cField) > d.InnerWidth()-1 {
		cField = cField[:len(cField)-2]
		cField = append(cField, smallFont['…'])
	}

	length := getCharFieldWidth(cField)

	x, ux := 0, 0
	y := d.PaddingTop

	switch d.Alignment {
	case ui.AlignLeft:
		x = 0
	case ui.AlignRight:
		x = d.InnerWidth() - length
	default:
		x = 0
	}

	ux = x
	for _, char := range cField {
		ux, _ = d.drawChar(ux, y, char, &buf)
	}

	return buf
}

func (d *CharField) drawChar(x, y int, char [][]rune, buf *ui.Buffer) (ux, uy int) {

	uy = y

	cell := ui.Cell{
		Bg: d.TextColor,
		Ch: ' ',
	}

	for _, line := range char {
		ux = x
		for _, c := range line {
			if c == '#' {
				buf.Set(d.InnerX()+ux, d.InnerY()+uy, cell)
			}
			ux++
		}
		uy++
	}

	return ux, uy
}

var colon string = `
..
#.
..
#.
..
`

var point string = `
..
..
..
..
#.
`

var dotdot string = `
....
....
....
....
#.#.
`

var space string = `
.......
.......
.......
.......
.......
`

var zero string = `
######.
#....#.
#....#.
#....#.
######.
`
var one string = `
....#.
....#.
....#.
....#.
....#.
`
var two string = `
######.
.....#.
######.
#......
######.
`

var three string = `
######.
.....#.
...###.
.....#.
######.
`

var four string = `
#......
#......
#...#..
######.
....#..
`

var five string = `
######.
#......
######.
.....#.
######.
`

var six string = `
######.
#......
######.
#....#.
######.
`

var seven string = `
######.
.....#.
.....#.
.....#.
.....#.
`

var height string = `
######.
#....#.
######.
#....#.
######.
`

var nine string = `
######.
#....#.
######.
.....#.
######.
`

var letterR string = `
#####..
#....#.
#####..
#...#..
#....#.
`

var letterA string = `
######.
#....#.
######.
#... #.
#....#.
`

var letterD string = `
#####.
#....#.
#....#.
#... #.
#####.
`

var letterI string = `
.#.
.#.
.#.
.#.
.#.
`

var letterO string = `
######.
#....#.
#....#.
#....#.
######.
`

var letterL string = `
#......
#......
#......
#......
######.
`

var letterN string = `
#....#.
##...#.
#.#..#.
#..#.#.
#...##.
`

var letterE string = `
######.
#......
####...
#......
######.
`

var letterF string = `
######.
#......
####...
#......
#......
`

// smallFont defines the font use to display the timer on termbox
var smallFont = map[rune][][]rune{
	'…': asArray(dotdot),
	';': asArray(colon),
	'.': asArray(point),
	' ': asArray(space),
	'1': asArray(one),
	'2': asArray(two),
	'3': asArray(three),
	'4': asArray(four),
	'5': asArray(five),
	'6': asArray(six),
	'7': asArray(seven),
	'8': asArray(height),
	'9': asArray(nine),
	'0': asArray(zero),
	'a': asArray(letterA),
	'A': asArray(letterA),
	'd': asArray(letterD),
	'D': asArray(letterD),
	'i': asArray(letterI),
	'I': asArray(letterI),
	'e': asArray(letterE),
	'E': asArray(letterE),
	'o': asArray(letterO),
	'O': asArray(letterO),
	'f': asArray(letterF),
	'F': asArray(letterF),
	'l': asArray(letterL),
	'L': asArray(letterL),
	'n': asArray(letterN),
	'N': asArray(letterN),
	'r': asArray(letterR),
	'R': asArray(letterR),
}

// Convert a character as an array of rune
func asArray(chars string) [][]rune {
	result := [][]rune{}
	line := []rune{}
	str := strings.TrimPrefix(chars, "\n")
	for _, c := range str {
		if c == '\n' {
			result = append(result, line)
			line = []rune{}
		} else {
			line = append(line, c)
		}
	}
	return result
}

// Return the width of a ascii char field
func getCharFieldWidth(cf [][][]rune) int {

	length := 0

	// iterate over all char fields
	for _, char := range cf {

		// all lines in the character array are equal long
		// so it's ok to check just the first line
		length += len(char[0])

	}

	return length
}
