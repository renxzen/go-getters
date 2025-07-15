package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/renxzen/go-getters/pkg/generator"
	"github.com/renxzen/go-getters/pkg/parser"
)

var (
	inputPath   = flag.String("input", ".", "Path to directory containing Go files")
	outputFile  = flag.String("output", "getters.gen.go", "Output file path")
	structNames = flag.String("structs", "", "Comma-separated list of struct names to generate getters for")
	help        = flag.Bool("help", false, "Show help message")
)

func main() {
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *structNames == "" {
		fmt.Fprintf(os.Stderr, "Error: -structs flag is required\n")
		showHelp()
		os.Exit(1)
	}

	// Parse struct names
	structs := strings.Split(*structNames, ",")
	for i, s := range structs {
		structs[i] = strings.TrimSpace(s)
	}

	// Parse the directory to extract struct information
	p := parser.New()
	result, err := p.ParseDirectory(*inputPath)
	if err != nil {
		log.Fatalf("Failed to parse directory: %v", err)
	}

	// Generate getters
	gen := generator.New()
	outBytes, err := gen.GenerateGetters(structs, result)
	if err != nil {
		log.Fatalf("Failed to generate getters: %v", err)
	}

	// Write output
	if err := os.WriteFile(*outputFile, outBytes, 0644); err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	fmt.Printf("Generated getters for %d struct(s) in %s\n", len(structs), *outputFile)
}

func showHelp() {
	fmt.Printf(`go-getters - Generate getter methods for Go structs

Usage:
  %s [options]

Options:
  -input string
        Path to directory containing Go files (default ".")
  -output string
        Output file name (default "getters.gen.go"). The file should be created in the same directory as the input directory as methods should live in the same package as the struct.
  -structs string
        Comma-separated list of struct names to generate getters for (required)
  -help
        Show this help message

Examples:
  %s -structs="User,Product"
  %s -input=./models -output=getters.go -structs="User,Product,Order"

`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
}
