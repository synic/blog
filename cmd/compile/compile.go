package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

func shouldCompile(sourceFn, destFn string) (bool, error) {
	inInfo, err := os.Stat(sourceFn)

	if err != nil {
		return false, err
	}

	outInfo, err := os.Stat(destFn)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true, nil
		}
		return false, err
	}

	if inInfo.ModTime().After(outInfo.ModTime()) {
		return true, nil
	}

	return false, nil
}

func main() {
	inputDir := flag.String("i", "", "Input directory")
	outputDir := flag.String("o", "articles.json", "Output directory")
	recompile := flag.Bool("recompile", false, "Recompile all articles instead of just new ones")
	verbose := flag.Bool("v", false, "Verbose output")
	deleteRemoved := flag.Bool("d", false, "Delete json files for removed markdown files")

	flag.Parse()

	if *inputDir == "" || *outputDir == "" {
		fmt.Println("Compile markdown articles to html")
		fmt.Println("")
		fmt.Println("Usage: compile -i [sourcedir] -o [destdir] [-recompile -v]")
		fmt.Println("")

		os.Exit(1)
	}

	files, err := os.ReadDir(*inputDir)

	if err != nil {
		log.Fatal(err)
	}

	begin := time.Now()
	validOutputFiles := make([]string, 0, len(files))
	compiledCount := 0
	skippedCount := 0
	deletedCount := 0

	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if ext != ".md" {
			continue
		}

		in := path.Join(*inputDir, file.Name())
		out := path.Join(*outputDir, fmt.Sprintf("%s.json", strings.TrimSuffix(file.Name(), ext)))
		validOutputFiles = append(validOutputFiles, out)

		shouldCompile, err := shouldCompile(in, out)

		if !shouldCompile && !*recompile {
			skippedCount += 1

			if *verbose {
				fmt.Printf("skipped %s...\n", in)
			}
			continue
		}

		article, err := parseArticle(in)

		if err != nil {
			log.Fatalf(`error parsing %s: %v`, file.Name(), err)
		}

		data, err := json.MarshalIndent(article, "", "  ")

		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(out, data, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

		if *verbose {
			fmt.Printf("compiled %s...\n", in)
		}
		compiledCount += 1
	}

	if *deleteRemoved {
		files, err := os.ReadDir(*outputDir)

		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			ext := filepath.Ext(file.Name())

			if ext != ".json" {
				continue
			}

			out := path.Join(*outputDir, file.Name())

			if !slices.Contains(validOutputFiles, out) {
				err := os.Remove(out)

				if err != nil {
					log.Fatal(err)
				}

				deletedCount += 1
			}
		}
	}

	end := time.Since(begin)

	var out strings.Builder

	out.WriteString(
		fmt.Sprintf("🎉 Article processing done in %s. compiled: %d", end, compiledCount),
	)

	if !*recompile {
		out.WriteString(fmt.Sprintf(", up-to-date: %d", skippedCount))
	}

	if *deleteRemoved {
		out.WriteString(fmt.Sprintf(", deleted: %d", deletedCount))
	}

	out.WriteString("\n")
	fmt.Printf(out.String())
}
