# PROGO #
**Pro**tein information gathering in **GO**

## Introduction ##

This program reads a CSV/TSV file with gene names and returns a TSV file containing all known PDB IDs per gene. Additionally, it creates individual TSV files with all known domains of the gene protein.

APIs used:

```text
EBI       : https://www.ebi.ac.uk/proteins/api/features
gConvert  : https://biit.cs.ut.ee/gprofiler/api/convert/convert/
```

## Installation ##

Compiled binaries are available for major operating systems and architectures in the **"releases"** section.

If your system is not listed or you prefer to compile `PROGO` on your own, make sure you have GO installed and available in your PATH. Additionally, you may need build or development tools specific to your operating system in order to compile the program.

To compile PROGO, follow these steps: 

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
Welcome to ProGo v.0.1.5-beta
Aditya Singh
Github: aditya-88

Usage of progo:
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
        Output folder location
  -skipdom
        Skip domain features
  -skippdb
        Skip PDB ID
 ```
## Description of options ##

**`-col`**      : Name of the column where the gene names are listed.

**`-delim`**    : The delimiter used in the file. By default, it uses "," as the delimiter, suitable for CSV files. You can change it to any delimiter used in your file. For example, for a TSV file, use "\t" as the delimiter.

**`-ebio`**     : Organism name as per the EBI nomenclature. For example, for "Homo sapiens," EBI uses "human".

**`-file`**     : Path to the input file.

**`-maxatt`**   : This defines the number of times PROGO should attempt to re-establish a link with the server in case of an error in the response.

**`-maxebi`**   : The EBI API enforces a strict usage limit of 20 requests per second per user. The default value of 20 concurrent EBI requests seems to work well. Modify it if you encounter errors.

**`-maxreq`**   : GO practically offers concurrency in thousands of channels, but to limit resource usage and prevent system overload, it is set to 1000 for PDB search. Keep in mind that the maximum number of GO routines running at a time will be equal to **-maxebi + -maxreq**.

**`-maxwait`**  : The program times out in 10 seconds by default. You can change this value if you expect delays in your request or do not wish to wait that long.

**`-org`**      : The gProfiler API uses different organism codes than EBI, so both need to be specified. The default value is set to "hsapiens" for "Homo sapiens".

**`-out`**      : Output folder name to store the results.
