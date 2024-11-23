package converter

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

func Convert(inputPath, outputPath string, reconvert bool) error {
	files, err := os.ReadDir(inputPath)

	if err != nil {
		return err
	}

	begin := time.Now()
	validOutputFiles := make([]string, 0, len(files))
	convertedCount := 0
	upToDateCount := 0
	deletedCount := 0

	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if ext != ".md" {
			continue
		}

		in := path.Join(inputPath, file.Name())
		out := path.Join(
			outputPath,
			fmt.Sprintf("%s.json", strings.TrimSuffix(file.Name(), ext)),
		)
		validOutputFiles = append(validOutputFiles, out)

		shouldConvert, err := shouldConvert(in, out)

		if !shouldConvert && !reconvert {
			upToDateCount += 1
			continue
		}

		article, err := Parse(in)

		if err != nil {
			return fmt.Errorf(`error parsing %s: %v`, file.Name(), err)
		}

		data, err := json.MarshalIndent(article, "", "  ")

		if err != nil {
			return err
		}

		err = os.WriteFile(out, data, os.ModePerm)

		if err != nil {
			return err
		}

		fmt.Printf("🎯 converted %s...\n", in)
		convertedCount += 1
	}

	files, err = os.ReadDir(inputPath)

	if err != nil {
		return err
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if ext != ".json" {
			continue
		}

		out := path.Join(outputPath, file.Name())

		if !slices.Contains(validOutputFiles, out) {
			fmt.Printf("⚠️ deleted %s...\n", out)
			err := os.Remove(out)

			if err != nil {
				return err
			}

			deletedCount += 1
		}
	}

	end := time.Since(begin)

	fmt.Printf("🎉 Article processing done in %s. converted: %d", end, convertedCount)

	if !reconvert {
		fmt.Printf(", up-to-date: %d", upToDateCount)
	}

	fmt.Printf(", deleted: %d\n", deletedCount)

	return nil
}

func shouldConvert(sourceFn, destFn string) (bool, error) {
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
