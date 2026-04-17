package resp

import (
	"strconv"
)

const (
	CRLF = "\r\n"
)

func EncodeString(value string) []byte {
	return []byte("+" + value + CRLF)
}

func EncodeError(value string) []byte {
	return []byte("-" + value + CRLF)
}

func EncodeInteger(value int) []byte {
	return []byte(":" + strconv.Itoa(value) + CRLF)
}

func EncodeBulk(value string) []byte {
	return []byte("$" + strconv.Itoa(len(value)) + CRLF + value + CRLF)
}

func EncodeArray(values []string) []byte {

	if len(values) == 0 {
		return []byte("*0\r\n")
	}

	var result []byte
	result = append(result, []byte("*"+strconv.Itoa(len(values))+CRLF)...)

	for _, val := range values {
		result = append(result, EncodeBulk(val)...)
	}

	return result
}
