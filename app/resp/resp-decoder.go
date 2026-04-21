package resp

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
		return p.parseString()
	case ERROR:
		return p.parseError()
	case INTEGER:
		return p.parseInteger()
	case BULK:
		return p.parseBulk()
	case ARRAY:
		return p.parseArray()
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

func (p *Parser) parseString() (Value, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}
	return Value{Type: STRING, Str: strings.TrimSuffix(line, CRLF)}, nil
}

func (p *Parser) parseError() (Value, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}
	return Value{Type: ERROR, Err: strings.TrimSuffix(line, CRLF)}, nil
}

func (p *Parser) parseInteger() (Value, error) {

	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	intStr := strings.TrimSuffix(line, CRLF)

	intVal, err := strconv.Atoi(intStr)
	if err != nil {
		return Value{}, err
	}

	return Value{Type: INTEGER, Int: intVal}, nil
}

func (p *Parser) parseBulk() (Value, error) {

	sizeStr, err := p.reader.ReadString('\n')

	if err != nil {
		return Value{}, err
	}

	sizeStr = strings.TrimSuffix(sizeStr, CRLF)

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

func (p *Parser) parseArray() (Value, error) {

	sizeStr, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	sz, err := strconv.Atoi(strings.TrimSuffix(sizeStr, CRLF))
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

func (p *Parser) Reset(rd io.Reader) {
	p.reader.Reset(rd)
}
