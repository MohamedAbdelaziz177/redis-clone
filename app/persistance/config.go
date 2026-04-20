package persistance

import (
	"bufio"
	"os"
	"strings"
)

type AppendFsync string

const (
	AppendFsyncAlways   AppendFsync = "always"
	AppendFsyncEverySec AppendFsync = "everysec"
	AppendFsyncNo       AppendFsync = "no"
)

type AOFConfig struct {
	Dir            string
	AppendOnly     bool
	AppendDirName  string
	AppendFileName string
	AppendFsync    AppendFsync
}

func NewAOFConfig() *AOFConfig {
	return &AOFConfig{
		Dir:            "./../data",
		AppendOnly:     true,
		AppendDirName:  "aofDir",
		AppendFileName: "appendonly.aof",
		AppendFsync:    AppendFsyncEverySec,
	}
}

func ReadAofConfig(path string) (*AOFConfig, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var config *AOFConfig = &AOFConfig{}

	for scanner.Scan() {
		ln := scanner.Text()
		lnParts := strings.SplitN(ln, " ", 2)
		if len(lnParts) < 2 {
			continue
		}
		key := strings.ToUpper(strings.TrimSpace(lnParts[0]))
		val := strings.ToUpper(strings.TrimSpace(lnParts[1]))
		switch key {
		case "APPENDDIRNAME":
			config.AppendDirName = val
		case "APPENDONLY":
			if val == "YES" {
				config.AppendOnly = true
			} else {
				config.AppendOnly = false
			}
		case "DIR":
			config.Dir = val
		case "APPENDFILENAME":
			config.AppendFileName = val
		case "APPENDFSYNC":
			config.AppendFsync = AppendFsync(val)
		}
	}

	return config, nil
}
