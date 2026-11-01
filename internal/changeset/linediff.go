package changeset

import "strings"

// LineOp marks a line of a line diff.
type LineOp byte

// The three marks of a line diff line.
const (
	LineSame    LineOp = ' '
	LineRemoved LineOp = '-'
	LineAdded   LineOp = '+'
)

// Line is one line of a line diff.
type Line struct {
	Op   LineOp
	Text string
}

// LineDiff compares old and new line by line, from a longest common
// subsequence: lines in both are kept, the rest are removed (-) or added
// (+). A removed run comes before the added run that replaces it.
func LineDiff(old, new string) []Line {
	a, b := strings.Split(old, "\n"), strings.Split(new, "\n")
	// lcs[i][j] is the common length of a[i:] and b[j:].
	lcs := make([][]int, len(a)+1)
	for i := range lcs {
		lcs[i] = make([]int, len(b)+1)
	}
	for i := len(a) - 1; i >= 0; i-- {
		for j := len(b) - 1; j >= 0; j-- {
			if a[i] == b[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else {
				lcs[i][j] = max(lcs[i+1][j], lcs[i][j+1])
			}
		}
	}
	var out []Line
	i, j := 0, 0
	for i < len(a) || j < len(b) {
		switch {
		case i < len(a) && j < len(b) && a[i] == b[j]:
			out = append(out, Line{LineSame, a[i]})
			i++
			j++
		case i < len(a) && (j == len(b) || lcs[i+1][j] >= lcs[i][j+1]):
			out = append(out, Line{LineRemoved, a[i]})
			i++
		default:
			out = append(out, Line{LineAdded, b[j]})
			j++
		}
	}
	return out
}
