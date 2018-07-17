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

var walk = walkByExt

var banned = map[string]bool{
	".pdf": true,
	".sql": true,
}
var sep = string(filepath.Separator)
var appGroup = map[string]bool{
	"Business":              true,
	"CareSystems":           true,
	"COE":                   true,
	"Encounters":            true,
	"ETG":                   true,
	"Incubator":             true,
	"OES":                   true,
	"ONM":                   true,
	"PCE":                   true,
	"QA":                    true,
	"Reporting":             true,
	"ReportingAndAnalytics": true,
	"ResearchAndInnovation": true,
	"TPCM":                  true,
}

var app string

// walkProject collects stats per TFS project
func walkByProjects(root string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	ext := filepath.Ext(root)
	ext = strings.ToLower(ext)

	// fmt.Println(">>>", ext)
	if !banned[ext] {
		return nil
	}

	dir := strings.SplitN(root, sep, 4)
	last := 2
	if appGroup[dir[1]] {
		last = 3
	}
	app1 := strings.Join(dir[0:last], sep)
	if app == "" {
		app = app1
	}

	if *flagPaths {
		fmt.Println("root:", root)
	}

	if strings.Compare(app, app1) != 0 {
		printExtCount(extCount)
		extCount = make(map[string]int)
		app = app1
	}
	extCount[ext]++
	count++
	return nil
}

// walkByExt collects stats by file extensions
func walkByExt(root string, info os.FileInfo, err error) error {
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

var program string
var version = "0.2"

var flagV = flag.Bool("version", false, "Print version and exit")

var flagApp = flag.Bool("app", true, "Count banned files by app")
var flagPaths = flag.Bool("path", false, "Print full path of each banned file")

// Main exports main()
func Main() {
	main()
}

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

	if *flagApp {
		walk = walkByProjects
	}

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

	printExtCount(extCount)
	fmt.Println("file count:", count)
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
	if *flagApp {
		fmt.Println("app:", app)
	}

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
