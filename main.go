package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
)

var count int
var extCount = make(map[string]int)

func walk(root string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	ext := filepath.Ext(root)
	ext = strings.ToLower(ext)
	extCount[ext]++

	// fmt.Println(root)
	count++
	return nil
}

var done = make(chan bool)

var program, version string

var flagV = flag.Bool("version", false, "Print version and exit")

func main() {
	program = path.Base(os.Args[0])
	if version == "" {
		version = "Unknown Version"
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s (%s)\n", program, version)
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... [STARTDIR]...\n", program)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *flagV {
		fmt.Println(program, version)
		return
	}
	args := flag.Args()

	var root string

	if len(args) == 0 {
		root = "."
	} else if len(args) == 1 {
		root = args[0]
	} else {
		flag.Usage()
		return
	}

	go markProgress()
	filepath.Walk(root, walk)
	done <- true

	fmt.Println(count)
	printExtCount(extCount)
}

func markProgress() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			fmt.Fprintf(os.Stderr, "Done!")
			return
		case <-ticker.C:
			fmt.Fprintf(os.Stderr, "Files processed %6s\n",
				humanize.Comma(int64(count)))
		}
	}
}

type pair struct {
	Key   string
	Value int
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func printExtCount(counts map[string]int) {
	fmt.Println()
	ranked := rankByExtCount(counts)
	fmt.Printf("%4s %s\n", "#", "extension")

	for _, pair := range ranked {
		fmt.Printf("%4d %s\n", pair.Value, pair.Key)
	}
}

func rankByExtCount(extFrequencies map[string]int) pairList {
	pl := make(pairList, len(extFrequencies))
	i := 0
	for k, v := range extFrequencies {
		pl[i] = pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}
