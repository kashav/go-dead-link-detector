package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kshvmdn/ded"
)

const defaultOutputFormat = `{{ .Filename }}:{{ .Line }}:{{ .Column }}: {{ .URL }} -> {{ .Result }}`

var (
	stdout     *log.Logger
	outputTmpl *template.Template

	version = "devel"

	format      = flag.String("f", "", fmt.Sprintf("Output format (default %q)", defaultOutputFormat))
	hideOK      = flag.Bool("n", false, "Hide \"200 OK\" responses")
	ignore      = flag.String("i", "", "Regular expression to ignore")
	outFile     = flag.String("o", "stdout", "Output file or [stderr|stdout]")
	showVersion = flag.Bool("v", false, "Show version and exit")
	workers     = flag.Int("w", 0, "Number of workers (0 = number of CPUs)")
)

func pathWorker(paths chan string, matches chan<- ded.Match, quit chan<- int) {
	for path := range paths {
		raw, err := ded.ReadTextFile(path)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(raw) == 0 {
			continue
		}
		for _, match := range ded.Parse(path, raw, *ignore) {
			matches <- match
		}
	}
	quit <- 0
}

func matchWorker(matches <-chan ded.Match, quit chan<- int) {
	var (
		count  int
		output bytes.Buffer
	)
	for match := range matches {
		count++
		match.Check()
		if match.Result == "200 OK" && *hideOK {
			continue
		}
		outputTmpl.Execute(&output, match)
		stdout.Println(output.String())
		output.Reset()
	}
	quit <- count
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("ded version %s\n", version)
		return
	}

	if len(*format) > 0 {
		t, err := template.New("custom").Parse(*format)
		if err != nil {
			log.Fatalf("unable to compile log format: %s", err)
		}
		outputTmpl = t
	} else {
		outputTmpl = template.Must(template.New("outputTmpl").Parse(defaultOutputFormat))
	}

	switch {
	case *outFile == "/dev/stderr" || *outFile == "stderr":
		stdout = log.New(os.Stderr, "", 0)
	case *outFile == "/dev/stdout" || *outFile == "stdout":
		fallthrough
	case *outFile == "" || *outFile == "-":
		stdout = log.New(os.Stdout, "", 0)
	default:
		fi, err := os.Create(*outFile)
		if err != nil {
			log.Fatalf("unable to create outfile %q: %s", *outFile, err)
		}
		defer fi.Close()
		stdout = log.New(fi, "", 0)
	}

	if *workers < 0 {
		log.Fatal("expected non-negative number of workers")
	}
	if *workers == 0 {
		*workers = runtime.NumCPU()
	}
	if *workers == 1 {
		*workers = 2
	}

	var (
		quit    = make(chan int, *workers)
		paths   = make(chan string, 64)
		matches = make(chan ded.Match, 8)
	)

	var (
		mid = float64(*workers) / 2

		pWorkers = int(math.Floor(mid))
		mWorkers = int(math.Ceil(mid))
	)

	for i := 0; i < *workers; i++ {
		if pWorkers < i {
			go matchWorker(matches, quit)
		} else {
			go pathWorker(paths, matches, quit)
		}
	}

	defer func() {
		for i := 0; i < *workers; i++ {
			<-quit
			if i == mWorkers {
				close(matches)
			}
		}
	}()

	if flag.NArg() == 0 {
		close(paths)

		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Printf("unable to read from stdin: %s", err)
			return
		}

		for _, match := range ded.Parse("stdin", string(b), *ignore) {
			matches <- match
		}

		return
	}

	for _, filename := range flag.Args() {
		filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				paths <- path
			}
			return nil
		})
	}
	close(paths)
}
