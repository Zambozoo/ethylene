package strs

import (
	"fmt"
)

func Strings[T fmt.Stringer](stringers []T) string {
	var result, spacer string
	for _, s := range stringers {
		result += spacer + s.String()
		spacer = ","
	}

	return fmt.Sprintf("[%s]", result)
}
