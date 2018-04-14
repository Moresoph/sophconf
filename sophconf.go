package sophconf

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

const (
	Separator = " "
)

func parseKeyValue(line string) (key, value string, err error) {
	// remove comments
	keyAndValue := strings.Split(line, Separator)
	if len(keyAndValue) != 2 {
		err = fmt.Errorf("can't split line=%v with space", line)
		return
	}
	key = keyAndValue[0]
	value = keyAndValue[1]
	return
}

func getOneLine(scanner *bufio.Scanner) (line string, ok bool) {
	ok = false
	var retLine string
	for scanner.Scan() {
		ok = true
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, "\\") {
			// merge line ends with backslash
			line = line[:len(line)-1]
			retLine += line
			continue
		}
		retLine += line
		break
	}
	return retLine, ok
}

func LoadConfFile(filename string) (ret map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	ret = make(map[string]string)
	err = nil

	basename := path.Dir(filename)
	scanner := bufio.NewScanner(file)
	lastKey := ""
	for {
		line, ok := getOneLine(scanner)
		if ok != true {
			break
		}

		key, value, err := parseKeyValue(line)
		if err != nil {
			return nil, err
		}
		if len(key) > 0 && len(value) > 0 {
			ret[key] = value
			lastKey = key
		}
	}

	if lastKey == "include" {
		includeFile := ret[lastKey]
		if strings.HasPrefix(includeFile, "/") == false {
			includeFile = basename + "/" + includeFile
		}
		includeConf, err := LoadConfFile(includeFile)
		if err != nil {
			err := fmt.Errorf("load include file=%v failed err=%v", includeFile, err)
			return nil, err
		}
		// merge sub conf
		for k, v := range includeConf {
			ret[k] = v
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return
}
