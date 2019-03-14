package abi

import (
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// batch of predefined reflect types
var (
	boolT         = reflect.TypeOf(bool(false))
	uint8T        = reflect.TypeOf(uint8(0))
	uint16T       = reflect.TypeOf(uint16(0))
	uint32T       = reflect.TypeOf(uint32(0))
	uint64T       = reflect.TypeOf(uint64(0))
	int8T         = reflect.TypeOf(int8(0))
	int16T        = reflect.TypeOf(int16(0))
	int32T        = reflect.TypeOf(int32(0))
	int64T        = reflect.TypeOf(int64(0))
	addressT      = reflect.TypeOf(common.Address{})
	stringT       = reflect.TypeOf("")
	dynamicBytesT = reflect.SliceOf(reflect.TypeOf(byte(0)))
	functionT     = reflect.ArrayOf(24, reflect.TypeOf(byte(0)))
	tupleT        = reflect.TypeOf(map[string]interface{}{})
	bigIntT       = reflect.TypeOf(new(big.Int))
)

// Kind represents the kind of abi type
type Kind int

const (
	// KindBool is a boolean
	KindBool Kind = iota

	// KindUInt is an uint
	KindUInt

	// KindInt is an int
	KindInt

	// KindString is a string
	KindString

	// KindArray is an array
	KindArray

	// KindSlice is a slice
	KindSlice

	// KindAddress is an address
	KindAddress

	// KindBytes is a bytes array
	KindBytes

	// KindFixedBytes is a fixed bytes
	KindFixedBytes

	// KindFixedPoint is a fixed point
	KindFixedPoint

	// KindTuple is a tuple
	KindTuple

	// KindFunction is a function
	KindFunction
)

func (k Kind) String() string {
	names := [...]string{
		"Bool",
		"Uint",
		"Int",
		"String",
		"Array",
		"Slice",
		"Address",
		"Bytes",
		"FixedBytes",
		"FixedPoint",
		"Tuple",
		"Function",
	}

	return names[k]
}

// TupleElem is an element of a tuple
type TupleElem struct {
	Name string
	Elem *Type
}

// Type is an ABI type
type Type struct {
	kind  Kind
	size  int
	elem  *Type
	raw   string
	tuple []*TupleElem
	t     reflect.Type
}

func (t *Type) isVariableInput() bool {
	return t.kind == KindSlice || t.kind == KindBytes || t.kind == KindString
}

func (t *Type) isDynamicType() bool {
	if t.kind == KindTuple {
		for _, elem := range t.tuple {
			if elem.Elem.isDynamicType() {
				return true
			}
		}
		return false
	}
	return t.kind == KindString || t.kind == KindBytes || t.kind == KindSlice || (t.kind == KindArray && t.elem.isDynamicType())
}

func parseType(arg *Argument) (string, error) {
	if !strings.HasPrefix(arg.Type, "tuple") {
		return arg.Type, nil
	}

	if len(arg.Components) == 0 {
		return "", fmt.Errorf("tuple type expects components but none found")
	}

	// parse the arg components from the tuple
	str := []string{}
	for _, i := range arg.Components {
		aux, err := parseType(i)
		if err != nil {
			return "", err
		}
		str = append(str, i.Name+" "+aux)
	}
	return fmt.Sprintf("tuple(%s)%s", strings.Join(str, ","), strings.TrimPrefix(arg.Type, "tuple")), nil
}

// NewTypeFromArgument parses an abi type from an argument
func NewTypeFromArgument(arg *Argument) (*Type, error) {
	str, err := parseType(arg)
	if err != nil {
		return nil, err
	}
	return NewType(str)
}

// NewType parses a type in string format
func NewType(s string) (*Type, error) {
	l := newLexer(s)
	return readType(l)
}

func getTypeSize(t *Type) int {
	if t.kind == KindArray && !t.elem.isDynamicType() {
		if t.elem.kind == KindArray || t.elem.kind == KindTuple {
			return t.size * getTypeSize(t.elem)
		}
		return t.size * 32
	} else if t.kind == KindTuple && !t.isDynamicType() {
		total := 0
		for _, elem := range t.tuple {
			total += getTypeSize(elem.Elem)
		}
		return total
	}
	return 32
}

var typeRegexp = regexp.MustCompile("^([[:alpha:]]+)([[:digit:]]*)$")

func expectedToken(t tokenType) error {
	return fmt.Errorf("expected token %s", t.String())
}

func readType(l *lexer) (*Type, error) {
	var tt *Type

	tok := l.nextToken()
	if tok.typ == tupleToken {
		if l.nextToken().typ != lparenToken {
			return nil, expectedToken(lparenToken)
		}

		var next token
		elems := []*TupleElem{}
		for {

			// read name
			name := l.nextToken()
			if name.typ != strToken {
				return nil, expectedToken(strToken)
			}

			elem, err := readType(l)
			if err != nil {
				return nil, err
			}
			elems = append(elems, &TupleElem{
				Name: name.literal,
				Elem: elem,
			})

			next = l.nextToken()
			if next.typ == commaToken {
				continue
			} else {
				break
			}
		}

		rawAux := []string{}
		for _, i := range elems {
			rawAux = append(rawAux, i.Elem.raw)
		}
		raw := fmt.Sprintf("(%s)", strings.Join(rawAux, ","))

		tt = &Type{kind: KindTuple, raw: raw, tuple: elems, t: tupleT}

	} else if tok.typ != strToken {
		return nil, expectedToken(strToken)

	} else {
		// Check normal types
		elem, err := decodeSimpleType(tok.literal)
		if err != nil {
			return nil, err
		}
		tt = elem
	}

	// check for arrays at the end of the type
	for {
		if l.ch != '[' {
			break
		}

		l.nextToken()
		n := l.nextToken()

		var tAux *Type
		if n.typ == rbracketToken {
			tAux = &Type{kind: KindSlice, elem: tt, raw: fmt.Sprintf("%s[]", tt.raw), t: reflect.SliceOf(tt.t)}

		} else if n.typ == numberToken {
			size, err := strconv.ParseUint(n.literal, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("failed to read array size '%s': %v", n.literal, err)
			}

			tAux = &Type{kind: KindArray, elem: tt, raw: fmt.Sprintf("%s[%d]", tt.raw, size), size: int(size), t: reflect.ArrayOf(int(size), tt.t)}
			if l.nextToken().typ != rbracketToken {
				return nil, expectedToken(rbracketToken)
			}
		} else {
			return nil, fmt.Errorf("unexpected token %s", n.typ.String())
		}

		tt = tAux
	}
	return tt, nil
}

func decodeSimpleType(str string) (*Type, error) {
	match := typeRegexp.FindStringSubmatch(str)
	if len(match) == 0 {
		return nil, fmt.Errorf("type format is incorrect. Expected 'type''bytes' but found '%s'", str)
	}
	match = match[1:]

	var err error
	t := match[0]

	bytes := 0
	ok := false

	if bytesStr := match[1]; bytesStr != "" {
		bytes, err = strconv.Atoi(bytesStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bytes '%s': %v", bytesStr, err)
		}
		ok = true
	}

	// Only int and uint need bytes for sure, 'bytes' may
	// have or not, the rest dont have bytes
	if t == "int" || t == "uint" {
		if !ok {
			return nil, fmt.Errorf("int and uint expect bytes")
		}
	} else if t != "bytes" && ok {
		return nil, fmt.Errorf("type %s does not expect bytes", t)
	}

	switch t {
	case "uint":
		var k reflect.Type
		switch bytes {
		case 8:
			k = uint8T
		case 16:
			k = uint16T
		case 32:
			k = uint32T
		case 64:
			k = uint64T
		default:
			if bytes%8 != 0 {
				panic(fmt.Errorf("number of bytes has to be M mod 8"))
			}
			k = bigIntT
		}
		return &Type{kind: KindUInt, size: int(bytes), t: k, raw: fmt.Sprintf("uint%d", bytes)}, nil

	case "int":
		var k reflect.Type
		switch bytes {
		case 8:
			k = int8T
		case 16:
			k = int16T
		case 32:
			k = int32T
		case 64:
			k = int64T
		default:
			if bytes%8 != 0 {
				panic(fmt.Errorf("number of bytes has to be M mod 8"))
			}
			k = bigIntT
		}
		return &Type{kind: KindInt, size: int(bytes), t: k, raw: fmt.Sprintf("int%d", bytes)}, nil

	case "byte":
		bytes = 1
		fallthrough

	case "bytes":
		if bytes == 0 {
			return &Type{kind: KindBytes, t: dynamicBytesT, raw: "bytes"}, nil
		}
		return &Type{kind: KindFixedBytes, size: int(bytes), raw: fmt.Sprintf("bytes%d", bytes), t: reflect.ArrayOf(int(bytes), reflect.TypeOf(byte(0)))}, nil

	case "string":
		return &Type{kind: KindString, t: stringT, raw: "string"}, nil

	case "bool":
		return &Type{kind: KindBool, t: boolT, raw: "bool"}, nil

	case "address":
		return &Type{kind: KindAddress, t: addressT, raw: "address"}, nil

	case "function":
		return &Type{kind: KindFunction, size: 24, t: functionT, raw: "function"}, nil

	default:
		return nil, fmt.Errorf("unknown type '%s'", t)
	}
}

type tokenType int

const (
	eofToken tokenType = iota
	strToken
	numberToken
	tupleToken
	lparenToken
	rparenToken
	lbracketToken
	rbracketToken
	commaToken
	invalidToken
)

func (t tokenType) String() string {
	names := [...]string{
		"eof",
		"string",
		"number",
		"tuple",
		"(",
		")",
		"[",
		"]",
		",",
		"<invalid>",
	}
	return names[t]
}

type token struct {
	typ     tokenType
	literal string
}

type lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func newLexer(input string) *lexer {
	l := &lexer{input: input}
	l.readChar()
	return l
}

func (l *lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *lexer) nextToken() token {
	var tok token

	// skip whitespace
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}

	switch l.ch {
	case ',':
		tok.typ = commaToken
	case '(':
		tok.typ = lparenToken
	case ')':
		tok.typ = rparenToken
	case '[':
		tok.typ = lbracketToken
	case ']':
		tok.typ = rbracketToken
	case 0:
		tok.typ = eofToken
	default:
		if isLetter(l.ch) {
			tok.literal = l.readIdentifier()
			if tok.literal == "tuple" {
				tok.typ = tupleToken
			} else {
				tok.typ = strToken
			}

			return tok
		} else if isDigit(l.ch) {
			return token{numberToken, l.readNumber()}
		} else {
			tok.typ = invalidToken
		}
	}

	l.readChar()
	return tok
}

func (l *lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}