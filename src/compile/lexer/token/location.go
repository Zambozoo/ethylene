package token

import (
	"fmt"
	"geth-cody/io/path"
)

type Location struct {
	Path        path.Path
	StartLine   int
	EndLine     int
	StartColumn int
	EndColumn   int
}

func (l *Location) Location() Location {
	return *l
}

func (l *Location) String() string {
	return fmt.Sprintf("Location{Path:%s, StartLine:%d, StartColumn:%d, EndLine:%d, EndColumn:%d}", l.Path, l.StartLine, l.StartColumn, l.EndLine, l.EndColumn)
}

type Locatable interface {
	Location() Location
}

func LocationBetween(start, end Locatable) Location {
	startLocation := start.Location()
	endLocation := end.Location()
	return Location{
		StartLine:   startLocation.StartLine,
		StartColumn: startLocation.StartColumn,
		EndLine:     endLocation.EndLine,
		EndColumn:   endLocation.EndColumn,
	}
}
