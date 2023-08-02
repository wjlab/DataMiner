package core

import (
	"fmt"
	"github.com/fatih/color"
)

func Logo() {
	logo := "                                                      ,,                              \n" +
		"`7MM\"\"\"Yb.            mm            `7MMM.     ,MMF'  db                              \n" +
		"  MM    `Yb.          MM              MMMb    dPMM                                    \n" +
		"  MM     `Mb  ,6\"Yb.mmMMmm  ,6\"Yb.    M YM   ,M MM  `7MM  `7MMpMMMb.  .gP\"Ya `7Mb,od8 \n" +
		"  MM      MM 8)   MM  MM   8)   MM    M  Mb  M' MM    MM    MM    MM ,M'   Yb  MM' \"' \n" +
		"  MM     ,MP  ,pm9MM  MM    ,pm9MM    M  YM.P'  MM    MM    MM    MM 8M\"\"\"\"\"\"  MM     \n" +
		"  MM    ,dP' 8M   MM  MM   8M   MM    M  `YM'   MM    MM    MM    MM YM.    ,  MM     \n" +
		".JMMmmmdP'   `Moo9^Yo.`Mbmo`Moo9^Yo..JML. `'  .JMML..JMML..JMML  JMML.`Mbmmd'.JMML.  "

	magenta := color.New(color.FgMagenta)
	hiMagenta := color.New(color.FgHiMagenta)
	cyan := color.New(color.FgCyan)

	colorIndex := 0
	for _, line := range splitLines(logo) {
		switch colorIndex {
		case 0, 3, 5, 6:
			magenta.Println(line)
		case 2:
			hiMagenta.Println(line)
		default:
			cyan.Println(line)
		}
		colorIndex++
	}
	color.Unset()
	fmt.Println()
}

// splitLines splits a string into lines and returns a slice of lines
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}