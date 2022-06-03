package main

import (
	"fmt"
	"strconv"
	"strings"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lineIndent(lines []string) string {
	var indents []int = make([]int, len(lines))
	const (
		ITYPE_NONE uint8 = iota
		ITYPE_TAB
		ITYPE_SPACE
	)
	itab := 0
	ispace := 0
	max := 0
	for line := range lines {
		if strings.HasPrefix(lines[line], " ") {
			ispace++
			spaces := 0
			for _, c := range lines[line] {
				if c == ' ' {
					spaces++
				} else {
					break
				}
			}
			indents[line] = spaces
			if spaces > max {
				max = spaces
			}
		} else if strings.HasPrefix(lines[line], "\t") {
			itab++
			tabs := 0
			for _, c := range lines[line] {
				if c == '\t' {
					tabs++
				} else {
					break
				}
			}
			indents[line] = tabs
			if tabs > max {
				max = tabs
			}
		}
	}

	iType := ITYPE_TAB
	if itab < ispace {
		iType = ITYPE_SPACE
	}

	indent := 0
	for line := range indents {
		d := gcd(indents[line], indent)
		indent = d
	}

	if indent == 0 {
		return ""
	}

	unitlen := max
	for line := range indents {
		if indents[line]/indent < unitlen && indents[line]/indent > 0 {
			unitlen = indents[line] / indent
		}
	}

	if iType == ITYPE_TAB {
		return strings.Repeat("\t", unitlen*indent)
	}
	return strings.Repeat(" ", unitlen*indent)
}

func CodeError(code []string, line, col, size int, filename, msg string) string {
	LineNSize := len(strconv.Itoa(line + 2))
	if size < 1 {
		size = 1
	}
	var pb strings.Builder

	var lines []string = make([]string, 0, 5)
	for i := line - 2; i <= line+2; i++ {
		if i >= 0 && i < len(code) {
			lines = append(lines, code[i])
		}
	}

	indent := lineIndent(lines)
	for i := line - 2; i <= line+2; i++ {
		if i < 0 || i >= len(code) {
			continue
		}

		fmt.Fprintf(&pb, "%d%s| %s\n",
			i+1,
			strings.Repeat(" ", LineNSize-len(strconv.Itoa(i+1))),
			strings.TrimPrefix(code[i], indent),
		)

		if i+1 == line {
			ocol := col
			if strings.HasPrefix(code[i], indent) {
				col -= len(indent)
			}
			if col > 0 {
				fmt.Fprintf(&pb, "%s| %s%s\n",
					strings.Repeat(" ", LineNSize),
					strings.Repeat(" ", col-1),
					strings.Repeat("^", size),
				)

				fmt.Fprintf(&pb, "%s| %s\n",
					strings.Repeat(" ", LineNSize),
					fmt.Sprintf("%s:%d:%d: %s", filename, line, ocol, msg),
				)
			}
		}
	}
	return pb.String()
}

func ErrorPrint(err error, file string) string {
	lines := strings.Split(file, "\n")
	var linesToPrint []string
	switch e := err.(type) {
	case *LexerError:
		Line := e.Line
		Col := e.Col
		return CodeError(lines, Line, Col, 1, e.Filename, e.Message)
	case *ParserError:
		Line := e.Position.Line
		Col := e.Position.Col
		size := len(e.Tokens[e.Pos].Value)
		return CodeError(lines, Line, Col, size, e.Position.File, e.Message)
	default:
		linesToPrint = append(linesToPrint, err.Error())
	}

	return strings.Join(linesToPrint, "\n")
}
