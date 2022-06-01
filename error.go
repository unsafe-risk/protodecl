package main

import (
	"fmt"
	"strconv"
	"strings"
)

func ErrorPrint(err error, file string) string {
	lines := strings.Split(file, "\n")
	var linesToPrint []string
	switch e := err.(type) {
	case *LexerError:
		Line := e.Line
		Col := e.Col
		maxlineNumSize := 0
		for i := Line - 2; i <= Line+2; i++ {
			if i < len(lines) {
				lnsize := len(strconv.Itoa(i))
				if lnsize > maxlineNumSize {
					maxlineNumSize = lnsize
				}
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
		for i := Line - 2; i <= Line+2; i++ {
			if i < len(lines) {
				lnsize := len(strconv.Itoa(i))
				if lnsize > maxlineNumSize {
					maxlineNumSize = lnsize
				}
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
