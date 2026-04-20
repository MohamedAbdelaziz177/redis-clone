package persistance

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type AOF struct {
	file *os.File
	rd   *bufio.Reader
	wr   *bufio.Writer
	mu   *sync.RWMutex
}

func NewAOF(config *AOFConfig) (*AOF, error) {

	os.MkdirAll(path.Join(config.Dir, config.AppendDirName), 0755)

	filepath := path.Join(config.Dir, config.AppendDirName, config.AppendFileName)
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)

	if err != nil {
		return nil, fmt.Errorf("Error Openning The AOF: %s", err)
	}

	return &AOF{
		file: file,
		rd:   bufio.NewReader(file),
		wr:   bufio.NewWriter(file),
		mu:   &sync.RWMutex{},
	}, nil
}

func (aof *AOF) Append(value *resp.Value) error {

	aof.mu.Lock()
	defer aof.mu.Unlock()

	entry := aof.serializeValue(value)
	_, err := aof.wr.Write(entry)

	if err != nil {
		return err
	}

	return nil
}

func (aof *AOF) ReadEntries() []resp.Value {

	aof.mu.RLock()
	defer aof.mu.RUnlock()

	values := make([]resp.Value, 0)
	parser := resp.NewParser(aof.file)

	for {

		val, err := parser.ParseValue()

		if err == io.EOF {
			break
		}

		if err != nil {
			return values
		}

		values = append(values, val)
	}

	return values
}

func (aof *AOF) serializeValue(value *resp.Value) []byte {

	bulksArr := make([]string, 0)

	if value.Type == resp.ARRAY && len(value.Array) > 0 {
		for _, val := range value.Array {
			bulksArr = append(bulksArr, val.Bulk)
		}
	}

	return resp.EncodeArray(bulksArr)
}
