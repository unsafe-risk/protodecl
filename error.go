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

	if iType == ITYPE_TAB {
		return strings.Repeat("\t", indent)
	}
	return strings.Repeat(" ", indent)
}

func ErrorPrint(err error, file string) string {
	lines := strings.Split(file, "\n")
	var linesToPrint []string
	switch e := err.(type) {
	case *LexerError:
		Line := e.Line
		Col := e.Col
		maxlineNumSize := 0
		var pLines []string = make([]string, 0, 5)
		for i := Line - 2; i <= Line+2; i++ {
			if i < len(lines) {
				pLines = append(pLines, lines[i])
				lnsize := len(strconv.Itoa(i))
				if lnsize > maxlineNumSize {
					maxlineNumSize = lnsize
				}
			}
		}
		indent := lineIndent(pLines)
		for i := Line - 2; i <= Line+2; i++ {
			if i < len(lines) {
				lines[i] = strings.TrimPrefix(lines[i], indent)
			}
		}

		for i := Line - 2; i <= Line+2; i++ {
			if i < len(lines) {
				linesToPrint = append(linesToPrint, fmt.Sprintf("%d%s| %s", i+1, strings.Repeat(" ", maxlineNumSize-len(strconv.Itoa(i))), lines[i]))
			}
			if i == Line-1 {
				linesToPrint = append(linesToPrint, strings.Repeat(" ", maxlineNumSize)+fmt.Sprintf("| %s", strings.Repeat(" ", Col-1)+"^"))
				linesToPrint = append(linesToPrint, strings.Repeat(" ", maxlineNumSize)+"| "+e.Error())
			}
		}
	case *ParserError:
		Line := e.Position.Line
		Col := e.Position.Col
		size := len(e.Tokens[e.Pos].Value)
		if size <= 0 {
			size = 1
		}
		maxlineNumSize := 0
		var pLines []string = make([]string, 0, 5)
		for i := Line - 2; i <= Line+2; i++ {
			if i < len(lines) {
				pLines = append(pLines, lines[i])
				lnsize := len(strconv.Itoa(i))
				if lnsize > maxlineNumSize {
					maxlineNumSize = lnsize
				}
			}
		}
		indent := lineIndent(pLines)
		for i := Line - 2; i <= Line+2; i++ {
			if i < len(lines) {
				lines[i] = strings.TrimPrefix(lines[i], indent)
			}
		}

		for i := Line - 2; i <= Line+2; i++ {
			if i < len(lines) {
				linesToPrint = append(linesToPrint, fmt.Sprintf("%d%s| %s", i+1, strings.Repeat(" ", maxlineNumSize-len(strconv.Itoa(i))), lines[i]))
			}
			if i == Line-1 {
				linesToPrint = append(linesToPrint, strings.Repeat(" ", maxlineNumSize)+fmt.Sprintf("| %s", strings.Repeat(" ", Col-1)+strings.Repeat("^", size)))
				linesToPrint = append(linesToPrint, strings.Repeat(" ", maxlineNumSize)+"| "+e.Error())
			}
		}
	default:
		linesToPrint = append(linesToPrint, err.Error())
	}

	return strings.Join(linesToPrint, "\n")
}
