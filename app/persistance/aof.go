package persistance

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/MohamedAbdelaziz177/redis-clone/app/resp"
)

type AOF struct {
	file *os.File
	Wr   *bufio.Writer
	mu   *sync.RWMutex
	p    *resp.Parser
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
		Wr:   bufio.NewWriter(file),
		mu:   &sync.RWMutex{},
		p:    resp.NewParser(file),
	}, nil
}

func (aof *AOF) Append(value *resp.Value) error {

	aof.mu.Lock()
	defer aof.mu.Unlock()

	entry := aof.serializeValue(value)
	_, err := aof.Wr.Write(entry)

	if err != nil {
		return err
	}

	return nil
}

func (aof *AOF) ReadEntries() ([]resp.Value, error) {

	aof.mu.Lock()
	defer aof.mu.Unlock()

	if _, err := aof.file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	aof.p.Reset(aof.file)

	values := make([]resp.Value, 0)
	for {

		val, err := aof.p.ParseValue()

		if err == io.EOF {
			break
		}

		if err != nil {
			return values, err
		}

		values = append(values, val)
	}

	if _, err := aof.file.Seek(0, io.SeekEnd); err != nil {
		return values, err
	}

	return values, nil
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
