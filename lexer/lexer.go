package lexer

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
)

///////////////////////////////////////////////////////////////////////////////

type Position lexer.Position

///////////////////////////////////////////////////////////////////////////////

const (
	EOF rune = -(iota + 1)
	Number
	String
	Op1Number
	Op2Number
	Op1Name
	Op2Name
	Assignment
	NameLpar
	Name
	Comma
	Rpar
	Keyword
)

var symbols = map[string]rune{
	"EOF":        EOF,
	"Number":     Number,
	"String":     String,
	"Op1Number":  Op1Number,
	"Op2Number":  Op2Number,
	"Op1Name":    Op1Name,
	"Op2Name":    Op2Name,
	"Assignment": Assignment,
	"NameLpar":   NameLpar,
	"Name":       Name,
	"Comma":      Comma,
	"Rpar":       Rpar,
	"Keyword":    Keyword,
}

///////////////////////////////////////////////////////////////////////////////

type BachLexerDefinition struct {
}

func (*BachLexerDefinition) Lex(r io.Reader) lexer.Lexer {
	return &BachLexer{
		reader: bufio.NewReader(r),
		pos: lexer.Position{
			Filename: lexer.NameOfReader(r),
			Line:     1,
			Column:   1,
		},
	}
}

func (*BachLexerDefinition) Symbols() map[string]rune {
	return symbols
}

///////////////////////////////////////////////////////////////////////////////

type BachLexer struct {
	reader *bufio.Reader
	pos    lexer.Position
}

func (l *BachLexer) Next() lexer.Token {
	for {
		r, size, err := l.reader.ReadRune()
		if err == io.EOF {
			return lexer.Token{
				Type:  EOF,
				Value: "",
				Pos:   l.pos,
			}
		}
		if err != nil {
			panic(err)
		}
		if r == ' ' {
			l.pos.Offset += size
			l.pos.Column += 1
			continue
		}
		if r == '\n' {
			l.pos.Offset += size
			l.pos.Line += 1
			l.pos.Column = 1
			continue
		}
		if unicode.IsLetter(r) || r == '_' {
			l.reader.UnreadRune()
			return l.nextWithLetter()
		}
		if strings.ContainsRune("0123456789", r) {
			l.reader.UnreadRune()
			return l.nextWithDigit()
		}
		if r == '"' {
			l.reader.UnreadRune()
			return l.nextWithDoubleQuote()
		}
		if strings.ContainsRune("+-*/%<>=", r) {
			l.reader.UnreadRune()
			return l.nextWithSpecial()
		}
		if r == ',' {
			token := lexer.Token{
				Type:  Comma,
				Value: ",",
				Pos:   l.pos,
			}
			l.pos.Offset += size
			l.pos.Column += 1
			return token
		}
		if r == ')' {
			token := lexer.Token{
				Type:  Rpar,
				Value: ")",
				Pos:   l.pos,
			}
			l.pos.Offset += size
			l.pos.Column += 1
			return token
		}
		panic(errors.E("lexer", l.pos, "invalid character"))
	}
}

func (l *BachLexer) nextWithLetter() lexer.Token {
	startPos := l.pos
	builder := strings.Builder{}
	for {
		r, size, err := l.reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if unicode.IsLetter(r) || r == '_' ||
			strings.ContainsRune("0123456789", r) {
			builder.WriteRune(r)
			l.pos.Offset += size
			l.pos.Column += 1
		} else {
			break
		}
	}
	l.reader.UnreadRune()
	name := builder.String()
	if isKeyword(name) {
		return lexer.Token{
			Type: Keyword,
			Value: name,
			Pos: startPos,
		}
	}
	return lexer.Token{
		Type: Name,
		Value: name,
		Pos: startPos,
	}
}

///////////////////////////////////////////////////////////////////////////////

func isKeyword(name string) bool {
	return name == "for" || name == "def" || name == "as" || name == "ok"
}

///////////////////////////////////////////////////////////////////////////////
