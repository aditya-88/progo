# PROGO #
**Pro**tein information gathering in **GO**

## Introduction ##

This program reads a `CSV/TSV` file with gene names and returns a `TSV` file of all known PDB IDs per gene.

Additionaly, it creats individual `TSV` files with all known domains of the gene protein.

## Installation ##

There are compiled binaries available for major OSes and archs in the `releases`.

In case your system isn't listed or you wish to compile PROGO on your own, make sure your have `GO` installed and available in `PATH`

Also, you might need `build` or `dev` tools specific to your OS in order to compile the program. 

```bash
git clone https://github.com/aditya-88/progo && \
cd progo && \
go build ./
```
You can also run the program without compiling.

```bash
cd progo && \
go run ./
```

## Usage ##

```bash
Welcome to ProGo v.0.1.0-beta
Aditya Singh
Github: aditya-88

Usage of /Users/aditya/Codes/progo/bin/progo_macos_arm64:
  -col string
    	Column name
  -delim string
    	Delimiter (default ",")
  -ebio string
    	EBI Organism (default "human")
  -file string
    	Input file path (CSV/TSV/ custom delimiter)
  -maxatt int
    	Max attempts to make a request (default 5)
  -maxebi uint
    	Maximum number of requests to EBI. Limited to 20 by default. (default 20)
  -maxreq uint
    	Maximum number of requests (default 1000)
  -maxwait uint
    	Max seconds to wait for a response in the final attempt
  -org string
    	Organism (default "hsapiens")
  -out string
    	Output file path
 ```
