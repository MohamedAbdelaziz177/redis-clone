package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type RespType byte

const (
	STRING  RespType = '+'
	ERROR   RespType = '-'
	INTEGER RespType = ':'
	BULK    RespType = '$'
	ARRAY   RespType = '*'
)

type Value struct {
	Type  RespType
	Str   string
	Bulk  string
	Int   int
	Array []Value
	Err   string
}

type Parser struct {
	reader *bufio.Reader
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{reader: bufio.NewReader(reader)}
}

func (p *Parser) ParseValue() (Value, error) {
	respType, err := p.getRespType()
	if err != nil {
		return Value{}, err
	}

	switch respType {
	case STRING:
		return p.ParseString()
	case ERROR:
		return p.ParseError()
	case INTEGER:
		return p.ParseInteger()
	case BULK:
		return p.ParseBulk()
	case ARRAY:
		return p.ParseArray()
	default:
		return Value{}, fmt.Errorf("unknown RESP type: %c", respType)
	}
}

func (p *Parser) getRespType() (RespType, error) {
	b, err := p.reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return RespType(b), nil
}

func (p *Parser) ParseString() (Value, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}
	return Value{Type: STRING, Str: strings.TrimSuffix(line, "\r\n")}, nil
}

func (p *Parser) ParseError() (Value, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}
	return Value{Type: ERROR, Err: strings.TrimSuffix(line, "\r\n")}, nil
}

func (p *Parser) ParseInteger() (Value, error) {

	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	intStr := strings.TrimSuffix(line, "\r\n")

	intVal, err := strconv.Atoi(intStr)
	if err != nil {
		return Value{}, err
	}

	return Value{Type: INTEGER, Int: intVal}, nil
}

func (p *Parser) ParseBulk() (Value, error) {

	sizeStr, err := p.reader.ReadString('\n')

	if err != nil {
		return Value{}, err
	}

	sizeStr = strings.TrimSuffix(sizeStr, "\r\n")

	sz, err := strconv.Atoi(sizeStr)
	if err != nil {
		return Value{}, err
	}

	if sz < 0 {
		return Value{Type: BULK, Bulk: ""}, nil
	}

	buff := make([]byte, sz)

	_, err = io.ReadFull(p.reader, buff)
	if err != nil {
		return Value{}, err
	}

	p.reader.ReadByte()
	p.reader.ReadByte()

	return Value{Type: BULK, Bulk: string(buff)}, nil
}

func (p *Parser) ParseArray() (Value, error) {

	sizeStr, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	sz, err := strconv.Atoi(strings.TrimSuffix(sizeStr, "\r\n"))
	if err != nil {
		return Value{}, err
	}
	if sz < 0 {
		return Value{Type: ARRAY, Array: []Value{}}, nil
	}

	tokens := make([]Value, 0)

	for i := 0; i < sz; i++ {
		val, err := p.ParseValue()
		if err != nil {
			return Value{}, err
		}
		tokens = append(tokens, val)
	}

	return Value{Type: ARRAY, Array: tokens}, nil
}

/*
func (p *Parser) ParseArrayOfBulks() (Value, error) {

	sizeVal, err := p.ParseInteger()
	if err != nil {
		return Value{}, err
	}

	sz := sizeVal.Int

	tonkens := make([]Value, 0, sz)

	for i := 0; i < sz; i++ {

		if _, err := p.reader.ReadByte(); err != nil {
			return Value{}, err
		}

		b, err := p.ParseBulk()
		if err != nil {
			return Value{}, err
		}
		tonkens = append(tonkens, b)
	}

	return Value{Type: ARRAY, Array: tonkens}, nil
}
*/
