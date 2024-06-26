package stringers

import (
	"fmt"
)

func Join[T fmt.Stringer](stringers []T, spacer string) string {
	var result string
	var sspacer string
	for _, s := range stringers {
		result += sspacer + s.String()
		sspacer = spacer
	}

	return fmt.Sprintf("[%s]", result)
}
