package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	s "strings"
	"sync"
	"sync/atomic"

	"github.com/heyyakash/go-grep/structs"
)

var wg sync.WaitGroup
var ops atomic.Uint64

const stopSignal = "stopSignal"

func HandleErr(e error) {
	if e != nil {
		panic(e)
	}
}

func findPattern(jobs <-chan string, pattern string, res *structs.Result) {
	for job := range jobs {
		if job == stopSignal {
			wg.Done()
			return
		}
		file, err := os.Open(job)
		HandleErr(err)
		scanner := bufio.NewScanner(file)
		i := 0
		for scanner.Scan() {
			line := scanner.Text()
			i++
			if s.Contains(line, pattern) {
				res.AddLine(fmt.Sprintf("File : %s \t Line %d : \t %s", job, i, line))
				ops.Add(1)
			}
		}
		file.Close()
	}
}

func main() {
	wordPtr := flag.String("word", "test", "The string to be searched")
	poolPtr := flag.Int("workers", 10, "The no. of workers for searching")
	flag.Parse()

	l, err := os.ReadDir(".")
	HandleErr(err)
	files := []string{}
	for _, v := range l {
		if !v.IsDir() {
			files = append(files, v.Name())
		}
	}
	jobs := make(chan string, len(files))
	res := structs.NewResultHolder()
	for i := 0; i < *poolPtr; i++ {
		wg.Add(1)
		go findPattern(jobs, *wordPtr, res)
	}

	for _, v := range files {
		jobs <- v
	}

	for i := 0; i < *poolPtr; i++ {
		jobs <- stopSignal
	}

	close(jobs)
	wg.Wait()

	lines := res.GetLines()
	for _, v := range lines {
		fmt.Println(v)
	}
}
