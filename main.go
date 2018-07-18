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

var bannedCount int
var fileCount int
var extCount = make(map[string]int)

var walk = walkByExt

var banned = map[string]bool{
	".doc":  true,
	".xdoc": true,
	".pdf":  true,
	".xls":  true,
	".xlsx": true,
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
	fileCount++
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

	i := strings.Index(root, sep)
	// fmt.Println("i", i)
	app1 := root[i+1:]
	// fmt.Println("app1", app1)
	j := strings.Index(app1, sep)
	// fmt.Println("j", j)
	app1 = app1[:j]
	// fmt.Println("app1:", app1)

	if *flagPaths {
		fmt.Println("root:", root)
	}

	if strings.Compare(app, app1) != 0 {
		printExtCount2(extCount)
		extCount = make(map[string]int)
		for ext := range banned {
			extCount[ext] = 0
		}
		app = app1
	}
	extCount[ext]++
	bannedCount++
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
	bannedCount++
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

	printExtCount2(extCount)
	fmt.Fprintf(os.Stderr, "banned count: %d \n", bannedCount)
	fmt.Fprintf(os.Stderr, "file count: %d \n", fileCount)
}

func markProgress() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			fmt.Fprintf(os.Stderr, "Done!\n")
			return
		case <-ticker.C:
			fmt.Fprintf(os.Stderr, "Files processed %6s\n",
				humanize.Comma(int64(bannedCount)))
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

type pair2 struct {
	Key   string
	Value int
}
type pairList2 []pair2

func (p pairList2) Len() int           { return len(p) }
func (p pairList2) Less(i, j int) bool { return p[i].Key < p[j].Key }
func (p pairList2) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var printedHeaders = false

func printExtCount2(counts map[string]int) {
	count := rankByExt(counts)
	if len(count) == 0 {
		return
	}

	if !printedHeaders {
		printedHeaders = true
		fmt.Print("application(s) ")
		for _, pair := range count {
			fmt.Print(pair.Key, " ")
		}
		fmt.Println()
	}

	fmt.Print(app, " ")
	for _, pair := range count {
		fmt.Print(pair.Value, " ")
	}
	fmt.Println()
}

func rankByExt(extFrequencies map[string]int) pairList2 {
	pl := make(pairList2, len(extFrequencies))
	i := 0
	for k, v := range extFrequencies {
		pl[i] = pair2{k, v}
		i++
	}
	sort.Sort(pl)
	return pl
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
