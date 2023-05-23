package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	software = "ProGo"
	version  = "0.1.0-beta"
	dev      = "Aditya Singh\nGithub: aditya-88\n"
)

var (
	GoProApi       = "https://biit.cs.ut.ee/gprofiler/api/convert/convert/"
	EbiApi         = "https://www.ebi.ac.uk/proteins/api/features"
	Client         = &http.Client{Timeout: 20 * time.Second}
	filePath	   string
	columnName     string
	delimiter      string
	organism       string
	ebiOrganism    string
	outputFilePath string
	maxReqs		   uint
	maxReqsEbi	   uint
	maxWaitTime    uint
	maxAttempts	   int
	skipDomain	   bool
	skipPdb        bool
)

func init() {
	flag.StringVar(&filePath, "file", "", "Input file path (CSV/TSV/ custom delimiter)")
	flag.StringVar(&columnName, "col", "", "Column name")
	flag.StringVar(&delimiter, "delim", ",", "Delimiter")
	flag.StringVar(&organism, "org", "hsapiens", "Organism")
	flag.StringVar(&ebiOrganism, "ebio", "human", "EBI Organism")
	flag.StringVar(&outputFilePath, "out", "", "Output file path")
	flag.UintVar(&maxReqs,"maxreq", 1000, "Maximum number of requests")
	flag.UintVar(&maxReqsEbi,"maxebi", 20, "Maximum number of requests to EBI. Limited to 20 by default.")
	flag.UintVar(&maxWaitTime,"maxwait", 0, "Max seconds to wait for a response in the final attempt")
	flag.IntVar(&maxAttempts,"maxatt", 5, "Max attempts to make a request")
	flag.BoolVar(&skipDomain,"skipdom", false, "Skip domain features")
	flag.BoolVar(&skipPdb,"skippdb", false, "Skip PDB ID")
	flag.Parse()
}

func main() {
	if skipDomain && skipPdb {
		fmt.Println("Don't take my name in vain.\nYou can't skip both PDB ID and domain features.\nThat's all I do.\nI'm a one trick pony.\nI'm outta here.")
		os.Exit(0)
	}	
	outputFolder := filepath.Dir(outputFilePath)
	fmt.Printf("Welcome to %s v.%s\n%s\n", software, version, dev)
	if filePath == "" || columnName == "" || outputFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}
	var wg sync.WaitGroup
	var wgFea sync.WaitGroup
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
	fmt.Println("Maximun number of concurrent requests for PDB IDs:", maxReqs)
	fmt.Println("Maximun number of concurrent requests for domain search:", maxReqsEbi)
	fmt.Println("Total available cores:", runtime.NumCPU())
	fmt.Println("Maximun number of attempts:", maxAttempts)
	fmt.Println("Maximun number of seconds to wait for a response in the final attempt:", maxWaitTime)
	fmt.Println("Output folder:", outputFolder)
	fmt.Println("----------------------------------------")
	guard := make(chan struct{}, maxReqs)
	guardFeat := make(chan struct{}, maxReqsEbi)
	bar := progressbar.Default(int64((len(genes) * 2) + 2))
	for _, gene := range genes {
		guardFeat <- struct{}{}
		saveLoc := fmt.Sprintf("/%s/%s_features.csv", outputFolder, gene)
		//fmt.Print(saveLoc)
		wgFea.Add(1)
		if !skipDomain {
			go func(gene string, organism string, api string, saveLoc string, wg *sync.WaitGroup) {
				defer func() { <-guardFeat }()
				saveFeats(gene, organism, api, saveLoc, wg)
			}(gene, organism, GoProApi, saveLoc, &wgFea)
			bar.Add(1)
		} else {
			wgFea.Done()
			bar.Add(1)
			<-guardFeat
		}

		wg.Add(1)
		guard <- struct{}{}
		if !skipPdb {
			go func(gene string, organism string, respch chan string, wg *sync.WaitGroup, maxWait uint, maxatt int) {
			defer func() { <-guard }()
			getID(gene, organism, respch, wg, maxWait, maxatt)
		}(gene, organism, respch, &wg, maxWaitTime, maxAttempts)
		bar.Add(1)	
		} else {
			wg.Done()
			bar.Add(1)
			<-guard
		}

	}
	wg.Wait()
	close(respch)
	for resp := range respch {
		gene2pdb += resp
	}
	bar.Add(1)
	writeToFile(gene2pdb, outputFilePath)
	if !skipPdb {
		failed := len(genes)-len(strings.Split(gene2pdb, "\n"))+1
		if failed > 0 {
		fmt.Printf("Failed to get PDB ID for %d genes\n", failed)
		}
	}
	wgFea.Wait()
	bar.Add(1)

	fmt.Printf("Done!\nCheck the output file:\n%s\n", outputFilePath)
	//saveFeats("APOE", "human", EbiApi, "./test_features.csv")
}
