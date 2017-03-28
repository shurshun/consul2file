package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/consul/structs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	OUT_DIR = flag.String("o", "./tmp", "Output dir")
	PREFIX  = flag.String("p", "storage/data", "Prefix to cut")
)

func init() {
	flag.Parse()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func getLastFileIndex(name string) string {
	if fileExists(name) {
		data, err := ioutil.ReadFile(name)
		check(err)

		return string(data)
	} else {
		return "0"
	}
}

func updateIndexForFile(name, idx string) {
	err := ioutil.WriteFile(name, []byte(idx), 0600)
	check(err)
}

func main() {
	var payload []structs.DirEntry
	var lastIndex, currIndex, indexFile, key, file, path, fullPath string

	dec := json.NewDecoder(os.Stdin)

	err := dec.Decode(&payload)
	check(err)

	for _, entry := range payload {
		key = strings.TrimPrefix(string(entry.Key), *PREFIX)
		file = filepath.Base(key)
		path = filepath.Dir(filepath.Clean(*OUT_DIR)+key) + "/"
		fullPath = path + file
		indexFile = fmt.Sprintf("%s.%s.index", path, file)

		currIndex = fmt.Sprintf("%v", entry.ModifyIndex)

		lastIndex = getLastFileIndex(indexFile)

		if currIndex != lastIndex {
			log.Printf("Saving '%s' to '%s'. ModifyIndex: %s\n", string(entry.Key), fullPath, currIndex)

			err = os.MkdirAll(path, 0755)
			check(err)

			err = ioutil.WriteFile(fullPath, entry.Value, 0600)
			check(err)

			updateIndexForFile(indexFile, currIndex)
		}
	}
}
