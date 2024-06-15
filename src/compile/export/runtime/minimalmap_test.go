package core

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type String string

func (s String) Serialize() []byte {
	return make([]byte, 8)
}

func TestMinimalMap(t *testing.T) {
	readFile, err := os.Open("/usr/share/dict/words")
	if err != nil {
		panic(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	words := map[uint64]String{}
	line := uint64(1)
	for fileScanner.Scan() {
		words[line] = String(strings.TrimSpace(fileScanner.Text()))
		line++
	}

	mm := NewMinimalMap(words, func(a, b String) bool { return string(a) == string(b) })
	word := mm.Get(500)
	assert.Equal(t, String("abscondence"), word)
	assert.Equal(t, len(mm.redirects)*4+len(mm.values)*8+8, len(mm.Bytes()))
}
