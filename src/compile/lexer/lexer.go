package lexer

import (
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/io/path"
	"strconv"
	"unicode"
	"unicode/utf8"

	"go.uber.org/zap"
)

type Lexer struct {
	path          path.Path
	input         string
	curRuneIndex  int  // current Location in input (points to current char)
	nextRuneIndex int  // current reading Location in input (after current char)
	ch            rune // current char under examination

	col  int
	line int
}

func NewLexer(input string, path path.Path) *Lexer {
	return &Lexer{
		input: input,
		path:  path,
		col:   -1,
	}
}

func (l *Lexer) readChar() {
	size := 1
	if l.nextRuneIndex >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL"
	} else {
		l.ch, size = utf8.DecodeRuneInString(l.input[l.nextRuneIndex:])
	}

	l.col++
	if l.ch == '\n' {
		l.col = -1
		l.line++
	}

	l.curRuneIndex = l.nextRuneIndex
	l.nextRuneIndex += size
}

func (l *Lexer) peekChar() rune {
	if l.nextRuneIndex >= len(l.input) {
		return 0
	}

	r, _ := utf8.DecodeRuneInString(l.input[l.nextRuneIndex:])
	return r
}

func (l *Lexer) nextToken() token.Token {
	l.skipWhitespace()

	startLine := l.line
	startColumn := l.col
	var tok token.Token
	switch l.ch {
	case 0:
		tok = token.Token{Type: token.TOK_EOF}
	case '"':
		tok = l.readString()
	case '\'':
		tok = l.readCharLiteral()
	default:
		if unicode.IsDigit(l.ch) {
			tok = l.readNumber()
		} else if unicode.IsLetter(l.ch) || l.ch == '_' {
			tok = l.readIdentifier()
		} else {
			tok = l.readSymbol()
		}
	}

	tok.Loc = token.Location{
		Path_:       l.path,
		StartLine:   startLine,
		EndLine:     l.line,
		StartColumn: startColumn,
		EndColumn:   l.col,
	}

	return tok
}

// readString captures and returns a string literal token, handling escaped characters and ensuring bounds.
func (l *Lexer) readString() token.Token {
	startPos := l.curRuneIndex
	valid := true
	var out string
	for {
		l.readChar()
		// End of string or file
		if l.ch == '"' || l.ch == 0 {
			break
		} else if l.ch == '\n' {
			valid = false
			l.readChar()
			break
		}
		if l.ch == '\\' {
			l.readChar() // Read the escape character
			switch l.ch {
			case 'n':
				out += "\n"
			case 'r':
				out += "\r"
			case 't':
				out += "\t"
			case 'b':
				out += "\b"
			case 'f':
				out += "\f"
			case '"':
				out += "\""
			case '\\':
				out += "\\"
			case 'u':
				if l.nextRuneIndex+4 > len(l.input) {
					// Not enough characters for a valid Unicode escape, break or handle error
					l.curRuneIndex = len(l.input) - 1 // Skip past the Unicode sequence.
					l.nextRuneIndex = len(l.input)
					break
				}
				// Attempt to parse Unicode escape sequence
				hex := l.input[l.curRuneIndex+1 : l.curRuneIndex+5] // Next 4 chars after \u
				unicodeVal, err := strconv.ParseInt(hex, 16, 32)
				if err != nil {
					// Handle invalid Unicode sequence
					valid = false
				} else {
					out += string(rune(unicodeVal))
				}
				l.readChar()
				l.readChar()
				l.readChar()
				l.readChar()
			default:
				break
			}
		} else {
			out += string(l.ch)
		}
	}

	if l.ch == '"' {
		l.readChar()
		if valid {
			return token.Token{Type: token.TOK_STRING, Value: out}
		}
	}

	return token.Token{Type: token.TOK_UNKOWN, Value: l.input[startPos:l.curRuneIndex]}
}

// readCharLiteral captures and returns a character literal token, handling escaped characters.
func (l *Lexer) readCharLiteral() token.Token {
	startPos := l.curRuneIndex
	valid := true
	l.readChar() // Move past the opening quote.
	var out rune
	if l.ch == '\\' { // Check for escape character.
		l.readChar() // Read the escape character or Unicode sequence.
		switch l.ch {
		case 'n':
			out = '\n'
		case 'r':
			out = '\r'
		case 't':
			out = '\t'
		case 'b':
			out = '\b'
		case 'f':
			out = '\f'
		case '\'':
			out = '\''
		case '\\':
			out = '\\'
		case 'u':
			if l.nextRuneIndex+4 > len(l.input) {
				// Not enough characters for a valid Unicode escape, handle as needed.
				l.curRuneIndex = len(l.input) - 1 // Skip past the Unicode sequence.
				l.nextRuneIndex = len(l.input)
			} else {
				// Attempt to parse Unicode escape sequence.
				hex := l.input[l.curRuneIndex+1 : l.curRuneIndex+5] // Next 4 chars after \u.
				unicodeVal, err := strconv.ParseInt(hex, 16, 32)
				if err != nil {
					valid = false
				} else {
					out = rune(unicodeVal)
				}
				l.readChar()
				l.readChar()
				l.readChar()
				l.readChar()
			}
		default:
			valid = false
		}
	} else if l.ch == '\n' {
		valid = false
	} else {
		out = l.ch
	}
	l.readChar() // Move past the character or escape sequence.

	if l.ch == '\'' {
		l.readChar()
		if valid {
			return token.Token{Type: token.TOK_CHARACTER, Rune: out, Value: string(out)}
		}
	}

	return token.Token{Type: token.TOK_UNKOWN, Value: l.input[startPos:l.curRuneIndex]}
}

func (l *Lexer) skipWhitespace() {
	for {
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			l.readChar()
		} else if l.ch == '/' {
			if l.peekChar() == '/' {
				for l.ch != '\n' && l.ch != 0 {
					l.readChar()
				}
				l.readChar()
			} else if l.peekChar() == '*' {
				for !(l.ch == '*' && l.peekChar() == '/') && l.ch != 0 {
					l.readChar()
				}
				l.readChar()
				l.readChar()
			} else {
				return
			}
		} else {
			return
		}
	}
}

func isHexDigit(r rune) bool {
	return unicode.IsDigit(r) || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')
}

func (l *Lexer) readNumber() token.Token {
	if l.ch == '0' {
		// Check for hexadecimal or binary
		peekChar := l.peekChar()
		if peekChar == 'x' || peekChar == 'X' {
			l.readChar() // Consume '0'
			l.readChar() // Consume 'x' or 'X'
			return l.readHexNumber()
		} else if peekChar == 'b' || peekChar == 'B' {
			// Binary
			l.readChar() // Consume '0'
			l.readChar() // Consume 'b' or 'B'
			return l.readBinaryNumber()
		}
	}

	return l.readDecimalNumber()
}

func (l *Lexer) readDecimalNumber() token.Token {
	Location := l.curRuneIndex
	tokenType := token.TOK_INTEGER
	for unicode.IsDigit(l.ch) || l.ch == '.' || l.ch == 'e' || l.ch == 'E' || l.ch == '+' || l.ch == '-' {
		if l.ch == '.' {
			if tokenType == token.TOK_FLOAT { // Break if we encounter a second decimal point
				break
			}
			tokenType = token.TOK_FLOAT
		} else if l.ch == 'e' || l.ch == 'E' {
			tokenType = token.TOK_FLOAT
			// Look ahead to check for exponent sign
			if peekChar := l.peekChar(); peekChar == '+' || peekChar == '-' || unicode.IsDigit(peekChar) {
				l.readChar() // Consume 'e' or 'E'
				if l.ch == '+' || l.ch == '-' {
					l.readChar() // Consume sign
				}
			} else {
				break // Break if 'e' or 'E' is not followed by a sign or digit
			}
		}
		l.readChar()
	}

	if tokenType == token.TOK_INTEGER {
		value, err := strconv.ParseUint(l.input[Location:l.curRuneIndex], 10, 64)
		if err == nil {
			return token.Token{Type: tokenType, Integer: value, Value: l.input[Location:l.curRuneIndex]}
		}
	} else {
		value, err := strconv.ParseFloat(l.input[Location:l.curRuneIndex], 64)
		if err == nil {
			return token.Token{Type: tokenType, Float: value, Value: l.input[Location:l.curRuneIndex]}
		}
	}
	return token.Token{Type: token.TOK_UNKOWN, Value: l.input[Location:l.curRuneIndex]}
}

func (l *Lexer) readHexNumber() token.Token {
	tokenType := token.TOK_INTEGER
	location := l.curRuneIndex
	hasP := false
	for isHexDigit(l.ch) || l.ch == '.' || (l.ch == 'p' || l.ch == 'P') {
		if l.ch == '.' {
			if tokenType == token.TOK_FLOAT { // Break if we encounter a second decimal point
				break
			}
			tokenType = token.TOK_FLOAT
		} else if l.ch == 'p' || l.ch == 'P' {
			if hasP { // Break if we encounter a second 'p' or 'P'
				break
			}
			hasP = true
			tokenType = token.TOK_FLOAT
			// Look ahead to check for exponent sign
			if peekChar := l.peekChar(); peekChar == '+' || peekChar == '-' || unicode.IsDigit(peekChar) {
				l.readChar() // Consume 'p' or 'P'
				if l.ch == '+' || l.ch == '-' {
					l.readChar() // Consume sign
				}
			} else {
				break // Break if 'p' or 'P' is not followed by a sign or digit
			}
		}
		l.readChar()
	}

	if tokenType == token.TOK_INTEGER {
		value, err := strconv.ParseUint(l.input[location:l.curRuneIndex], 16, 64)
		if err == nil {
			return token.Token{Type: tokenType, Integer: value, Value: l.input[location-2 : l.curRuneIndex]}
		}
	} else {
		value, err := strconv.ParseFloat(l.input[location-2:l.curRuneIndex], 64)
		if err == nil {
			return token.Token{Type: tokenType, Float: value, Value: l.input[location-2 : l.curRuneIndex]}
		}
	}
	return token.Token{Type: token.TOK_UNKOWN, Value: l.input[location-2 : l.curRuneIndex]}
}

func (l *Lexer) readBinaryNumber() token.Token {
	Location := l.curRuneIndex
	for l.ch == '0' || l.ch == '1' {
		l.readChar()
	}

	value, err := strconv.ParseUint(l.input[Location:l.curRuneIndex], 2, 64)
	if err != nil {
		return token.Token{Type: token.TOK_UNKOWN, Value: l.input[Location-2 : l.curRuneIndex]}
	}

	return token.Token{Type: token.TOK_INTEGER, Integer: value, Value: l.input[Location-2 : l.curRuneIndex]}
}

func (l *Lexer) readIdentifier() token.Token {
	Location := l.curRuneIndex
	l.readChar()

	for unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}

	identifier := l.input[Location:l.curRuneIndex]
	if _, isKeyword := token.KeywordMap[token.Type(identifier)]; isKeyword {
		return token.Token{Type: token.Type(identifier)}
	}

	return token.Token{Type: token.TOK_IDENTIFIER, Value: identifier}
}

func (l *Lexer) readSymbol() token.Token {
	startLocation := l.curRuneIndex

	for {
		symbol := l.input[startLocation:l.nextRuneIndex]
		l.readChar()
		if _, isSymbol := token.SymbolMap[token.Type(symbol)]; isSymbol {
			if l.nextRuneIndex > len(l.input) {
				return token.Token{Type: token.Type(symbol)}
			}

			nextSymbol := l.input[startLocation:l.nextRuneIndex]
			if _, okNext := token.SymbolMap[token.Type(nextSymbol)]; !okNext {
				// If extending doesn't lead to a longer match, return the current match
				return token.Token{Type: token.Type(symbol)}
			}
		} else {
			return token.Token{Type: token.TOK_UNKOWN, Value: l.input[startLocation:l.nextRuneIndex]}
		}
	}
}

func (l *Lexer) Lex() ([]token.Token, io.Error) {
	l.readChar()

	var tokens []token.Token
	var err io.Error
	for {
		t := l.nextToken()
		tokens = append(tokens, t)
		switch t.Type {
		case token.TOK_EOF:
			return tokens, err
		case token.TOK_UNKOWN:
			err = io.JoinError(err, io.NewError("malformed token", zap.String("token", t.String())))
		}
	}
}
