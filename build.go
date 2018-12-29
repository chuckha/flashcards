package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var clean bytes.Buffer

	out := "vocab.txt"

	var vocab bytes.Buffer
	w := csv.NewWriter(&vocab)
	w.Comma = '\t'
	for _, arg := range os.Args[1:] {
		// read the file
		f, _ := os.Open(arg)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			txt := scanner.Text()
			// skip empty lines
			if len(strings.TrimSpace(txt)) == 0 {
				continue
			}
			// skip comments
			if strings.HasPrefix(txt, "#") {
				continue
			}
			// write any other line into a buffer
			fmt.Fprintln(&clean, txt)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		r := csv.NewReader(&clean)
		records, err := r.ReadAll()
		if err != nil {
			panic(err)
		}
		basename := filepath.Base(arg)
		tag := strings.TrimSuffix(basename, filepath.Ext(basename))
		for _, record := range records {
			record = append(record, tag)
			w.Write(record)
		}
		w.Flush()
	}

	if err := ioutil.WriteFile(out, vocab.Bytes(), os.FileMode(int(0644))); err != nil {
		fmt.Println(err)
	}
}
