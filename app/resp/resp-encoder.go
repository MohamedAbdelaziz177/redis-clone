package resp

import (
	"strconv"
)

func EncodeString(value string) []byte {
	return []byte("+" + value + "\r\n")
}

func EncodeError(value string) []byte {
	return []byte("-" + value + "\r\n")
}

func EncodeInteger(value int) []byte {
	return []byte(":" + strconv.Itoa(value) + "\r\n")
}

func EncodeBulk(value string) []byte {
	return []byte("$" + strconv.Itoa(len(value)) + "\r\n" + value + "\r\n")
}

func EncodeArray(values []string) []byte {
	var result []byte
	result = append(result, []byte("*"+strconv.Itoa(len(values))+"\r\n")...)
	for _, val := range values {
		result = append(result, EncodeBulk(val)...)
	}
	return result
}
