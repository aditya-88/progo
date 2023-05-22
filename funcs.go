package main

import (
	"encoding/csv"
	"encoding/json"
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
			err = json.NewDecoder(resp.Body).Decode(&Pdbfile)
			if err != nil {
				panic(err)
			}
			break
		}
		if i == maxAtt-1 {
			//respch <- gene + "\tErrorCode_" + resp.Status + "\n"
			wg.Done()
			return
		}
	}
	names := ""
	for i, name := range Pdbfile.Result {
		if i == len(Pdbfile.Result)-1 {
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

func writeToFile(gene2pdb string, path string) { // This function writes the PDB ID to a file
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(gene2pdb)
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
