package namecase

import (
	"unicode"
)

type runeKind uint8

const (
	_ runeKind = iota
	lowerRuneCase
	upperRuneCase
	splitRuneCase
	ohterRuneCase
	eofRuneCase
)

func getKind(r rune) runeKind {
	switch {
	case r == 0:
		return eofRuneCase
	case unicode.IsUpper(r):
		return upperRuneCase
	case unicode.IsLower(r):
		return lowerRuneCase
	case unicode.IsSpace(r), unicode.IsPunct(r):
		return splitRuneCase
	default:
		return ohterRuneCase
	}
}

type wordKind uint8

const (
	_ wordKind = iota
	lowerWordCase
	upperWordCase
	splitWordCase
	ohterWordCase
	titleWordCase
	initialismsWordCase
)
