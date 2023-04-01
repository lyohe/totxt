package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Read .totxtignore file and return a list of patterns to ignore.
func getIgnoreList(ignoreFilePath string) ([]string, error) {
	ignoreList := []string{}
	file, err := os.Open(ignoreFilePath)
	if err != nil {
		return ignoreList, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ignoreList = append(ignoreList, strings.TrimSpace(line))
	}
	return ignoreList, scanner.Err()
}

func shouldIgnore(filePath string, ignoreList []string) bool {
	for _, pattern := range ignoreList {
		matched, _ := filepath.Match(pattern, filePath)
		if matched {
			return true
		}
	}
	return false
}

// Recursively process a directory and write the contents to the output file.
func processDirectory(dirPath string, ignoreList []string, outputFile *os.File) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativeFilePath, _ := filepath.Rel(dirPath, path)
		if shouldIgnore(relativeFilePath, ignoreList) {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintf(outputFile, "----\n%s\n%s\n", relativeFilePath, string(content))
		return nil
	})
	return err
}

func main() {
	var (
		preambleFile string
		outputFile   string
	)
	flag.StringVar(&preambleFile, "p", "preamble.txt", "Path to the preamble file")
	flag.StringVar(&outputFile, "o", "output.txt", "Path to the output file")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: ./totxt /path/to/git/directory [-p /path/to/preamble.txt] [-o /path/to/output.txt]")
		os.Exit(1)
	}

	dirPath := flag.Arg(0)
	ignoreFilePath := filepath.Join(dirPath, ".totxtignore")

	ignoreList, err := getIgnoreList(ignoreFilePath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error reading .totxtignore file:", err)
		os.Exit(1)
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer outFile.Close()

	if preambleFile != "" {
		preambleContent, err := os.ReadFile(preambleFile)
		if err != nil {
			fmt.Println("Error reading preamble file:", err)
			os.Exit(1)
		}
		_, _ = outFile.WriteString(string(preambleContent) + "\n")
	} else {
		defaultPreamble := "The following text is a diretory structure with code. The structure of the text are sections that begin with ----, followed by a single line containing the file path and file name, followed by a variable amount of lines containing the file contents. The text representing the directory ends when the symbols --END-- are encounted. Any further text beyond --END-- are meant to be interpreted as instructions using the aforementioned directory as context.\n"
		_, _ = outFile.WriteString(defaultPreamble)
	}

	if err := processDirectory(dirPath, ignoreList, outFile); err != nil {
		fmt.Println("Error processing directory:", err)
		os.Exit(1)
	}

	_, _ = outFile.WriteString("--END--")
	fmt.Printf("Directory contents written to %s.\n", outputFile)
}
