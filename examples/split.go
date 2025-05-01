package examples

import (
	"strings"

	"github.com/deadlysurgeon/testfinder/extra"
)

// Split slices s into all substrings separated by sep and
// returns a slice of the substrings between those separators.
func Split(spl extra.ID, sep string) []string {
	s := string(spl)
	var result []string
	i := strings.Index(s, sep)
	for i > -1 {
		result = append(result, s[:i])
		s = s[i+len(sep):]
		i = strings.Index(s, sep)
	}
	return append(result, s)
}
