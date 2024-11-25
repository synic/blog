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

type ConvertResult struct {
	ConvertedPaths []string
	Duration       time.Duration
	ConvertedCount int
	UpToDateCount  int
	DeletedCount   int
	reconvert      bool
}

func (r ConvertResult) String() string {
	var b strings.Builder

	b.WriteString(
		fmt.Sprintf("🎉 Article processing done in %s. converted: %d", r.Duration, r.ConvertedCount),
	)

	if !r.reconvert {
		b.WriteString(fmt.Sprintf(", up-to-date: %d", r.UpToDateCount))
	}

	b.WriteString(fmt.Sprintf(", deleted: %d\n", r.DeletedCount))

	return b.String()
}

func Convert(inputPath, outputPath string, reconvert bool) (ConvertResult, error) {
	res := ConvertResult{reconvert: reconvert}
	files, err := os.ReadDir(inputPath)

	if err != nil {
		return ConvertResult{}, err
	}

	res.ConvertedPaths = make([]string, 0, len(files))
	begin := time.Now()
	validOutputFiles := make([]string, 0, len(files))

	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if ext != ".md" {
			continue
		}

		in := path.Join(inputPath, file.Name())
		out := path.Join(outputPath, strings.TrimSuffix(file.Name(), ext)+".json")
		validOutputFiles = append(validOutputFiles, out)

		shouldConvert, err := shouldConvert(in, out)

		if !shouldConvert && !reconvert {
			res.UpToDateCount += 1
			continue
		}

		article, err := Parse(in)

		if err != nil {
			return res, fmt.Errorf(`error parsing %s: %v`, file.Name(), err)
		}

		data, err := json.MarshalIndent(article, "", "  ")

		if err != nil {
			return res, err
		}

		err = os.WriteFile(out, data, os.ModePerm)

		if err != nil {
			return res, err
		}

		res.ConvertedPaths = append(res.ConvertedPaths, in)
		res.ConvertedCount += 1
	}

	files, err = os.ReadDir(inputPath)

	if err != nil {
		return res, err
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
				return res, err
			}

			res.DeletedCount += 1
		}
	}

	res.Duration = time.Since(begin)

	return res, nil
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
