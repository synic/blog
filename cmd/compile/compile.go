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
	count := 0
	skipCount := 0

	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if ext != ".md" {
			continue
		}

		in := path.Join(*inputDir, file.Name())
		out := path.Join(*outputDir, fmt.Sprintf("%s.json", strings.TrimSuffix(file.Name(), ext)))

		shouldCompile, err := shouldCompile(in, out)

		if !shouldCompile && !*recompile {
			skipCount += 1

			if *verbose {
				fmt.Printf("skipped %s...\n", in)
			}
			continue
		}

		article, err := parseArticle(in)

		if err != nil {
			log.Fatal(err)
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
		count += 1
	}

	end := time.Since(begin)

	if count <= 0 {
		fmt.Printf("Checked %d files in %s, but all were up-to-date.\n", skipCount, end)
		return
	}

	if *recompile {
		fmt.Printf("Done! Compiled %d articles in %s.\n", count, end)
	} else {
		fmt.Printf("Done! Compiled %d articles in %s, and skipped %d that were up-to-date.\n", count, end, skipCount)
	}
}
