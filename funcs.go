package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

func getID(gene string, organism string, respch chan string, wg *sync.WaitGroup, maxWait uint, maxAtt int) { // This function returns the PDB ID of the gene
	type Request struct {
		Organism string `json:"organism"`
		Query    string `json:"query"`
		Target   string `json:"target"`
	}
	var pdbFile struct {
		Result []struct {
			Converted string `json:"converted"`
		} `json:"result"`
	}
	var request = new(Request)
	request.Organism = organism
	request.Query = gene
	request.Target = "PDB"
	send_request, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	for i := 0; i < maxAtt; i++ {
		resp, err := Client.Post(GoProApi, "application/json", strings.NewReader(string(send_request)))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			err = json.NewDecoder(resp.Body).Decode(&pdbFile)
			if err != nil {
				panic(err)
			}
			break
		}
		if i == maxAtt-1 {
			wg.Done()
			return
		}
	}
	names := ""
	for i, name := range pdbFile.Result {
		if i == len(pdbFile.Result)-1 {
			names = names + name.Converted
			break
		}
		names = names + name.Converted + ","
	}
	// If the names have None followed by a comma, delete "None"
	names = strings.Replace(names, "None,", "", -1)
	names = gene + "\t" + names + "\n"
	respch <- names
	wg.Done()
}

func writeToFile(varList string, path string) { // This function writes the PDB ID to a file
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Add a header to the file
	_, _ = file.WriteString("Gene\tPDB_ID\n")
	_, err = file.WriteString(varList)
	if err != nil {
		panic(err)
	}
}

// Function to read the CSV file and return the column as a string with values separated by spaces
func parseCSV(file string, columnName string, delim string) []string {
	// Open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// Create a new reader
	r := csv.NewReader(f)
	// Set the delimiter
	r.Comma = []rune(delim)[0]
	// Read all the records
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	// Get the column index
	var columnIndex int
	for i, record := range records {
		if i == 0 {
			for j, column := range record {
				if column == columnName {
					columnIndex = j
					break
				}
			}
		}
	}
	// Get the column string
	var columnString []string
	for i, record := range records {
		if i == 0 {
			continue
		}
		columnString = append(columnString, record[columnIndex])
	}
	// Return the column string
	return columnString
}

func removeDuplicates(genes []string) []string { // This function removes duplicates from the slice of strings
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range genes {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func removeEmpty(genes []string) []string { // This function removes empty strings from the slice of strings
	var list []string
	for _, gene := range genes {
		if gene != "" {
			list = append(list, gene)
		}
	}
	return list
}

func saveFeats(gene string, organism string, api string, saveLoc string, wg *sync.WaitGroup) {
	var protFeatures []struct {
		Features []struct {
			Begin       string `json:"begin"`
			End         string `json:"end"`
			Description string `json:"description"`
		} `json:"features"`
	}
	request := fmt.Sprintf("%v?offset=0&size=100&reviewed=true&gene=%v&organism=%v&types=DOMAIN", api, gene, organism)
	resp, err := Client.Get(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		err = json.NewDecoder(resp.Body).Decode(&protFeatures)
		if err != nil {
			panic(err)
		}
	}
	if len(protFeatures) > 0 {
		// Write to file
		file, err := os.Create(saveLoc)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		_, err = file.WriteString("Gene\tStart\tEnd\tDomain\n")
		if err != nil {
			panic(err)
		}

		for _, feature := range protFeatures[0].Features {
			_, err = file.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\n", gene, feature.Begin, feature.End, feature.Description))
			if err != nil {
				panic(err)
			}
		}
	}
	wg.Done()
}
