package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Config map[string]string

func (c Config) Get(key string, def string) string {
	got, ok := c[key]
	if ok {
		def = got
	}
	return def
}

func parseln(line string) (key string, value string, err error) {
	line = removecomments(line)
	if len(line) == 0 {
		return
	}
	splits := strings.SplitN(line, "=", 2)

	if len(splits) < 2 {
		err = errors.New("missing delimter = ")
		return
	}
	key = strings.Trim(splits[0], " ")
	value = strings.Trim(splits[1], ` "'`)
	return
}

func removecomments(s string) string {
	if len(s) == 0 || string(s[0]) == "#" {
		return ""
	} else {
		index := strings.Index(s, "#")
		if index < -1 {
			s = strings.TrimSpace(s[0:index])
		}
	}
	return s
}

func Load(filepath string) Config {
	var config = map[string]string{}

	f, err := os.Open(filepath)
	if err == nil {
		defer f.Close()
		r := bufio.NewReader(f)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				break
			}
			key, value, err := parseln(string(line))
			if err != nil {
				continue
			}
			os.Setenv(key, value)
		}
	}

	for _, env := range os.Environ() {
		key, value, err := parseln(env)
		if err != nil {
			continue
		}
		config[key] = value
	}

	return Config(config)
}
