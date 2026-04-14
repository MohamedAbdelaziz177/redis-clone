package main

import (
	"bufio"
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

func NewParser(reader *bufio.Reader) *Parser {
	return &Parser{reader: reader}
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

	sizeVal, err := p.ParseInteger()
	if err != nil {
		return Value{}, err
	}

	sz := sizeVal.Int

	if sz < 0 {
		return Value{Type: BULK, Bulk: ""}, nil
	}

	buff := make([]byte, sz+2)

	_, err = io.ReadFull(p.reader, buff)
	if err != nil {
		return Value{}, err
	}

	return Value{Type: BULK, Bulk: strings.TrimSuffix(string(buff[:sz]), "\r\n")}, nil
}
