package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	software = "ProGo"
	version  = "0.0.1"
	dev      = "Aditya Singh\nGithub: aditya-88\n"
)

var (
	GoProApi       = "https://biit.cs.ut.ee/gprofiler/api/convert/convert/"
	Client         = &http.Client{Timeout: 20 * time.Second}
	Pdbfile        struct{ Result []struct{ Converted string `json:"converted"` } `json:"result"` }
	ProtFeatures   struct{ Features []struct{ Begin, End int; Name string } `json:"features"` }
	filePath	   string
	columnName     string
	delimiter      string
	organism       string
	outputFilePath string
	maxReqs		   uint
	maxWaitTime    uint
	maxAttempts	int
)

func init() {
	flag.StringVar(&filePath, "file", "", "File path")
	flag.StringVar(&columnName, "col", "", "Column name")
	flag.StringVar(&delimiter, "delim", ",", "Delimiter")
	flag.StringVar(&organism, "org", "", "Organism")
	flag.StringVar(&outputFilePath, "out", "", "Output file path")
	flag.UintVar(&maxReqs,"maxreq", 1000, "Maximum number of requests")
	flag.UintVar(&maxWaitTime,"maxwait", 0, "Max seconds to wait for a response in the final attempt")
	flag.IntVar(&maxAttempts,"maxatt", 5, "Max attempts to make a request")
	flag.Parse()
}

func main() {
	fmt.Printf("Welcome to %s v.%s\n%s\n", software, version, dev)

	if filePath == "" || columnName == "" || organism == "" || outputFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Check if output file already exists, if so, read it to a variable as string

	var wg sync.WaitGroup
	genes := parseCSV(filePath, columnName, delimiter)
	respch := make(chan string, len(genes))
	var gene2pdb string
	fmt.Println("----------------------------------------")
	fmt.Println("------------Basic config----------------")
	fmt.Println("----------------------------------------")
	fmt.Println("Organism:", organism)
	fmt.Println("Total number of input genes:", len(genes))
	// Remove any empty/ duplicate genes
	genes = removeEmpty(genes)
	genes = removeDuplicates(genes)
	fmt.Println("Total number of genes after dropping empty and duplicates:", len(genes))
	fmt.Println("Maximun number of concurrent requests:", maxReqs)
	fmt.Println("Total available cores:", runtime.NumCPU())
	fmt.Println("Maximun number of attempts:", maxAttempts)
	fmt.Println("Maximun number of seconds to wait for a response in the final attempt:", maxWaitTime)
	fmt.Println("----------------------------------------")
	//wg.Add(len(genes))
	guard := make(chan struct{}, maxReqs)
	bar := progressbar.Default(int64(len(genes) + 2))
	for _, gene := range genes {
		guard <- struct{}{}
		wg.Add(1)
		go func(gene string, organism string, respch chan string, wg *sync.WaitGroup, maxWait uint, maxatt int) {
			defer func() { <-guard }()
			getID(gene, organism, respch, wg, maxWait, maxatt)
		}(gene, organism, respch, &wg, maxWaitTime, maxAttempts)
		bar.Add(1)
	}
	wg.Wait()
	close(respch)

	for resp := range respch {
		gene2pdb += resp
	}
	bar.Add(1)
	writeToFile(gene2pdb, outputFilePath)
	bar.Add(1)
	failed := len(genes)-len(strings.Split(gene2pdb, "\n"))+1
	if failed > 0 {
		fmt.Printf("Failed to get PDB ID for %d genes\n", failed)
	}
	fmt.Printf("Done!\nCheck the output file:\n%s\n", outputFilePath)
}
