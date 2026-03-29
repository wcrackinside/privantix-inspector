package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"privantix-source-inspector/catalog"
)

func main() {
	input := flag.String("input", "", "Path to analysis.json from privantix-inspector (required)")
	output := flag.String("output", "./catalog_output", "Output directory for catalog files")
	flag.Parse()

	if *input == "" {
		fmt.Println("Usage: privantix-catalog --input <path-to-analysis.json> [--output <dir>]")
		fmt.Println("")
		fmt.Println("Transforms privantix-inspector output into a structured dataset catalog.")
		fmt.Println("Generates: catalog.json, catalog_datasets.csv, catalog_columns.csv, catalog.html")
		os.Exit(1)
	}

	if _, err := os.Stat(*input); os.IsNotExist(err) {
		log.Fatalf("input file not found: %s", *input)
	}

	started := time.Now()
	log.Printf("reading analysis from %s", *input)
	result, err := catalog.ParseAnalysis(*input)
	if err != nil {
		log.Fatalf("parse error: %v", err)
	}

	c := catalog.Build(result)
	log.Printf("catalog: %d datasets, %d columns, %d rules", c.TotalDatasets, c.TotalColumns, c.TotalRules)

	if err := catalog.ExportAll(*output, c); err != nil {
		log.Fatalf("export error: %v", err)
	}

	duration := time.Since(started).Seconds()
	log.Printf("catalog written to %s (%.2fs)", *output, duration)
	fmt.Printf("✔ Catalog generated: %s/catalog.html (%.2fs)\n", *output, duration)
}
