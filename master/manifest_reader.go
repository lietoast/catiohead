package master

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

// read key-value pairs from manifest
func readManifest(path string) map[int]string {
	config := make(map[int]string)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key, err := strconv.Atoi(strings.TrimSpace(s[:index]))
		if err != nil {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		config[key] = value
	}
	return config
}

// read all characters from a json file
func readJSON(path string) string {
	buf, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return string(buf)
}
