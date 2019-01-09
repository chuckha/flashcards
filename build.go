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

	"github.com/pkg/errors"
)

func walker(collector *csv.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "vocab.txt") {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		var clean bytes.Buffer
		f, err := os.Open(path)
		if err != nil {
			return errors.Wrap(err, "failed to open: "+path)
		}
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
			return errors.Wrap(err, path)
		}

		r := csv.NewReader(&clean)
		records, err := r.ReadAll()
		if err != nil {
			return errors.Wrap(err, path)
		}
		basename := filepath.Base(path)
		tag := strings.TrimSuffix(basename, filepath.Ext(basename))
		for _, record := range records {
			record = append(record, tag)
			collector.Write(record)
		}
		collector.Flush()
		return nil
	}
}

func main() {
	out := os.Args[1] + ".txt"
	var collection bytes.Buffer
	w := csv.NewWriter(&collection)
	w.Comma = '\t'

	if err := filepath.Walk(os.Args[1], walker(w)); err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	fmt.Println(filepath.Join(os.Args[1], out))
	if err := ioutil.WriteFile(filepath.Join(os.Args[1], out), collection.Bytes(), os.FileMode(int(0644))); err != nil {
		fmt.Println(err)
	}
}
